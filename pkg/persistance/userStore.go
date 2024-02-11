package persistance

import (
	"errors"
	persistance "main/pkg/persistance/models"
	"time"
)

func getUser(userId string) (persistance.User, error) {

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

func updateUser(updateUser persistance.User) bool {
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

func getNewUser(userId string, messageCount int, goodBotCount int, badBotCount int, imageCount int, imageTokens float64) persistance.User {

	newUser := persistance.User{UserId: userId, UserStats: persistance.UserStats{MessageCount: messageCount, GoodBotCount: goodBotCount, BadBotCount: badBotCount, ImageCount: imageCount, ImageTokens: imageTokens, LastBonus: time.Time{}, BonusStreak: 0, HoldStreakTimer: time.Time{}, SaveStreakTokens: 0, Stocks: []persistance.Stock{}}}

	return newUser
}

func addUser(user persistance.User) bool {
	botTracking.UserStats = append(botTracking.UserStats, user)
	return true
}

func GetUserStats(userId string) (persistance.UserStats, error) {

	user, err := getUser(userId)

	if err != nil {
		return persistance.UserStats{}, err
	}

	return user.UserStats, nil
}

func UpdateUserStats(userId string, stats persistance.UserStats) bool {

	user, err := getUser(userId)

	if err != nil {
		return false
	}

	user.UserStats = stats

	return updateUser(user)
}

// TODO - Delete User Function
