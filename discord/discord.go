package discord

import (
	"Weather-Bot-Discord/mylogger"
	"github.com/bwmarrin/discordgo"
)

func WeatherHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	frequency := i.ApplicationCommandData().Options[0]
	hourly := false
	if frequency.Name == "hourly" {
		hourly = true
	}

	subCommand := frequency.Options[0]
	if subCommand.Name == "zip" {
		zipCode := subCommand.Options[0].StringValue()
		zipHandler(s, i, zipCode, hourly)
	} else if subCommand.Name == "coordinates" {
		lat := subCommand.Options[0].FloatValue()
		long := subCommand.Options[1].FloatValue()
		coordinatesHandler(s, i, lat, long, hourly)
	} else if subCommand.Name == "address" {
		address := subCommand.Options[0].StringValue()
		addressHandler(s, i, address, hourly)
	}
}

func sendSlashCommandResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	return err
}

func sendSlashCommandResponseAndLogError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := sendSlashCommandResponse(s, i, message)
	if err != nil {
		mylogger.Errorln("could not send slash command message:", err)
	}
}
