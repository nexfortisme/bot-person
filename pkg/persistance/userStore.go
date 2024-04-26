package persistance

import (
	"fmt"
	persistance "main/pkg/persistance/models"

	stateService "main/pkg/state/services"

	"github.com/surrealdb/surrealdb.go"
)

func GetUser(userId string) (*persistance.User, error) {

	db := GetDB()

	// Get user by ID
	data, err := db.Query("SELECT * FROM users WHERE UserId = $userId", map[string]interface{}{
		"userId": userId,
	})
	if err != nil {
		panic(err)
	}

	// Unmarshal data
	selectedUser := make([]persistance.User, 1)
	_, err = surrealdb.UnmarshalRaw(data, &selectedUser)
	if err != nil {
		panic(err)
	}

	if selectedUser[0].ID == "" || err != nil {

		newUser := persistance.User{}

		if userId != "SYSTEM" {
			// discordSession := stateService.GetDiscordSession()
			// discordUser, _ := discordSession.User(userId)
			// try {
			// 	newUser.Username = discordUser.
			// } catch(err Error){
			// 	newUser.Username = "SYSTEM"
			// }
		} else {
			newUser.Username = "SYSTEM"
		}

		newUser.UserId = userId
		newUser.UserStats.ImageTokens = 50

		resp, err := db.Create("users", newUser)
		if err != nil {
			return nil, err
		}

		// Unmarshal data
		createdUser := make([]persistance.User, 1)
		err = surrealdb.Unmarshal(resp, &createdUser)
		if err != nil {
			panic(err)
		}

		return &createdUser[0], nil
	}

	fmt.Printf("User: %+v\n", selectedUser[0])

	return &selectedUser[0], nil
}

func UpdateUser(updateUser persistance.User) bool {

	db := GetDB()

	if _, err := db.Update(updateUser.ID, updateUser); err != nil {
		return false
	}
	return true
}
