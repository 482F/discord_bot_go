package main

import (
	"github.com/bwmarrin/discordgo"
    "strings"
    "fmt"
)

//メッセージが投稿されたら
func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	if mc.Author.ID == config.BotID {
		return
	}
	var messages []string = strings.Split(mc.Content, " ")
	if messages[0] != "<@!"+config.BotID+">" {
		return
	}
	var commands []string = messages[1:]
	var channel *discordgo.Channel
	var channelName string
	channel, err := s.Channel(mc.ChannelID)
	checkErr(err)
	channelName = channel.Name
	switch channelName {
	case "general":
		generalCommand(s, channel, commands)
	case "terraria":
		terrariaCommand(s, channel, commands)
	case "minecraft":
		minecraftCommand(s, channel, commands)
	case "test":
		testCommand(s, channel, commands)
	default:
		sendMessage(s, channel.ID, fmt.Sprint(commands))
	}
	return
}

//メンバーのステータスが変更されたら
func onPresenceUpdate(s *discordgo.Session, pu *discordgo.PresenceUpdate) {
	var selfUpdateIndex int
	var message string
	var before *discordgo.Presence
	var after *discordgo.Presence
	before = presences[pu.User.ID]
	after = &pu.Presence
	presences[pu.User.ID] = after
	fmt.Println("after: ", after)
	fmt.Println("Username: ", after.User.Username)
	fmt.Println("Nick: ", after.Nick)
	fmt.Println("ID: ", after.User.ID)
	if before == nil {
		return
	}

	if arrayHasString(config.IgnoreUsers, before.User.ID) {
		return
	}

	var name string = config.UserMap[after.User.ID]
	if before.Status != after.Status {
		message = name + " is " + string(after.Status)
	} else if before.Game != after.Game {
		var gameName string = "None"
		if after.Game != nil {
			gameName = after.Game.Name
		}
		message = name + " playing " + gameName
	} else {
		return
	}
	if _, err := updateIndex[after.User.ID]; !err {
		updateIndex[after.User.ID] = 0
	}
	updateIndex[after.User.ID] += 1
	selfUpdateIndex = updateIndex[after.User.ID]
	sleep(5000)
	if selfUpdateIndex != updateIndex[after.User.ID] {
		return
	}

	sendMessage(s, testChannelID, message)
	updateIndex[after.User.ID] = 0
	return
}
