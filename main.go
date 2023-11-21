package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var (
// db *sql.DB
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "abc123"
	dbname   = "bot-person"
)

type User struct {
	id           uuid.UUID `json:"id"`
	dateCreated  string    `json:"dateCreated"`
	dateModified string    `json:"dateModified"`
	username     string    `json:"username"`
	password     string    `json:"password"`
	salt         string    `json:"salt"`
}

func main() {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable dbname=%s", host, port, user, password, dbname)

	r := gin.Default()

	fmt.Println("Connecting to database...")
	fmt.Println(connString)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	r.GET("/users", func(c *gin.Context) {
		users, err := getAllUsers(db)
		fmt.Println(users)
		fmt.Println("Number of users: ", len(users));
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"users": users})
	})

	r.Run(":9000")
}

func getAllUsers(db *sql.DB) ([]User, error) {

	if db == nil {
		panic("db is nil")
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.id, &user.dateCreated, &user.dateModified, &user.username, &user.password, &user.salt); err != nil {
			fmt.Println("Error scanning rows: ", err)
			return nil, err
		}

		fmt.Println(user)
		users = append(users, user)
	}

	fmt.Printf("users: %v", users)

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
