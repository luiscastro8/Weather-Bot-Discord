package main

import (
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/token"
	"Weather-Bot-Discord/weather"
	"Weather-Bot-Discord/weather/forecast"
	"Weather-Bot-Discord/weather/points"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var Logger *mylogger.MyLogger

func main() {
	Logger = mylogger.New()

	botToken, err := token.GetToken(Logger)
	if err != nil {
		Logger.Fatalln("Unable to get bot token:", err)
	}

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		Logger.Fatalln("Unable to create discord session:", err)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		Logger.Fatalln("Unable to open discord connection:", err)
	}

	err = weather.OpenZipFile("zip-codes.csv")
	if err != nil {
		Logger.Fatalln("Unable to cache zip codes and coordinates:", err)
	}

	Logger.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = dg.Close()
	if err != nil {
		Logger.Fatalln("Unable to close discord connection:", err)
	}
	Logger.Println("Closed discord connection")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	words := strings.Fields(m.Content)
	if len(words) != 2 {
		return
	}
	if words[0] != "!weather" {
		return
	}
	if len(words[1]) != 5 {
		return
	}
	for _, c := range words[1] {
		if c < '0' || c > '9' {
			return
		}
	}

	lat, long, err := weather.GetCoordsFromZip(words[1])
	if err != nil {
		Logger.Errorln(err)
		_, _ = s.ChannelMessageSend(m.ChannelID, "There was an error getting the forecast")
		return
	}

	url, err := points.GetForecastURLFromCoords(lat, long)
	if err != nil {
		Logger.Errorln(err)
		_, _ = s.ChannelMessageSend(m.ChannelID, "There was an error getting the forecast")
		return
	}

	forecastMessage, err := forecast.GetForecastFromURL(url)
	if err != nil {
		Logger.Errorln(err)
		_, _ = s.ChannelMessageSend(m.ChannelID, "There was an error getting the forecast")
		return
	}

	_, err = s.ChannelMessageSend(m.ChannelID, forecastMessage)
	if err != nil {
		Logger.Errorln("error sending message", err)
		_, _ = s.ChannelMessageSend(m.ChannelID, "There was an error getting the forecast")
		return
	}
	Logger.Println(fmt.Sprintf("Sent message for zip %s in channel %s for guild %s", words[1], m.ChannelID, m.GuildID))
}
