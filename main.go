package main

import (
	"Weather-Bot-Discord/discord"
	"Weather-Bot-Discord/mylogger"
	"Weather-Bot-Discord/token"
	"Weather-Bot-Discord/weather"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
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

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.ApplicationCommandData().Name == "weather" {
			discord.WeatherHandler(s, i)
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		Logger.Fatalln("Unable to open discord connection:", err)
	}

	_, err = dg.ApplicationCommandCreate(dg.State.User.ID, "", discord.WeatherCommand)
	if err != nil {
		Logger.Fatalln("Cannot create command", err)
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
