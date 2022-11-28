package persistance

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func AddImageTokens(tokenAmount float64, userId string) bool {

	user, err := getUser(userId)

	if err != nil {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, tokenAmount, time.Time{}}})
		return true
	} else {
		user.UserStats.ImageTokens += tokenAmount
		return updateUser(user)
	}

}

func TransferrImageTokens(tokenAmount float64, fromUserId string, toUserId string) bool {

	fromUser, fromErr := getUser(fromUserId)
	toUser, toErr := getUser(toUserId)

	// Checking to see if user exists
	if fromErr != nil {
		return false
	} else {

		// The User exists but is trying to transferr more tokens then they have
		if fromUser.UserStats.ImageTokens < tokenAmount {
			return false
		}

		// Checking to see if the toUser exists
		if toErr != nil {

			// toUser doesn't exist
			// Creates user and assigns them the number of tokens that is being transferred
			botTracking.UserStats = append(botTracking.UserStats, UserStruct{toUserId, UserStatsStruct{0, 0, 0, 0, tokenAmount, time.Time{}}})
			fromUser.UserStats.ImageTokens -= tokenAmount
			return updateUser(fromUser)
		} else {
			toUser.UserStats.ImageTokens += tokenAmount
			fromUser.UserStats.ImageTokens -= tokenAmount
			return updateUser(toUser) && updateUser(fromUser)
		}
	}

}

func UseImageToken(userId string) bool {

	user, err := getUser(userId)

	if err != nil {
		return false
	} else {
		if user.UserStats.ImageTokens < 1 {
			return false
		} else {
			user.UserStats.ImageTokens--
			return updateUser(user)
		}
	}

}

// Is this needed, can be canned
func UserHasTokens(userId string) bool {
	user, err := getUser(userId)

	if err != nil {
		return false
	} else {
		if user.UserStats.ImageTokens <= 0 {
			return false
		} else {
			return true
		}
	}

}

func GetUserTokenCount(userId string) float64 {
	user, err := getUser(userId)

	if err != nil {
		return 0
	} else {
		return user.UserStats.ImageTokens
	}
}

func SetUserTokenCount(userId string, tokenAmount float64) bool {
	user, err := getUser(userId)

	if err != nil {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, tokenAmount, time.Time{}}})
		return true
	} else {
		user.UserStats.ImageTokens = tokenAmount
		return updateUser(user)
	}
}

func RemoveUserTokens(userId string, tokenAmount float64) bool {
	user, err := getUser(userId)

	if err != nil {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 0, 0, 0, 0, time.Time{}}})
		return true
	} else {
		if (user.UserStats.ImageTokens - tokenAmount) <= 0 {
			user.UserStats.ImageTokens = 0
			return updateUser(user)
		} else {
			user.UserStats.ImageTokens -= tokenAmount
			return updateUser(user)
		}
	}
}

func GetUserReward(userId string) (float64, error) {

	user, err := getUser(userId)

	// Checking for error and setting userId on returned user struct as needed
	if err != nil {
		user.UserId = userId
	}

	// Checking to see if the user has a LastBonus time
	if (user.UserStats.LastBonus != time.Time{}) {
		
		// Checking diff between now and lastBonus time
		diff := time.Now().Sub(user.UserStats.LastBonus)


		// if diff is less than 1 day (86400 seconds), then throws error
		if diff.Seconds() <= 86400 {

			// Doing math to for countdown to next bonus
			nextBonus := user.UserStats.LastBonus.AddDate(0, 0, 1)
			timeToNextBonus := nextBonus.Sub(time.Now())

			// TODO - Parse timeToNextBonus better
			errString := "Please try again in: " + timeToNextBonus.String()

			// returning error
			return -1.0, errors.New(errString)
		}
	}

	// Setting random seed and generating a, value safe, token amount
	randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))
	reward := randomizer.Intn(45) + 5
	rewardf64 := float64(reward) / 10.0
	rewardString := fmt.Sprintf("%.2f", rewardf64)
	finalReward, _ := strconv.ParseFloat(rewardString, 64)

	// Updating User Record
	user.UserStats.LastBonus = time.Now()
	user.UserStats.ImageTokens += finalReward

	if !updateUser(user) {
		return -1, errors.New("Error updating user record")
	} else {
		return finalReward, nil
	}

}
