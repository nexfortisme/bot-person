package persistance

import (
	"errors"
	"fmt"

	persistance "main/pkg/persistance/models"
	"main/pkg/util"

	"github.com/go-pg/pg/v10"
)

var (
	databaseConnection *pg.DB
)

func Connect() {
	// Create a database connection
	databaseConnection := pg.Connect(&pg.Options{
		Addr:     util.GetDBHost() + ":5432",
		User:     util.GetDBUser(),
		Password: util.GetDBPassword(),
		Database: util.GetDBName(),
	})

	// Ensure the database connection is closed when main() exits
	defer databaseConnection.Close()

	fmt.Println("Database Connected.")
}

func GetUser(discord_user_id string) (persistance.DBUser, error) {

	var dbUser persistance.DBUser

	err := databaseConnection.Model(dbUser).Where("discord_user_id = ?", discord_user_id).Select()
	if err != nil {
		return persistance.DBUser{}, err
	}

	return dbUser, nil
}

func CreateUser(dbUser persistance.DBUser) (persistance.DBUser, error) {

	result, err := databaseConnection.Model(&dbUser).Insert()
	if err != nil {
		return persistance.DBUser{}, err
	}

	if result.RowsAffected() < 1 {
		return persistance.DBUser{}, errors.New("db user - creation failed")
	}

	return dbUser, nil
}

func UpdateUser(dbUser persistance.DBUser) (bool, error) {

	result, err := databaseConnection.Model(&dbUser).Update()
	if err != nil {
		return false, err
	}

	if result.RowsAffected() < 1 {
		return false, errors.New("db user - update failed")
	}

	return true, nil
}

func GetDatabaseConnection() *pg.DB {
	return databaseConnection
}
