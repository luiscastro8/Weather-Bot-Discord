package discord

import (
	addressApi "Weather-Bot-Discord/api/address"
	"Weather-Bot-Discord/api/weather/forecast"
	"Weather-Bot-Discord/api/weather/points"
	"Weather-Bot-Discord/myerrors"
	"Weather-Bot-Discord/mylogger"
	"github.com/bwmarrin/discordgo"
)

func addressHandler(s *discordgo.Session, i *discordgo.InteractionCreate, address string, hourly bool) {
	lat, long, matchedAddress, err := addressApi.GetCoordsFromAddress(address)
	if err != nil {
		mylogger.Errorln(err)
		if addressNotFoundError, ok := err.(myerrors.AddressNotFoundError); ok {
			sendSlashCommandResponseAndLogError(s, i, "Could not find address for input: "+addressNotFoundError.UnmatchedAddress)
			return
		}
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	forecastUrl, err := points.GetDailyForecastURLFromCoords(lat, long)
	if err != nil {
		mylogger.Errorln(err)
		sendSlashCommandResponseAndLogError(s, i, "There was an error getting the forecast")
		return
	}

	forecastMessage, err := forecast.GetForecastFromURL(forecastUrl, "Found Address: "+matchedAddress+"\n", hourly)
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
