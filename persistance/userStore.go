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

	return getNewUser("bad user", -1, -1, -1, -1, -1), errors.New("unable to find user")
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

func createAndAddUser(userId string, messageCount int, goodBotCount int, badBotCount int, imageCount int, imageTokens float64) bool {
	botTracking.UserStats = append(botTracking.UserStats, getNewUser(userId, messageCount, goodBotCount, badBotCount, imageCount, imageTokens))
	return true
}

func getNewUser(userId string, messageCount int, goodBotCount int, badBotCount int, imageCount int, imageTokens float64) UserStruct {

	newUser := UserStruct{userId, UserStatsStruct{messageCount, goodBotCount, badBotCount, imageCount, imageTokens, time.Time{}, 0, time.Time{}, 0, []UserStock{}}}

	return newUser
}

func GetUserStats(userId string) (UserStatsStruct, error) {

	user, err := getUser(userId)

	if err != nil {
		return UserStatsStruct{}, err
	}

	return user.UserStats, nil
}

func UpdateUserStats(userId string, stats UserStatsStruct) bool {

	user, err := getUser(userId)

	if err != nil {
		return false
	}

	user.UserStats = stats

	return updateUser(user)
}

// TODO - Delete User Function
