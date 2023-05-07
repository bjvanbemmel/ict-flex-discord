package main

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "register-channel",
			Description: "Register a channel that you would like the announcements to be posted in",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "The target channel",
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
						discordgo.ChannelTypeGuildNews,
					},
					Required: true,
				},
			},
		},
		{
			Name:        "register-role",
			Description: "Register the role you want to be pinged to notify users about new announcements",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "The target role",
					Required:    true,
				},
			},
		},
		{
			Name:        "inspect",
			Description: "Shows the current configuration",
		},
		{
			Name:        "embed",
			Description: "Sends an embed... I think..?",
		},
	}
)
