package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/bwmarrin/discordgo"
	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var s *discordgo.Session
var f *firestore.Client

type Kill struct {
	Killer string    `firestore:"killer,omitempty"`
	Victim string    `firestore:"victim,omitempty"`
	Reason string    `firestore:"reason,omitempty"`
	Date   time.Time `firestore:"date,omitempty"`
}

type PlayerCount struct {
	Player string
	Count  int
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "serve",
		Short:   "Serves Tarkov TK Discord bot",
		Example: "tarkov-tk-bot serve",
		RunE:    serveRun,
	}
}

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"tklog": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			killer := i.ApplicationCommandData().Options[0].UserValue(s)
			victim := i.ApplicationCommandData().Options[1].UserValue(s)
			reason := ""

			if len(i.ApplicationCommandData().Options) > 2 {
				reason = i.ApplicationCommandData().Options[2].StringValue()
			}

			_, _, err := f.Collection("kills").Add(context.Background(), map[string]interface{}{
				"date":     time.Now(),
				"killer":   killer.ID,
				"victim":   victim.ID,
				"reason":   reason,
				"serverId": i.GuildID,
			})
			if err != nil {
				log.Printf("err: failed logging kill: %v\n", err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Could not log kill by **%s** on **%s**. Please try again\n", killer.Username, victim.Username),
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Kill by **%s** on **%s** logged\n", killer.Username, victim.Username),
				},
			})
		},
		"tkkills": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx := context.Background()
			iter := f.Collection("kills").Where("serverId", "==", i.GuildID).Documents(ctx)

			var kills []Kill

			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Printf("err: failed getting kills: %v\n", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Could not get kills for server. Please try again\n"),
						},
					})
					return
				}
				var k Kill
				doc.DataTo(&k)
				kills = append(kills, k)
			}
			players := getPlayersFromKills(kills, true)

			playerKills := []PlayerCount{}

			for _, p := range players {
				kc := 0
				for _, k := range kills {
					if k.Killer == p {
						kc++
					}
				}
				playerKills = append(playerKills, PlayerCount{Player: p, Count: kc})
			}

			sort.Slice(playerKills, func(i, j int) bool {
				return playerKills[i].Count > playerKills[j].Count
			})

			msgformat := "**Most Team Kills**\n"
			margs := []interface{}{}

			for x, k := range playerKills {
				playerName := ""
				player, err := s.GuildMember(i.GuildID, k.Player)
				if err != nil {
					fmt.Printf("Failed getting member: %v\n", err)
					return
				}
				if player.Nick != "" {
					playerName = player.Nick
				} else {
					playerName = player.User.Username
				}

				margs = append(margs, x+1, playerName, k.Count)
				msgformat += "%v. **%s** - %v TKs\n"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(msgformat, margs...),
				},
			})
		},
		"tkdeaths": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx := context.Background()
			iter := f.Collection("kills").Where("serverId", "==", i.GuildID).Documents(ctx)

			var kills []Kill

			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Printf("Failed getting kills: %v\n", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Could not get kills for server. Please try again\n"),
						},
					})
					return
				}
				var k Kill
				doc.DataTo(&k)
				kills = append(kills, k)
			}
			players := getPlayersFromKills(kills, false)

			playerDeaths := []PlayerCount{}

			for _, p := range players {
				dc := 0
				for _, k := range kills {
					if k.Victim == p {
						dc++
					}
				}
				playerDeaths = append(playerDeaths, PlayerCount{Player: p, Count: dc})
			}

			sort.Slice(playerDeaths, func(i, j int) bool {
				return playerDeaths[i].Count > playerDeaths[j].Count
			})

			msgformat := "**Most Team Deaths**\n"
			margs := []interface{}{}

			for x, k := range playerDeaths {
				playerName := ""
				player, err := s.GuildMember(i.GuildID, k.Player)
				if err != nil {
					fmt.Printf("Failed getting member: %v\n", err)
					return
				}
				if player.Nick != "" {
					playerName = player.Nick
				} else {
					playerName = player.User.Username
				}

				margs = append(margs, x+1, playerName, k.Count)
				msgformat += "%v. **%s** - %v TDs\n"
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(msgformat, margs...),
				},
			})
		},
		"tkstats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			users, err := s.GuildMembers(i.GuildID, "", 1000)
			if err != nil {
				fmt.Printf("Failed getting guild members: %v\n", err)
				s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Could not get kills for server. Please try again\n"))
				return
			}

			var killer *discordgo.User
			if len(i.ApplicationCommandData().Options) > 0 {
				killer = i.ApplicationCommandData().Options[0].UserValue(s)
			}
			ctx := context.Background()
			var iter *firestore.DocumentIterator
			if killer != nil {
				iter = f.Collection("kills").Where("serverId", "==", i.GuildID).Where("killer", "==", killer.ID).OrderBy("date", firestore.Desc).Documents(ctx)
			} else {
				iter = f.Collection("kills").Where("serverId", "==", i.GuildID).OrderBy("date", firestore.Desc).Documents(ctx)
			}

			kills := []Kill{}

			// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
			// 	Data: &discordgo.InteractionResponseData{
			// 		Content: fmt.Sprintf("Generating CSV..."),
			// 	},
			// })

			// var wg sync.WaitGroup

			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Printf("Failed getting kills: %v\n", err)
					s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Could not get kills for server. Please try again\n"))
					return
				}
				// wg.Add(1)
				// go func() {
				// defer wg.Done()
				var k Kill
				doc.DataTo(&k)
				for _, kU := range users {
					if kU.User.ID == k.Killer {
						if kU.Nick != "" {
							k.Killer = kU.Nick
						} else {
							k.Killer = kU.User.Username
						}
					}
				}
				for _, vU := range users {
					if vU.User.ID == k.Killer {
						if vU.Nick != "" {
							k.Killer = vU.Nick
						} else {
							k.Killer = vU.User.Username
						}
					}
				}
				kills = append(kills, k)
				// }()
			}

			// wg.Wait()

			csvBuffer := &bytes.Buffer{}
			err = gocsv.Marshal(kills, csvBuffer)

			if err != nil {
				fmt.Printf("Failed marshalling data: %v\n", err)
				return
			}
			guild, _ := s.Guild(i.GuildID)

			var fileName string

			if killer != nil {
				fileName = fmt.Sprintf("%s-%s-TarkovTKStats-%s.csv", killer.Username, guild.Name, time.Now().Format("2006-01-02"))
			} else {
				fileName = fmt.Sprintf("%s-TarkovTKStats-%s.csv", guild.Name, time.Now().Format("2006-01-02"))
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Successfully generated server stats"),
					Files: []*discordgo.File{
						{
							Name:        fileName,
							ContentType: "text/csv",
							Reader:      csvBuffer,
						},
					},
				},
			})
		},
		"tkreset": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx := context.Background()
			iter := f.Collection("kills").Where("serverId", "==", i.GuildID).Documents(ctx)

			for {
				doc, err := iter.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Printf("Failed resetting server data: %v\n", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Could not reset kills for server. Please try again\n"),
						},
					})
					return
				}
				_, err = doc.Ref.Delete(ctx)
				if err != nil {
					fmt.Printf("Failed resetting server data: %v\n", err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Could not reset kills for server. Please try again\n"),
						},
					})
					return
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Serve data successfully reset"),
				},
			})
		},
		"tkinfo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("**Thank you for using Tarkov TK**\n\nTarkov TK is still a work in progress, so if you have any suggestions or issues, pleaset let me know via Twitter https://twitter.com/KyleShepherdDev\n\nAlso if you enjoy the bot and want to support the development and maintenance, any help would be appreciated https://patreon.com/tarkovtk"),
				},
			})
		},
		"tkremove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Test"),
				},
			})
		},
	}
)

func serveRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	opt := option.WithCredentialsFile("./service-account-file.json")
	conf := &firebase.Config{ProjectID: "tarkov-tk-dev"}
	app, err := firebase.NewApp(ctx, conf, opt)

	if err != nil {
		return (err)
	}

	f, err = app.Firestore(ctx)
	if err != nil {
		return (err)
	}
	defer f.Close()

	s, err = discordgo.New("Bot " + cfg.Discord.BotToken)
	if err != nil {
		fmt.Printf("Invalid bot parameters: %v\n", err)
		return err
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		fmt.Printf("Cannot open the session: %v\n", err)
		return err
	}

	if cfg.Discord.GuildID != "" {
		fmt.Printf("Adding commands to guild: %s...\n", cfg.Discord.GuildID)
	} else {
		fmt.Printf("Adding commands...\n")
	}
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, cfg.Discord.GuildID, v)
		if err != nil {
			fmt.Printf("Cannot create the %v command: %v\n", v.Name, err)
			return err
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Printf("Press Ctrl+C to exit\n")
	<-stop

	if cfg.Discord.RemoveCommands {
		fmt.Printf("Removing commands...\n")
		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, cfg.Discord.GuildID, v.ID)
			if err != nil {
				fmt.Printf("Cannot delete %v command: %v\n", v.Name, err)
				return err
			}
		}
	}

	fmt.Printf("Gracefully shutdowning\n")
	return nil
}

func getPlayersFromKills(kills []Kill, getKiller bool) []string {
	players := []string{}
	keys := make(map[string]bool)

	if getKiller {
		for _, k := range kills {
			if _, value := keys[k.Killer]; !value {
				keys[k.Killer] = true
				players = append(players, k.Killer)
			}
		}
	} else {
		for _, k := range kills {
			if _, value := keys[k.Victim]; !value {
				keys[k.Victim] = true
				players = append(players, k.Victim)
			}
		}
	}

	return players
}
