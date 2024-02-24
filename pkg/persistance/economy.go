package persistance

import (
	"errors"
	"fmt"
	persistance "main/pkg/persistance/models"
	"main/pkg/util"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
)

func AddBotPersonTokens(tokenAmount float64, userId string) bool {

	user, err := getUser(userId)

	if err != nil {
		createAndAddUser(userId, 0, 0, 0, 0, util.LowerFloatPrecision(tokenAmount))
		return true
	} else {
		user.UserStats.ImageTokens += tokenAmount
		return updateUser(user)
	}

}

func TransferBotPersonTokens(tokenAmount float64, fromUserId string, toUserId string) bool {

	fromUser, fromErr := getUser(fromUserId)
	toUser, toErr := getUser(toUserId)

	// Checking to see if user exists
	if fromErr != nil {
		return false
	} else {

		// The User exists but is trying to transfer more tokens than they have
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
			user.UserStats.ImageTokens -= 10
			return updateUser(user)
		}
	}

}

// Is this needed, can be canned
func UserHasTokens(userId string) bool {
	user, err := getUser(userId)

	if err != nil {
		return false
	}

	return user.UserStats.ImageTokens > 0
}

func GetUserTokenCount(userId string) float64 {
	user, err := getUser(userId)

	if err != nil {
		return 0
	} else {
		return user.UserStats.ImageTokens
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

	// Getting user and setting necessary variables
	user, err := getUser(userId)
	returnString := ""
	modifier := 1

	// Checking for Get User Error and setting appropriate values
	if err != nil {
		user.UserId = userId
		addUser(user)
	}

	// Checking to see if the user has a LastBonus time
	if (user.UserStats.LastBonus != time.Time{}) {

		// Checking diff between now and lastBonus time
		diff := time.Since(user.UserStats.LastBonus)

		// if diff is less than 1 day (86400 seconds), then throws error
		if diff.Seconds() <= 86400 {

			// Doing math to for countdown to next bonus
			nextBonus := user.UserStats.LastBonus.AddDate(0, 0, 1)
			timeToNextBonus := time.Until(nextBonus)
			formattedString := durafmt.Parse(timeToNextBonus).LimitFirstN(3)

			errString := "Please try again in: " + formattedString.String()

			return -1.0, "", errors.New(errString)
		}

		timeWindow := user.UserStats.LastBonus.Add(time.Hour * 48)
		timeWindowDiff := time.Since(timeWindow)

		// Missed Window
		if timeWindowDiff > 0 {

			// If the user has had a time set for trying to save their streak
			if (user.UserStats.HoldStreakTimer != time.Time{}) {
				returnString = fmt.Sprintf("Streak Not Saved in time. Streak Reset. \nCurrent Streak: %d", 1)
				user.UserStats.BonusStreak = 1
				user.UserStats.HoldStreakTimer = time.Time{}
			} else {

				// Used for displaying the discord relative time for the user to save their streak
				inFiveMinutes := time.Now().Add(time.Minute * 5).Unix()

				// Not sure if the relative time string work for users outside the time zone of the bot
				errString := fmt.Sprintf("Bonus Not Redeemed within 24 hours. To save your streak, use `/save-streak` <t:%d:R>. `/save-streak` will use a save token or purchase one for 1/2 of your current tokens. \nCurrent Streak: %d", inFiveMinutes, user.UserStats.BonusStreak)

				user.UserStats.HoldStreakTimer = time.Now()

				if !updateUser(user) {
					return -1, "", errors.New("error updating user record")
				} else {
					return -1, "", errors.New(errString)
				}
			}

		} else {
			user.UserStats.BonusStreak++
			streak := user.UserStats.BonusStreak

			returnString, modifier = util.GetStreakStringAndModifier(streak)
		}

	}

	// Get Final Bonus Reward
	finalReward := util.GetUserBonus(5, 50, modifier)

	// Updating User Record
	user.UserStats.LastBonus = time.Now()
	user.UserStats.ImageTokens += finalReward

	if !updateUser(user) {
		return -1, "", errors.New("error updating user record")
	} else {
		return finalReward, returnString, nil
	}

}

func BuyLootbox(userId string) (float64, int, error) {

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
	reward := 0.0

	if val <= 7992 {
		reward += 3.63
	} else if val > 7992 && val <= 9590 {
		reward += 8
	} else if val > 9590 && val <= 9910 {
		reward += 15
	} else if val > 9910 && val <= 9974 {
		reward += 25
	} else if val > 9974 && val <= 10000 {
		reward += 50
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

	AddBotPersonTokens(tokenAddAmount, m.Author.ID)
}

func AddStock(userId string, stockTicker string, quantity float64) error {

	user, err := getUser(userId)

	if err != nil {
		return err
	}

	userPortfolio := user.UserStats.Stocks

	for index, element := range userPortfolio {
		if element.StockTicker == stockTicker {
			element.StockCount += quantity
			userPortfolio[index] = element
			updateUser(user)
			return nil
		}
	}

	newStock := persistance.Stock{StockTicker: stockTicker, StockCount: quantity}
	user.UserStats.Stocks = append(user.UserStats.Stocks, newStock)
	updateUser(user)

	return nil
}

func RemoveStock(userId string, stockTicker string, quantity float64) error {

	user, err := getUser(userId)

	if err != nil {
		return err
	}

	userPortfolio := user.UserStats.Stocks

	for index, element := range userPortfolio {
		if element.StockTicker == stockTicker {
			element.StockCount -= quantity
			userPortfolio[index] = element
			updateUser(user)
			return nil
		}
	}

	return errors.New("stock not found")
}

func GetUserStock(userId string, stockTicker string) (persistance.Stock, error) {

	user, err := getUser(userId)

	if err != nil {
		return persistance.Stock{}, err
	}

	userPortfolio := user.UserStats.Stocks

	for _, element := range userPortfolio {
		if element.StockTicker == stockTicker {
			return element, nil
		}
	}

	return persistance.Stock{}, errors.New("stock not found")

}
