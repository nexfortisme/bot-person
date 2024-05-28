package persistance

import (
	persistance "main/pkg/persistance/models"

	"github.com/surrealdb/surrealdb.go"
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

func GetUserStats(userId string) persistance.MyStats {

	user, _ := GetUser(userId)

	interactionCountData, err := db.Query("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId NOT IN [12, 13] GROUP ALL", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		panic(err)
	}

	goodBotCountData, err := db.Query("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [34] GROUP ALL", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		panic(err)
	}

	badBotCountData, err := db.Query("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [33] GROUP ALL", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		panic(err)
	}

	lootBoxCountData, err := db.Query("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [9] GROUP ALL", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		panic(err)
	}

	imageCountData, err := db.Query("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [16] GROUP ALL", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		panic(err)
	}

	chatCountData, err := db.Query("SELECT count() AS count FROM events WHERE eventUser = $userId AND eventId IN [12, 13] GROUP ALL", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		panic(err)
	}

	// Unmarshal data
	interactionCount := make([]persistance.UserEventCount, 1)
	_, err = surrealdb.UnmarshalRaw(interactionCountData, &interactionCount)
	if err != nil {
		panic(err)
	}

	goodBotCount := make([]persistance.UserEventCount, 1)
	_, err = surrealdb.UnmarshalRaw(goodBotCountData, &goodBotCount)
	if err != nil {
		panic(err)
	}

	badBotCount := make([]persistance.UserEventCount, 1)
	_, err = surrealdb.UnmarshalRaw(badBotCountData, &badBotCount)
	if err != nil {
		panic(err)
	}

	lootBoxCount := make([]persistance.UserEventCount, 1)
	_, err = surrealdb.UnmarshalRaw(lootBoxCountData, &lootBoxCount)
	if err != nil {
		panic(err)
	}

	imageCount := make([]persistance.UserEventCount, 1)
	_, err = surrealdb.UnmarshalRaw(imageCountData, &imageCount)
	if err != nil {
		panic(err)
	}

	chatCount := make([]persistance.UserEventCount, 1)
	_, err = surrealdb.UnmarshalRaw(chatCountData, &chatCount)
	if err != nil {
		panic(err)
	}

	var myStats persistance.MyStats

	myStats.InteractionCount = interactionCount[0].Count
	myStats.GoodBotCount = goodBotCount[0].Count
	myStats.BadBotCount = badBotCount[0].Count
	myStats.LootBoxCount = lootBoxCount[0].Count
	myStats.ImageCount = imageCount[0].Count
	myStats.ChatCount = chatCount[0].Count

	myStats.ImageTokens = user.UserStats.ImageTokens
	myStats.BonusStreak = user.UserStats.BonusStreak
	myStats.LastBonus = user.UserStats.LastBonus

	return myStats
}
