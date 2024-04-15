package persistance

import (
	"fmt"
	persistance "main/pkg/persistance"
	persistanceModels "main/pkg/persistance/models"

	"github.com/bwmarrin/discordgo"
	"github.com/surrealdb/surrealdb.go"
)

func GetUserByDiscordUserId(discordUserId string, s *discordgo.Session) (persistanceModels.User, error) {

	databaseConnection := persistance.GetDB()
	userQueryString := fmt.Sprintf("users:%s", discordUserId)

	userResponse, err := databaseConnection.Select(userQueryString)
	if err != nil {
		newUser := persistanceModels.User{}
		newUser.DiscordUserId = discordUserId
		databaseConnection.Create("users", newUser)
		return newUser, nil
	}

	selectedUser := new(persistanceModels.User)
	err = surrealdb.Unmarshal(userResponse, &selectedUser)
	if err != nil {
		return persistanceModels.User{}, err
	}

	return *selectedUser, nil
}

func UpsertUser(dbUser persistanceModels.User) (bool, error) {

	databaseConnection := persistance.GetDB()
	userQueryString := fmt.Sprintf("users:%s", dbUser.DiscordUserId)

	if _, err := databaseConnection.Update(userQueryString, dbUser); err != nil {
		return false, err
	}

	return true, nil
}
