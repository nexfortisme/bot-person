package util

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleErrors(err error) {
	if err != nil {
		fmt.Println("Error: " + err.Error())
	}
}

func HandleFatalErrors(err error, message string) {
	if err != nil {
		log.Fatalf(message)
	}
}

func ReplaceIDsWithNames(m *discordgo.MessageCreate, s *discordgo.Session) string {
	id := s.State.User.ID
	toReplace := fmt.Sprintf("<@%s> ", id)
	msg := strings.Replace(m.Message.Content, toReplace, "", 1)
	msg = replaceMentionsWithNames(m.Mentions, msg)

	return msg
}

// The message string that the bot receives reads mentions of other users as
// an ID in the form of "<@000000000000>", instead iterate over each mention and
// replace the ID with the user's username
func replaceMentionsWithNames(mentions []*discordgo.User, message string) string {
	retStr := strings.Clone(message)
	for _, mention := range mentions {
		idStr := fmt.Sprintf("<@%s>", mention.ID)
		retStr = strings.ReplaceAll(retStr, idStr, mention.Username)
	}
	return retStr
}

func LowerFloatPrecision(num float64) float64 {
	floatString := fmt.Sprintf("%.2f", num)
	returnFloatValue, _ := strconv.ParseFloat(floatString, 64)
	return returnFloatValue
}

func IntToFloat(num int) float64 {
	return LowerFloatPrecision(float64(num))
}

func GetOofResponse() string {
	options := [4]string{"oof.", "That sucks.", "Should have rolled better.", "I saw a person roll a better number once. It wasn't you but I saw someone else do it. \n"}

	return options[rand.Intn(len(options))]
}

func GetBadBotResponse() string {
	badBotResponses := make([]string, 0)
	badBotResponses = append(badBotResponses, "I'm sorry")
	badBotResponses = append(badBotResponses, "It won't happen again")
	badBotResponses = append(badBotResponses, "Eat Shit")
	badBotResponses = append(badBotResponses, "Ok.")
	badBotResponses = append(badBotResponses, "Sure Thing.")
	badBotResponses = append(badBotResponses, "Like you are the most perfect being in existance. Pound sand pal.")
	badBotResponses = append(badBotResponses, "https://youtu.be/4X7q87RDSHI")

	return badBotResponses[rand.Intn(len(badBotResponses))]
}

func CleanUpImages(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// time.Sleep(time.Hour * 8);
	// i.Interaction.Intre

	time.AfterFunc(time.Hour*8, func() {
		s.InteractionResponseDelete(i.Interaction)
	})
}
