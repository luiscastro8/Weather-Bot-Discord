package discord

import (
	"Weather-Bot-Discord/api"
	"Weather-Bot-Discord/myerrors"
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/weather/forecast"
	"Weather-Bot-Discord/weather/points"
	"github.com/bwmarrin/discordgo"
)

func addressHandler(s *discordgo.Session, i *discordgo.InteractionCreate, address string) {
	lat, long, matchedAddress, err := api.GetCoordsFromAddress(address)
	if err != nil {
		mylogger.Errorln(err)
		if addressNotFoundError, ok := err.(myerrors.AddressNotFoundError); ok {
			sendSlashCommandResponseAndLogError(s, i, "Could not find address for input: "+addressNotFoundError.UnmatchedAddress)
			return
		}
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	forecastUrl, err := points.GetForecastURLFromCoords(lat, long)
	if err != nil {
		mylogger.Errorln(err)
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	forecastMessage, err := forecast.GetForecastFromURL(forecastUrl, "Found Address: "+matchedAddress+"\n")
	if err != nil {
		mylogger.Errorln(err)
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	err = sendSlashCommandResponse(s, i, forecastMessage)
	if err != nil {
		mylogger.Errorln("could not send slash command message:", err)
	} else {
		mylogger.Println("sent weather forecast for address", matchedAddress)
	}
}
