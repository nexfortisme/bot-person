package persistance

import (
	"strings"
	"time"
)

type BPInteraction int

const (
	BPChatInteraction       BPInteraction = iota // anything with the divinci chat model
	BPImageRequest                               // any image request
	BPBadBotInteraction                          // bad bot
	BPGoodBotInteraction                         // good bot
	BPBasicInteraction                           // any basic bot interaction
	BPSystemInteraction                          // TODO - Figure out what this means
	BPCreateUserInteraction                      // used by system for creating users
	BPLennyFaceInteracton
)

func GetUserStats(userId string) MyStats {

	user, err := GetUser(userId)
	if err != nil {
		panic(err)
	}

	var interactionCountData int64
	var goodBotCountData int64
	var badBotCountData int64
	var lootBoxCountData int64
	var imageCountData int64
	var chatCountData int64

	// Interaction Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType NOT IN (12, 13)", &interactionCountData, userId)
	if err != nil {
		panic(err)
	}

	// Good Bot Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (34)", &goodBotCountData, userId)
	if err != nil {
		panic(err)
	}

	// Bad Bot Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (33)", &badBotCountData, userId)
	if err != nil {
		panic(err)
	}

	// Loot Box Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (9)", &lootBoxCountData, userId)
	if err != nil {
		panic(err)
	}

	// Image Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (16)", &imageCountData, userId)
	if err != nil {
		panic(err)
	}

	// Chat Count
	err = RunQuery("SELECT count(*) AS count FROM events WHERE EventUser = ? AND EventType IN (12, 13)", &chatCountData, userId)
	if err != nil {
		panic(err)
	}

	var myStats MyStats

	myStats.InteractionCount = int(interactionCountData)
	myStats.GoodBotCount = int(goodBotCountData)
	myStats.BadBotCount = int(badBotCountData)
	myStats.LootBoxCount = int(lootBoxCountData)
	myStats.ImageCount = int(imageCountData)
	myStats.ChatCount = int(chatCountData)

	myStats.ImageTokens = user.ImageTokens
	myStats.BonusStreak = user.BonusStreak


	// Convert the string timestamp to a time.Time object
	lastBonus, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", strings.Split(user.LastBonus, " m=")[0])
	if err != nil {
		// Set a default time if parsing fails
		lastBonus = time.Now()
	}

	// Set the LastBonus field to the parsed time
	myStats.LastBonus = lastBonus

	return myStats
}
