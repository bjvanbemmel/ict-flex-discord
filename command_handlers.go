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

			TargetRole = optionMap["role"].RoleValue(nil, "")
			log.Info("Role has been configured.")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Registered %v as the target role.", TargetRole.Mention()),
				},
			})
		},
		"register-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			TargetChannel = optionMap["channel"].ChannelValue(nil)
			log.Info("Channel has been configured.")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: fmt.Sprintf("Registered %v as the target channel.", TargetChannel.Mention()),
				},
			})
		},
		"inspect": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var tr string = "<nil>"
			var tc string = "<nil>"

			if TargetRole != nil {
				tr = TargetRole.Mention()
			}

			if TargetChannel != nil {
				tc = TargetChannel.Mention()
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
