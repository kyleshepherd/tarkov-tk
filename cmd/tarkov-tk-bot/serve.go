package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gocarina/gocsv"
	"github.com/kyleshepherd/discord-tk-bot/internal/storage"
	"github.com/kyleshepherd/discord-tk-bot/internal/storage/firestore"
	"github.com/spf13/cobra"
)

var s *discordgo.Session
var store storage.KillStore

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

func listen(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handler() error {
	http.HandleFunc("/", listen)
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		return err
	}
	return nil
}

func serveRun(cmd *cobra.Command, args []string) error {
	go handler()
	ctx := context.Background()
	ks, err := firestore.NewKillStore(ctx, cfg.Firebase.ProjectID, cfg.Firebase.ServiceAccountFilePath)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	store = ks
	defer store.Close()
	s, err := discordgo.New("Bot " + cfg.Discord.BotToken)
	if err != nil {
		log.Error().Err(err)
		return err
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	cmds, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, cfg.Discord.GuildID, commands)
	if cfg.Discord.GuildID != "" {
		log.Printf("Commands added to guild: %s...\n", cfg.Discord.GuildID)
	} else {
		log.Printf("Commands added...\n")
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	sig := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, os.Interrupt, os.Kill}
	signal.Notify(stop, sig...)
	fmt.Printf("Press Ctrl+C to exit\n")
	<-stop

	if cfg.Discord.RemoveCommands {
		log.Printf("Removing commands...\n")
		for _, v := range cmds {
			err := s.ApplicationCommandDelete(s.State.User.ID, cfg.Discord.GuildID, v.ID)
			if err != nil {
				log.Error().Err(err)
				return err
			}
		}
	}

	log.Printf("Gracefully shutdowning\n")
	return nil
}

func getPlayersFromKills(kills []*storage.Kill, getKiller bool) []string {
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

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"tklog": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx := context.Background()
			killer := i.ApplicationCommandData().Options[0].UserValue(s)
			victim := i.ApplicationCommandData().Options[1].UserValue(s)
			reason := ""

			if len(i.ApplicationCommandData().Options) > 2 {
				reason = i.ApplicationCommandData().Options[2].StringValue()
			}

			kill := storage.Kill{
				ServerID: i.GuildID,
				Killer:   killer.ID,
				Victim:   victim.ID,
				Reason:   reason,
				Date:     time.Now(),
			}

			_, err := store.CreateKill(ctx, &kill)
			if err != nil {
				log.Error().Err(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Could not log kill by **%s** on **%s**. Please try again\n", killer.Username, victim.Username),
					},
				})
				return
			}

			msgContent := fmt.Sprintf("Kill by **%s** on **%s** logged", killer.Username, victim.Username)
			if kill.Reason != "" {
				msgContent += fmt.Sprintf(": \"**%s**\"", kill.Reason)
			}
			msgContent += fmt.Sprintf("\n")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msgContent,
				},
			})
		},
		"tkkills": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx := context.Background()
			kills, err := store.ListKillsForServer(ctx, i.GuildID)
			users, err := s.GuildMembers(i.GuildID, "", 1000)

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

			playerKills := []PlayerCount{}
			players := getPlayersFromKills(kills, true)

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
				for _, u := range users {
					if k.Player == u.User.ID {
						if u.Nick != "" {
							playerName = u.Nick
						} else {
							playerName = u.User.Username
						}
					}
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
			kills, err := store.ListKillsForServer(ctx, i.GuildID)
			users, err := s.GuildMembers(i.GuildID, "", 1000)
			if err != nil {
				log.Error().Err(err)
				s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Could not get kills for server. Please try again\n"))
				return
			}

			if err != nil {
				log.Error().Err(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Could not get kills for server. Please try again\n"),
					},
				})
				return
			}

			playerDeaths := []PlayerCount{}
			players := getPlayersFromKills(kills, false)

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
				for _, u := range users {
					if k.Player == u.User.ID {
						if u.Nick != "" {
							playerName = u.Nick
						} else {
							playerName = u.User.Username
						}
					}
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
			ctx := context.Background()
			users, err := s.GuildMembers(i.GuildID, "", 1000)
			if err != nil {
				log.Error().Err(err)
				s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Could not get kills for server. Please try again\n"))
				return
			}
			var killer *discordgo.User
			if len(i.ApplicationCommandData().Options) > 0 {
				killer = i.ApplicationCommandData().Options[0].UserValue(s)
			}

			var kills []*storage.Kill

			if killer != nil {
				kills, err = store.ListPlayerKillsForServer(ctx, i.GuildID, killer.ID)
				if err != nil {
					log.Error().Err(err)
					s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Could not get kills for server. Please try again\n"))
					return
				}
			} else {
				kills, err = store.ListKillsForServer(ctx, i.GuildID)
				if err != nil {
					log.Error().Err(err)
					s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Could not get kills for server. Please try again\n"))
					return
				}
			}

			for _, k := range kills {
				for _, u := range users {
					if k.Killer == u.User.ID {
						if u.Nick != "" {
							k.Killer = u.Nick
						} else {
							k.Killer = u.User.Username
						}
					}

					if k.Victim == u.User.ID {
						if u.Nick != "" {
							k.Victim = u.Nick
						} else {
							k.Victim = u.User.Username
						}
					}
				}
			}

			csvBuffer := &bytes.Buffer{}
			err = gocsv.Marshal(kills, csvBuffer)

			if err != nil {
				log.Error().Err(err)
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
			err := store.DeleteKillsForServer(ctx, i.GuildID)
			if err != nil {
				log.Error().Err(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("Could not reset kills for server. Please try again\n"),
					},
				})
				return
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
					Content: fmt.Sprintf("TO DO"),
				},
			})
		},
	}
)
