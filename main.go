package main

import (
	"Weather-Bot-Discord/mylogger"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

var Logger *mylogger.MyLogger

func getToken() (string, error) {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		return "", errors.New("DISCORD_TOKEN environment variable is empty")
	}
	return token, nil
}

func main() {
	Logger = mylogger.New()

	token, err := getToken()
	if err != nil {
		Logger.Fatalln("Unable to get bot token:", err)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		Logger.Fatalln("Unable to create discord session:", err)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		Logger.Fatalln("Unable to open discord connection:", err)
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

	var err error
	var message *discordgo.Message
	if m.Content == "ping" {
		message, err = s.ChannelMessageSend(m.ChannelID, "pong")
	} else if m.Content == "pong" {
		message, err = s.ChannelMessageSend(m.ChannelID, "ping")
	}

	if err != nil {
		Logger.Errorln("error sending message", err)
	}
	Logger.Println(fmt.Sprintf("Sent message \"%s\" in channel %s for guild %s", message.Content, message.ChannelID, m.GuildID))
}
