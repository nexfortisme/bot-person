package persistance

import (
	"fmt"
	"sync"

	"github.com/surrealdb/surrealdb.go"
)

var (
	db   *surrealdb.DB
	once sync.Once
)

func initDB() {
	// Create a database connection
	db, err := surrealdb.New("ws://localhost:8000/rpc")
	if err != nil {
		fmt.Println("Error connecting to database.")
		panic(err)
	}

	if _, err = db.Signin(map[string]interface{}{
		"user": "501",
		"pass": "root",
	}); err != nil {
		panic(err)
	}

	if _, err = db.Use("botPerson", "botPerson"); err != nil {
		panic(err)
	}

	fmt.Println("Database Connected.")
}

// func HelloWorld() {
// 	var greeting string
// 	ctx := context.Background()

// 	_, err := db.QueryOneContext(ctx, pg.Scan(&greeting), "SELECT 'Hello, World!'")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(greeting)
// }

// TODO - Delete User

func GetDB() *surrealdb.DB {
	once.Do(func() {
		initDB()
	})
	return db
}
