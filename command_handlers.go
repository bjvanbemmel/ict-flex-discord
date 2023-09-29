package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"register-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if err := Settings.SetRole(optionMap["role"].RoleValue(nil, "")); err != nil {
				log.Errorf("Something went wrong while setting role: `%s`", err.Error())
				return
			} else {
				log.Info("Role has been configured.")
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Registered %v as the target role.", Settings.TargetRole.Mention()),
				},
			})
		},
		"register-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if err := Settings.SetChannel(optionMap["channel"].ChannelValue(nil)); err != nil {
				log.Errorf("Something went wrong while setting channel: `%s`", err.Error())
				return
			} else {
				log.Info("Channel has been configured.")
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Registered %v as the target channel.", Settings.TargetChannel.Mention()),
				},
			})
		},
		"inspect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var tr string = "<nil>"
			var tc string = "<nil>"

			role, err := Settings.Role()
			if err != nil {
				log.Errorf("Something went wrong while fetching role: `%s`", err.Error())
			}

			channel, err := Settings.Channel()
			if err != nil {
				log.Errorf("Something went wrong while fetching channel: `%s`", err.Error())
			}

			if role != nil {
				tr = role.Mention()
			}

			if channel != nil {
				tc = channel.Mention()
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
					Embeds: []*discordgo.MessageEmbed{
						{
							Type:        discordgo.EmbedTypeRich,
							Title:       "Configuration",
							Description: "This is the current configuration:",
							Color:       0xEF7E05,
							Fields: []*discordgo.MessageEmbedField{
								{
									Name:  "Target channel:",
									Value: tc,
								},
								{
									Name:  "Target role:",
									Value: tr,
								},
							},
						},
					},
				},
			})
		},
	}
)
