package main

import (
	"github.com/bwmarrin/discordgo"
    "log"
)

func sendMessage(s *discordgo.Session, channelID string, msg string) {
	if msg == "" {
		return
	}
	_, err := s.ChannelMessageSend(channelID, msg)

	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func getChannelByName(s *discordgo.Session, name string) *discordgo.Channel {
	channels, err := s.GuildChannels(config.GuildID)
	checkErr(err)
	for _, channel := range channels {
		if channel.Name == name {
			return channel
		}
	}
	return nil
}
