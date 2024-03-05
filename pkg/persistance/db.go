package persistance

import (
	"fmt"
	"sync"

	"main/pkg/util"

	"github.com/go-pg/pg/v10"
)

var (
	db   *pg.DB
	once sync.Once
)

func initDB() {
	// Create a database connection
	db = pg.Connect(&pg.Options{
		Addr:     util.GetDBHost() + ":5432",
		User:     util.GetDBUser(),
		Password: util.GetDBPassword(),
		Database: util.GetDBName(),
	})

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

func GetDB() *pg.DB {
	once.Do(func() {
		initDB()
	})
	return db
}
