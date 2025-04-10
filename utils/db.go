package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() {
	// Open (or create) a database
	db, err := sql.Open("sqlite3", "projectdb.db")
	// checkError(err)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// drop tables
	// dropTables(db, err)

	// Create the table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		time DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS books (
		bookId INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		pageNumber INTEGER NOT NULL,
		author TEXT NOT NULL,
		time DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS movies(
		movieId INTEGER PRIMARY KEY AUTOINCREMENT, 
		title TEXT NOT NULL,
		rating TEXT NOT NULL,
		releaseYear INTERGER NOT NULL,
		time DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS videoGames (
		videoGameId INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		rating TEXT NOT NULL,
		releaseYear INTERGER NOT NULL,
		time DATETIME NOT NULL
	)`

	// Creating the db tables
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database and tables created successfully!")
}

// func dropTables(db *sql.DB, err error) {
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	tables := []string{"users", "books", "movies", "videoGames"}

// 	for _, table := range tables {
// 		query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", table)
// 		_, err := db.Exec(query)
// 		if err != nil {
// 			log.Fatalf("Failed to delete table %s: %v", table, err)
// 		}
// 		fmt.Printf("Table %s deleted successfully!\n", table)
// 	}

// 	log.Println("Tables deleted successfully!")

// }
