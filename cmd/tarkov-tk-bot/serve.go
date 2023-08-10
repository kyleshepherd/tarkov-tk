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
		log.Info().Msgf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	err = s.Open()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, cfg.Discord.GuildID, commands)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	allCmds, _ := s.ApplicationCommands(s.State.User.ID, cfg.Discord.GuildID)

	if cfg.Discord.GuildID != "" {
		log.Info().Msgf("Commands added to guild: %s...\n", cfg.Discord.GuildID)
	} else {
		log.Info().Msg("Commands added...\n")
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	sig := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, os.Interrupt, os.Kill}
	signal.Notify(stop, sig...)
	log.Info().Msg("Press Ctrl+C to exit\n")
	<-stop

	if cfg.Discord.RemoveCommands {
		log.Info().Msg("Removing commands...\n")
		for _, v := range allCmds {
			err := s.ApplicationCommandDelete(s.State.User.ID, cfg.Discord.GuildID, v.ID)
			if err != nil {
				log.Error().Err(err)
				return err
			}
		}
	}

	log.Info().Msg("Gracefully shutdowning\n")
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
	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"removeKill": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx := context.Background()
			kills, err := store.ListKillsForServer(ctx, i.GuildID)
			if err != nil {
				log.Error().Err(err)
				return
			}

			if len(kills) < 1 {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "No kills to remove",
					},
				})

				if err != nil {
					log.Error().Err(err)
					return
				}
				return
			}

			kill := kills[0]

			err = store.DeleteKill(ctx, kill.ID)
			if err != nil {
				log.Error().Err(err)
				return
			}

			err = s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			if err != nil {
				log.Error().Err(err)
				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "**Kill removed**",
				},
			})

			if err != nil {
				log.Error().Err(err)
				return
			}
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"tklog": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ctx := context.Background()
			killer := i.ApplicationCommandData().Options[0].UserValue(s)
			victim := i.ApplicationCommandData().Options[1].UserValue(s)
			reason := ""

			if len(i.ApplicationCommandData().Options) > 2 {
				reason = i.ApplicationCommandData().Options[2].StringValue()
			}

			if len(reason) > 500 {
				reason = reason[:500]
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
				log.Error().Err(err)
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
			var shouldCsv bool
			for _, o := range i.ApplicationCommandData().Options {
				switch o.Type {
				case discordgo.ApplicationCommandOptionUser:
					killer = o.UserValue(s)
				case discordgo.ApplicationCommandOptionBoolean:
					shouldCsv = o.BoolValue()
				}
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

			guild, _ := s.Guild(i.GuildID)

			if shouldCsv {
				csvBuffer := &bytes.Buffer{}
				err = gocsv.Marshal(kills, csvBuffer)

				if err != nil {
					log.Error().Err(err)
					return
				}

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

				if err != nil {
					log.Error().Err(err)
				}
			} else {
				response := fmt.Sprintf("**Stats for %s**\n\n", guild.Name)

				hasSentIResponse := false

				for _, k := range kills {
					date := fmt.Sprintf("%d/%d/%d", k.Date.Day(), k.Date.Month(), k.Date.Year())
					kMsg := fmt.Sprintf("%s - **%s** killed **%s**", date, k.Killer, k.Victim)
					if k.Reason != "" {
						kMsg += fmt.Sprintf(": \"%s\"", k.Reason)
					}
					kMsg += fmt.Sprint("\n\n")
					if len(response+kMsg) > 2000 {
						if !hasSentIResponse {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content:  response,
									CustomID: "test",
								},
							})
							hasSentIResponse = true
						} else {
							s.ChannelMessageSend(i.ChannelID, response)
						}
						response = "\n"
					}
					response += kMsg
				}

				if !hasSentIResponse {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: response,
						},
					})
				} else {
					s.ChannelMessageSend(i.ChannelID, response)
				}
			}
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
			ctx := context.Background()
			kills, err := store.ListKillsForServer(ctx, i.GuildID)
			if err != nil {
				log.Error().Err(err)
				s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("Could not get kills for server. Please try again\n"))
				return
			}

			if len(kills) == 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf("No kills currently logged"),
					},
				})
				return
			}

			users, err := s.GuildMembers(i.GuildID, "", 1000)
			kill := kills[0]

			for _, u := range users {
				if kill.Killer == u.User.ID {
					if u.Nick != "" {
						kill.Killer = u.Nick
					} else {
						kill.Killer = u.User.Username
					}
				}

				if kill.Victim == u.User.ID {
					if u.Nick != "" {
						kill.Victim = u.Nick
					} else {
						kill.Victim = u.User.Username
					}
				}
			}

			kMsg := fmt.Sprintf("**%s** killed **%s**", kill.Killer, kill.Victim)

			if kill.Reason != "" {
				kMsg += fmt.Sprintf(": \"%s\"", kill.Reason)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: kMsg,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🗑️",
									},
									Label:    "Remove Kill",
									Style:    discordgo.DangerButton,
									CustomID: "removeKill",
								},
							},
						},
					},
				},
			})
		},
		"tkthanks": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("**Thank you to everyone who uses this bot but especially to these Patreon supporters, past and present:**\nThomas Pimentel\nDale Schlegel\nWill Eudave\nJasmine\nGrndControl\nBlake Freeman\nMatryx PX\nChris James\nDenver Welch\nJared Heisler\nRyan Miers\nArch Andrews\nNorman Golden\nChris Laquidara\nG\nKyle Melton\nMichael Li\nCoop Diddy\nPBRbiter"),
				},
			})
		},
	}
)
