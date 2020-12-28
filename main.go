//usr/bin/env go run "$(readlink $0)" $@; exit
package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	BotID       string            `json:"BotID"`
	GuildID     string            `json:"GuildID"`
	Token       string            `json:"Token"`
	IgnoreUsers []string          `json:"IgnoreUsers"`
	UserMap     map[string]string `json:"UserMap"`
}

var config Config
var presences map[string]*discordgo.Presence = map[string]*discordgo.Presence{}
var updateIndex map[string]int = map[string]int{}
var testChannelID string

func main() {
	var err error
	checkErr(err)
	var stopBot <-chan bool = make(chan bool)
	config, err = readConfigJSON("./config.json")
	checkErr(err)

	discord, err := discordgo.New()
	discord.Token = "Bot " + config.Token
	checkErr(err)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildPresences | discordgo.IntentsGuildMessages)
	discord.AddHandler(onMessageCreate)
	discord.AddHandler(onPresenceUpdate)
	err = discord.Open()
	checkErr(err)
	defer discord.Close()

	testChannelID = getChannelByName(discord, "test").ID

	fmt.Println("Listening...")
	getChannelByName(discord, "General")
	<-stopBot //プログラムが終了しないようロック
	return
}

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

func generalCommand(s *discordgo.Session, channel *discordgo.Channel, commands []string) bool {
	return true
}

func terrariaCommand(s *discordgo.Session, channel *discordgo.Channel, commands []string) bool {
	return true
}

func minecraftCommand(s *discordgo.Session, channel *discordgo.Channel, commands []string) bool {
	return true
}

func testCommand(s *discordgo.Session, channel *discordgo.Channel, commands []string) bool {
	for _, command := range commands {
		sendMessage(s, channel.ID, fmt.Sprint(command))
	}
	return true
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

func arrayHasString(ss []string, target string) bool {
	for _, s := range ss {
		if s == target {
			return true
		}
	}
	return false
}

func readTxtFile(filepath string) (string, error) {
	bt, err := ioutil.ReadFile(filepath)
	var str string = string(bt)
	str = str[:len(str)-1] //改行が入ってしまうので取り除く
	return str, err
}

func readConfigJSON(filepath string) (Config, error) {
	var cs []Config
	bytes, err := ioutil.ReadFile(filepath)
	checkErr(err)
	checkErr(json.Unmarshal(bytes, &cs))
	return cs[0], nil
}

func checkErr(e error) bool {
	if e != nil {
		fmt.Println(os.Stderr, e)
		os.Exit(1)
	}
	return true
}

func sleep(ms time.Duration) {
	time.Sleep(time.Millisecond * ms)
	return
}
