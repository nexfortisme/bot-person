package util

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

func GetGoodBotResponse() string {
	goodBotResponses := make([]string, 0)
	goodBotResponses = append(goodBotResponses, "Thank you, I'm here to serve.")
	goodBotResponses = append(goodBotResponses, "I'm glad I could assist you.")
	goodBotResponses = append(goodBotResponses, "Your satisfaction is my top priority.")
	goodBotResponses = append(goodBotResponses, "Thanks for the compliment, it means a lot.")
	goodBotResponses = append(goodBotResponses, "It's my pleasure to be of help.")
	goodBotResponses = append(goodBotResponses, "I strive to be the best bot I can be.")
	goodBotResponses = append(goodBotResponses, "I appreciate the positive feedback.")
	goodBotResponses = append(goodBotResponses, "I'm happy to hear that you found my service helpful.")
	goodBotResponses = append(goodBotResponses, "Thanks for acknowledging my hard work.")
	goodBotResponses = append(goodBotResponses, "I'm always here if you need me.")

	return goodBotResponses[rand.Intn(len(goodBotResponses))]
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

//func GetStreakBonus(userStats persistance.UserStatsStruct) (string, float64, error) {
//
//	var returnString string
//	var modifier int
//
//	userStats.BonusStreak++
//	streak := userStats.BonusStreak
//
//	returnString, modifier = GetStreakStringAndModifier(streak)
//
//	// Setting random seed and generating a, value safe, token amount
//	randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))
//	reward := randomizer.Intn(45) + 5
//	reward *= modifier
//	rewardf64 := float64(reward) / 10.0
//	finalReward := LowerFloatPrecision(rewardf64)
//
//	// Updating User Record
//	userStats.LastBonus = time.Now()
//	userStats.ImageTokens += finalReward
//
//	return "", -1, nil
//}

func GetStreakStringAndModifier(streak int) (string, int) {

	var returnString string
	var modifier int

	if streak%10 == 0 && streak%100 != 0 && streak%50 != 0 {
		returnString = fmt.Sprintf("Congrats on keeping the streak alive. Current Streak: %d. Bonus Modifier: 2x", streak)
		modifier = 2
	} else if streak%25 == 0 && streak%50 != 0 && streak%100 != 0 {
		returnString = fmt.Sprintf("Great work on keeping the streak alive! Current Streak: %d. Bonus Modifier: 5x", streak)
		modifier = 5
	} else if streak%50 == 0 && streak%100 != 0 {
		returnString = fmt.Sprintf("Wow! That's a long time. Current Streak: %d. Bonus Modifier: 10x", streak)
		modifier = 10
	} else if streak%69 == 0 {
		returnString = fmt.Sprintf("Nice, Congratulations! Current Streak: %d. Bonus Modifier: 15x", streak)
		modifier = 15
	} else if streak%100 == 0 {
		returnString = fmt.Sprintf("Few people ever reach is this far, Congratulations! Current Streak: %d. Bonus Modifier: 50x", streak)
		modifier = 50
	} else {
		returnString = fmt.Sprintf("Current Bonus Streak: %d", streak)
		modifier = 1
	}

	return returnString, modifier
}

func GetUserBonus(min int, max int, modifier int) float64 {
	randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))
	reward := randomizer.Intn(max-min) + min
	reward *= modifier
	rewardF64 := float64(reward) / 10.0
	finalReward := LowerFloatPrecision(rewardF64)
	return finalReward
}
