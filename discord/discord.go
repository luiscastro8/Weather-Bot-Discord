package discord

import (
	"Weather-Bot-Discord/api"
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/weather"
	"Weather-Bot-Discord/weather/forecast"
	"Weather-Bot-Discord/weather/points"
	"Weather-Bot-Discord/weather/zip"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var WeatherCommand = &discordgo.ApplicationCommand{
	Name:        "weather",
	Description: "get the weather",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "zip",
			Description: "get weather by zip code",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "zip-code",
					Description: "5 digit zip code",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "coordinates",
			Description: "get weather using latitude and longitude",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "latitude",
					Description: "latitude",
					Type:        discordgo.ApplicationCommandOptionNumber,
					Required:    true,
				},
				{
					Name:        "longitude",
					Description: "longitude",
					Type:        discordgo.ApplicationCommandOptionNumber,
					Required:    true,
				},
			},
		},
		{
			Name:        "address",
			Description: "get weather using an address",
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "address",
					Description: "address",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
	},
}

func WeatherHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	subCommand := i.ApplicationCommandData().Options[0]
	if subCommand.Name == "zip" {
		zipCode := subCommand.Options[0].StringValue()
		if !isValidZip(zipCode) {
			err := sendSlashCommandResponse(s, i, "Error: Zip code must be exactly 5 digits long")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		forecastUrl, ok := weather.GetUrlFromCache(zipCode)
		if !ok {
			lat, long, err := zip.GetCoordsFromZip(zipCode)
			if err != nil {
				mylogger.Errorln(err)
				err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
				if err != nil {
					mylogger.Errorln("could not send slash command message:", err)
				}
				return
			}

			weather.AcquireLockForCaching()
			forecastUrl, err = points.GetForecastURLFromCoords(lat, long)
			if err != nil {
				weather.ReleaseLockForCaching()
				mylogger.Errorln(err)
				err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
				if err != nil {
					mylogger.Errorln("could not send slash command message:", err)
				}
				return
			}
			weather.WriteToCache(zipCode, forecastUrl)
		}

		forecastMessage, err := forecast.GetForecastFromURL(forecastUrl)
		if err != nil {
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		err = sendSlashCommandResponse(s, i, forecastMessage)
		if err != nil {
			mylogger.Errorln("could not send slash command message:", err)
		} else {
			mylogger.Println("sent weather forecast for zip", zipCode)
		}
	} else if subCommand.Name == "coordinates" {
		lat := subCommand.Options[0].FloatValue()
		long := subCommand.Options[1].FloatValue()

		forecastUrl, err := points.GetForecastURLFromCoords(fmt.Sprintf("%.4f", lat), fmt.Sprintf("%.4f", long))
		if err != nil {
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		forecastMessage, err := forecast.GetForecastFromURL(forecastUrl)
		if err != nil {
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		err = sendSlashCommandResponse(s, i, forecastMessage)
		if err != nil {
			mylogger.Errorln("could not send slash command message:", err)
		} else {
			mylogger.Println("sent weather forecast for coords", lat, long)
		}
	} else if subCommand.Name == "address" {
		address := subCommand.Options[0].StringValue()
		lat, long, _, err := api.GetCoordsFromAddress(address) // todo use matching address
		if err != nil {
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		forecastUrl, err := points.GetForecastURLFromCoords(lat, long)
		if err != nil {
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		forecastMessage, err := forecast.GetForecastFromURL(forecastUrl)
		if err != nil {
			mylogger.Errorln(err)
			err = sendSlashCommandResponse(s, i, "There was an error getting the forecast")
			if err != nil {
				mylogger.Errorln("could not send slash command message:", err)
			}
			return
		}

		err = sendSlashCommandResponse(s, i, forecastMessage)
		if err != nil {
			mylogger.Errorln("could not send slash command message:", err)
		} else {
			mylogger.Println("sent weather forecast for coords", lat, long)
		}
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

func isValidZip(zipCode string) bool {
	if len(zipCode) != 5 {
		return false
	}
	for _, c := range zipCode {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
