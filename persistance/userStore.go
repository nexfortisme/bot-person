package persistance

import (
	"errors"
	"time"
)

func getUser(userId string) (UserStruct, error) {

	// TODO - Change this to use a Map instead
	for _, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {
			return element, nil
		}
	}
	
	return UserStruct{"bad user", UserStatsStruct{-1, -1, -1, -1, -1, time.Time{}}}, errors.New("unable to find user")
}

func updateUser(updateUser UserStruct) bool {
	// TODO - Change this to use a Map instead
	for index, element := range botTracking.UserStats {
		if element.UserId != updateUser.UserId {
			continue
		} else {
			botTracking.UserStats[index] = updateUser
			return true
		}
	}
	return false
}

func createUser(userId string, messageCount int, goodBotCount int, badBotCount int, imageCount int, imageTokens int) bool {
	botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{messageCount, goodBotCount, badBotCount, imageCount, float64(imageTokens), time.Time{}}})
	return true
}

// TODO - Delete User Function
