package util

import (
	"fmt"
	"log"
	"strings"

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

func ReplaceIDsWithNames(m *discordgo.MessageCreate, s *discordgo.Session) string{
	id := s.State.User.ID;
	toReplace := fmt.Sprintf("<@%s> ", id)
	msg := strings.Replace(m.Message.Content, toReplace, "", 1)
	msg = replaceMentionsWithNames(m.Mentions, msg)

	return msg;
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
