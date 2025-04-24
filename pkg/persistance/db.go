package persistance

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

var (
	db   *sqlite.Conn
	once sync.Once

	err error
)

func initDB() {
	var err error
	db, err = sqlite.OpenConn("db.sqlite", 0)
	if err != nil {
		fmt.Println("Error connecting to database.")
		panic(err)
	}

	InitializeDatabase(db)

	fmt.Println("Database Connected.")
}

func GetDB() *sqlite.Conn {
	once.Do(func() {
		initDB()
	})
	return db
}

// Initializing the DB with the necessary tables for the bot to function
func InitializeDatabase(db *sqlite.Conn) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS Users (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		UserId TEXT,
		Username TEXT,
		ImageTokens REAL,
		BonusStreak INTEGER,
		LastBonus TIMESTAMP
	);`

	createUserStatsTable := `
	CREATE TABLE IF NOT EXISTS UserStats (
		UserID TEXT PRIMARY KEY,
		InteractionCount INTEGER,
		ChatCount INTEGER,
		GoodBotCount INTEGER,
		BadBotCount INTEGER,
		ImageCount INTEGER,
		LootBoxCount INTEGER,
		FOREIGN KEY (UserID) REFERENCES Users(ID)
	);`

	createStocksTable := `
	CREATE TABLE IF NOT EXISTS Stocks (
		UserID TEXT,
		StockTicker TEXT,
		StockCount REAL,
		PRIMARY KEY (UserID, StockTicker),
		FOREIGN KEY (UserID) REFERENCES Users(ID)
	);`

	createLogsTable := `
	CREATE TABLE IF NOT EXISTS Logs (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		EventType TEXT,
		UserID TEXT,
		Description TEXT,
		GuildID TEXT,
		Timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (UserID) REFERENCES Users(ID)
	);`

	createRewardStatusTable := `
	CREATE TABLE IF NOT EXISTS RewardStatus (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		StatusID INTEGER,
		StatusName TEXT
	);`

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS Events (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		EventType TEXT,
		EventTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		EventUser TEXT,
		EventData TEXT,
		EventServer TEXT,
		FOREIGN KEY (EventUser) REFERENCES Users(ID)
	);`

	createUserAttributesTable := `
	CREATE TABLE IF NOT EXISTS UserAttributes (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		UserId TEXT,
		Attribute TEXT,
		Value TEXT,
		FOREIGN KEY (UserId) REFERENCES Users(ID)
	);`

	// Execute the table creation statements
	tables := []string{createUsersTable, createUserStatsTable, createStocksTable, createLogsTable, createRewardStatusTable, createEventsTable, createUserAttributesTable}

	for _, table := range tables {
		err := sqlitex.Execute(db, table, nil)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}

	fmt.Println("Database tables initialized successfully.")
}

// RunQuery executes a given SQL query with optional parameters and returns the results.
func RunQuery(query string, output interface{}, params ...interface{}) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("error preparing query: %w", err)
	}
	defer stmt.Finalize()

	for i, param := range params {
		var paramIndex int = i + 1

		switch v := param.(type) {
		case string:
			stmt.BindText(paramIndex, v)
		case int64:
			stmt.BindInt64(paramIndex, v)
		case int:
			stmt.BindInt64(paramIndex, int64(v))
		case float64:
			stmt.BindFloat(paramIndex, v)
		case bool:
			stmt.BindBool(paramIndex, v)
		case []byte:
			stmt.BindBytes(paramIndex, v)
		case nil:
			stmt.BindNull(paramIndex)
		default:
			return fmt.Errorf("unsupported parameter type %T for parameter %d", param, paramIndex)
		}
	}

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return fmt.Errorf("error executing query: %w", err)
		}
		if !hasRow {
			break
		}

		// Skip result processing if output is nil
		if output == nil {
			continue
		}

		outputValue := reflect.ValueOf(output)
		if outputValue.Kind() != reflect.Ptr {
			return fmt.Errorf("output must be a pointer")
		}

		elemValue := outputValue.Elem()

		// Handle basic types (int64, string, etc.)
		if elemValue.Kind() == reflect.Int64 || elemValue.Kind() == reflect.Int {
			elemValue.SetInt(stmt.ColumnInt64(0))
			continue
		}

		if elemValue.Kind() == reflect.Struct {
			// Populate a single struct
			for i := 0; i < elemValue.NumField(); i++ {
				field := elemValue.Field(i)
				switch field.Kind() {
				case reflect.Int64:
					field.SetInt(stmt.ColumnInt64(i))
				case reflect.Int:
					field.SetInt(stmt.ColumnInt64(i))
				case reflect.String:
					field.SetString(stmt.ColumnText(i))
				case reflect.Float64:
					field.SetFloat(stmt.ColumnFloat(i))
				default:
					return fmt.Errorf("unsupported field type %v", field.Kind())
				}
			}
		} else if elemValue.Kind() == reflect.Slice {
			// Populate a slice
			elementType := elemValue.Type().Elem()
			newElement := reflect.New(elementType).Elem()

			for i := 0; i < elementType.NumField(); i++ {
				field := elementType.Field(i)
				switch field.Type.Kind() {
				case reflect.Int64:
					newElement.Field(i).SetInt(stmt.ColumnInt64(i))
				case reflect.Int:
					newElement.Field(i).SetInt(stmt.ColumnInt64(i))
				case reflect.String:
					newElement.Field(i).SetString(stmt.ColumnText(i))
				case reflect.Float64:
					newElement.Field(i).SetFloat(stmt.ColumnFloat(i))
				default:
					return fmt.Errorf("unsupported field type %v for field %s", field.Type, field.Name)
				}
			}

			elemValue.Set(reflect.Append(elemValue, newElement))
		} else {
			return fmt.Errorf("output must be a pointer to a struct, slice, or basic type")
		}
	}

	return nil
}
