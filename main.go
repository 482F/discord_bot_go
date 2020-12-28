//usr/bin/env go run "$(readlink $0)" $@; exit
package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
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

func readConfigJSON(filepath string) (Config, error) {
	var cs []Config
	bytes, err := ioutil.ReadFile(filepath)
	checkErr(err)
	checkErr(json.Unmarshal(bytes, &cs))
	return cs[0], nil
}
