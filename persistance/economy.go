package persistance

import (
	"errors"
	"fmt"
	"main/util"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
)

func AddImageTokens(tokenAmount float64, userId string) bool {

	user, err := getUser(userId)

	if err != nil {
		createAndAddUser(userId, 0, 0, 0, 0, util.LowerFloatPrecision(tokenAmount))
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
			createAndAddUser(toUserId, 0, 0, 0, 0, util.LowerFloatPrecision(tokenAmount))
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
		createAndAddUser(userId, 0, 0, 0, 0, util.LowerFloatPrecision(tokenAmount))
		return true
	} else {
		user.UserStats.ImageTokens = tokenAmount
		return updateUser(user)
	}
}

func RemoveUserTokens(userId string, tokenAmount float64) bool {
	user, err := getUser(userId)

	if err != nil {
		createAndAddUser(userId, 0, 0, 0, 0, 0)
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

func GetUserReward(userId string) (float64, string, error) {

	user, err := getUser(userId)
	returnString := ""

	modifier := 1

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
			formattedString := durafmt.Parse(timeToNextBonus).LimitFirstN(3)

			errString := "Please try again in: " + formattedString.String()

			// returning error
			return -1.0, "", errors.New(errString)
		}

		timeWindow := user.UserStats.LastBonus.Add(time.Hour * 48)
		timeWindowDiff := time.Now().Sub(timeWindow)

		// Missed Window
		if timeWindowDiff > 0 {
			returnString = fmt.Sprintf("Bonus Not Redeemed within 24 hours. Streak Reset. \nCurrent Streak: %d", 1)
			user.UserStats.BonusStreak = 1
		} else {
			user.UserStats.BonusStreak++
			streak := user.UserStats.BonusStreak

			if streak%10 == 0 && streak%100 != 0 && streak%50 != 0 {
				returnString = fmt.Sprintf("Congrats on keeping the streak alive. Current Streak: %d. Bonus Modifier: 2x", streak)
				modifier = 2
			} else if streak%25 == 0 && streak%50 != 0 && streak%100 != 0 {
				returnString = fmt.Sprintf("Great work on keeping the streak alive! Current Streak: %d. Bonus Modifier: 5x", streak)
				modifier = 5
			} else if streak%50 == 0 && streak%100 != 0 {
				returnString = fmt.Sprintf("Wow! That's a long time. Current Streak: %d. Bonus Modifier: 10x", streak)
				modifier = 10
			} else if streak%100 == 0 {
				returnString = fmt.Sprintf("Few people ever reach is this far, Congratulations! Current Streak: %d. Bonus Modifier: 50x", streak)
				modifier = 50
			} else {
				returnString = fmt.Sprintf("Current Bonus Streak: %d", streak)
			}
		}

	}

	// Setting random seed and generating a, value safe, token amount
	randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))
	reward := randomizer.Intn(45) + 5
	reward *= modifier
	rewardf64 := float64(reward) / 10.0
	finalReward := util.LowerFloatPrecision(rewardf64)

	// Updating User Record
	user.UserStats.LastBonus = time.Now()
	user.UserStats.ImageTokens += finalReward

	if !updateUser(user) {
		return -1, "", errors.New("error updating user record")
	} else {
		return finalReward, returnString, nil
	}

}

func BuyLootbox(userId string) (int, int, error) {

	user, err := getUser(userId)

	if err != nil {
		return -1, -1, err
	}

	if user.UserStats.ImageTokens < 2.5 {
		return -1, -1, errors.New("you do not have the 2.5 tokens needed to purchase a lootbox")
	} else {
		user.UserStats.ImageTokens -= 2.5
	}

	random := rand.New(rand.NewSource(time.Now().UnixMilli()))
	lootboxSeed := random.Intn(9999999999) + 1000000000

	val := hashLootBoxSeed(lootboxSeed)
	reward := 0

	if val <= 7992 {
		reward += 1
	} else if val > 7992 && val <= 9590 {
		reward += 3
	} else if val > 9590 && val <= 9910 {
		reward += 10
	} else if val > 9910 && val <= 9974 {
		reward += 50
	} else if val > 9974 && val <= 10000 {
		reward += 250
	}

	user.UserStats.ImageTokens += float64(reward)

	if !updateUser(user) {
		return -1, -1, errors.New("error updating user record")
	} else {
		return reward, lootboxSeed, nil
	}

}

func hashLootBoxSeed(bar int) int {

	/*
	* Blue – 100 items – 79.92%
	* Purple – 20 items – 15.98%
	* Pink – 4 items – 3.2%
	* Red – 0.8 items – 0.64%
	* Yellow – 0.32 items – 0.26%
	 */

	modifier := 1000000000

	total := 1

	for i := 1; i < 11; i++ {
		num := (bar / modifier) % 10
		// fmt.Println(num)
		modifier /= 10

		if i <= 8 {
			if num != 0 {
				total *= num
			} else {
				total += 10
			}
		} else {
			if num != 0 {
				total *= num
			} else {
				total += 10
			}
		}

	}

	return total % 10000
}

func APictureIsWorthAThousand(incomingMessage string, m *discordgo.MessageCreate) {

	// Looking at messages in the channel and returning WORD_COUNT / 1000 number of tokens
	// ie. A picture is worth 1000 words
	wordCount := len(strings.Fields(incomingMessage))
	tokenValue := fmt.Sprintf("%.2f", (float64(wordCount) / 1000.0))
	tokenAddAmount, _ := strconv.ParseFloat(tokenValue, 64)

	AddImageTokens(tokenAddAmount, m.Author.ID)
}
