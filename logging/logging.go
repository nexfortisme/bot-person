package logging

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/util"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	botTracking BotTracking
)

func InitBotStatistics() {
	var trackingFile []byte

	trackingFile, err := ioutil.ReadFile("botTracking.json")
	if err != nil {

		log.Printf("Error Reading botTracking. Creating File")
		os.WriteFile("botTracking.json", []byte("{\"BadBotCount\":0,\"MessageCount\":0, \"UserStats\":[]}"), 0666)

		trackingFile, err = ioutil.ReadFile("botTracking.json")
		if err != nil {
			log.Fatalf("Could not read config file: botTracking.json")
		}

	}

	err = json.Unmarshal(trackingFile, &botTracking)
	if err != nil {
		log.Fatalf("Could not parse: botTracking.json")
	}
}

func LogOutGoingMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {

}

func LogError(err string) {
	log.Fatalf(err)
}

func LogIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name
	message := util.ReplaceIDsWithNames(m, s)

	log.Printf("%s (%s) < %s\n", requestUser, rGuildName, message)
}

func IncrementTracker(flag int, userId string, username string) {
	foundUser := false
	foundUser = handleUserStatIncrementing(flag, userId)

	if !foundUser {
		createNewUserTracking(userId, username, flag)
	}

	incrementBotTracking(flag)
}

func GetUserStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, element := range botTracking.UserStats {
		if element.UserId != m.Author.ID {
			continue
		} else {
			msg := "You have interacted with the bot " + strconv.Itoa(element.UserStats.MessageCount) + " times, praised the bot " + strconv.Itoa(element.UserStats.GoodBotCount) + " times, and scolded the bot " + strconv.Itoa(element.UserStats.BadBotCount) + " times. You have requested an image " + strconv.Itoa(element.UserStats.ImageCount) + " times."
			_, err := s.ChannelMessageSend(m.ChannelID, msg)
			util.HandleErrors(err)
			LogOutGoingMessage(s, m, msg)
			return
		}
	}

	msg := "Sorry, you don't have any recorded interactions with the bot."
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	util.HandleErrors(err)
	LogOutGoingMessage(s, m, msg)
}

func SlashGetUserStats(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	userId := i.Interaction.Member.User.ID

	for _, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {
			msg := "You have interacted with the bot " + strconv.Itoa(element.UserStats.MessageCount) + " times, praised the bot " + strconv.Itoa(element.UserStats.GoodBotCount) + " times, and scolded the bot " + strconv.Itoa(element.UserStats.BadBotCount) + " times. You have requested an image " + strconv.Itoa(element.UserStats.ImageCount) + " times."
			return msg
		}
	}
	msg := "Sorry, you don't have any recorded interactions with the bot."

	return msg
}

func GetBotStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildCount := len(s.State.Guilds)
	msg := "Across " + strconv.Itoa(guildCount) + " servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times, praised " + strconv.Itoa(botTracking.GoodBotCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	util.HandleErrors(err)
	LogOutGoingMessage(s, m, msg)
}

func SlashGetBotStats(s *discordgo.Session) string {
	guildCount := len(s.State.Guilds)
	msg := "Across " + strconv.Itoa(guildCount) + " servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times, praised " + strconv.Itoa(botTracking.GoodBotCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	return msg
}

func GetBadBotCount() int {
	return botTracking.BadBotCount
}

func ShutDown() {
	log.Println("Writing botTracking.json...")
	fle, _ := json.Marshal(botTracking)
	os.WriteFile("botTracking.json", fle, 0666)
}

func handleUserStatIncrementing(flag int, userId string) bool {

	// TODO - Handle this better. I don't like traversing an array each time.
	// Convert to a map....eventually
	for index, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {

			if flag == 1 {
				element.UserStats.MessageCount++
				element.UserStats.GoodBotCount++
			} else if flag == 2 {
				element.UserStats.MessageCount++
				element.UserStats.BadBotCount++
			} else if flag == 3 {
				element.UserStats.MessageCount++
				element.UserStats.ImageCount++
			} else {
				element.UserStats.MessageCount++
			}

			// Is this necessary?
			botTracking.UserStats[index] = element
			return true
		}
	}

	return false
}

func incrementBotTracking(flag int) {
	if flag == 1 {
		botTracking.MessageCount++
		botTracking.GoodBotCount++
	} else if flag == 2 {
		botTracking.MessageCount++
		botTracking.BadBotCount++
	} else {
		botTracking.MessageCount++
	}
}

func createNewUserTracking(userId string, username string, flag int) {
	log.Println("Creating New User For: " + username)
	if flag == 1 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 1, 0, 0, 0}})
	} else if flag == 2 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 0, 1, 0, 0}})
	} else if flag == 3 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 0, 1, 1, 0}})
	} else {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 0, 0, 0, 0}})
	}
}

func getUser(userId string) (UserStatsStruct, error) {
	for _, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {
			return element.UserStats, nil
		}
	}
	return UserStatsStruct{-1, -1, -1, -1, -1}, nil
}

func updateUser(userStats UserStatsStruct, userId string) bool {
	for index, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {

			// Is this necessary?
			botTracking.UserStats[index].UserStats = userStats
			return true
		}
	}
	return false
}

func AddImageTokens(tokenAmount int, userId string) bool {

	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
		return true
	} else {
		user.ImageTokens += tokenAmount
		return updateUser(user, userId)
	}

}

func TransferrImageTokens(tokenAmount int, fromUserId string, toUserId string) bool {

	fromUser, _ := getUser(fromUserId)
	toUser, _ := getUser(toUserId)

	// Checking to see if user exists
	if fromUser.MessageCount == -1 {
		return false
	} else {

		// The User exists but is trying to transferr more tokens then they have
		if fromUser.ImageTokens < tokenAmount {
			return false
		}

		// Checking to see if the toUser exists
		if toUser.MessageCount == -1 {

			// toUser doesn't exist
			// Creates user and assigns them the number of tokens that is being transferred
			botTracking.UserStats = append(botTracking.UserStats, UserStruct{toUserId, UserStatsStruct{0, 0, 0, 0, tokenAmount}})
			fromUser.ImageTokens -= tokenAmount
			return updateUser(fromUser, fromUserId)
		} else {
			toUser.ImageTokens += tokenAmount
			fromUser.ImageTokens -= tokenAmount
			return updateUser(toUser, toUserId) && updateUser(fromUser, fromUserId)
		}
	}

}

func UseImageToken(userId string) bool {

	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		return false
	} else {
		if user.ImageTokens <= 0 {
			return false
		} else {
			user.ImageTokens--
			return updateUser(user, userId);
		}
	}

}

func UserHasTokens(userId string) bool {
	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		return false
	} else {
		if user.ImageTokens <= 0 {
			return false
		} else {
			return true
		}
	}

}

func GetUserTokenCount(userId string) int {
	user, _ := getUser(userId)

	if user.MessageCount == -1 {
		return 0
	} else {
		return user.ImageTokens
	}
}
