package main

import "github.com/bwmarrin/discordgo"

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "tklog",
		Description: "Log a TK",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "killer",
				Description: "Killer",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "victim",
				Description: "Victim",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "reason",
				Description: "Reason for TK (Maximum 500 characters)",
				Required:    false,
			},
		},
	},
	{
		Name:        "tkkills",
		Description: "Get kills leaderboard for server",
	},
	{
		Name:        "tkdeaths",
		Description: "Get deaths leaderboard for server",
	},
	{
		Name:        "tkstats",
		Description: "Get all TKs for server or a single player in server exported into a CSV",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "killer",
				Description: "User whose kills to retrieve",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "csv",
				Description: "Generate CSV file of stats",
				Required:    false,
			},
		},
	},
	{
		Name:        "tkreset",
		Description: "Reset the TK server data for your channel. **THIS WILL DELETE ALL TK LOGS**",
	},
	{
		Name:        "tkinfo",
		Description: "Some info about the project and the creator, Kyle",
	},
	{
		Name:        "tkremove",
		Description: "Removes the last logged kill",
	},
	{
		Name:        "tkthanks",
		Description: "Shout out to my Patreons!",
	},
}
