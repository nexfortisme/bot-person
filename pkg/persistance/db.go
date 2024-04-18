package persistance

import (
	"fmt"
	"main/pkg/util"
	"sync"

	"github.com/surrealdb/surrealdb.go"
)

var (
	db   *surrealdb.DB
	once sync.Once

	err error
)

func initDB() {
	// Create a database connection
	db, err = surrealdb.New(util.GetDBHost())
	if err != nil {
		fmt.Println("Error connecting to database.")
		panic(err)
	}

	if _, err = db.Use(util.GetDBNamespace(), util.GetDBName()); err != nil {
		panic(err)
	}

	if _, err = db.Signin(map[string]interface{}{
		"user": util.GetDBUser(),
		"pass": util.GetDBPassword(),
	}); err != nil {
		panic(err)
	}

	fmt.Println("Database Connected.")
}

func GetDB() *surrealdb.DB {
	once.Do(func() {
		initDB()
	})
	return db
}
