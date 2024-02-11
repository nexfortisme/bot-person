package persistance

import (
	"fmt"

	"main/pkg/util"

	"github.com/go-pg/pg/v10"
)

func Connect() {
	// Create a database connection
	db := pg.Connect(&pg.Options{
		Addr:     util.GetDBHost() + ":5432",
		User:     util.GetDBUser(),
		Password: util.GetDBPassword(),
		Database: util.GetDBName(),
	})
	defer db.Close() // Ensure the database connection is closed when main() exits

	// // Example: Using the DB object to perform a query
	// var greeting string
	// _, err := db.QueryOne(pg.Scan(&greeting), "SELECT 'Hello, world!'")
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	// 	os.Exit(1)
	// }

	fmt.Println("Database Connected.")
}
