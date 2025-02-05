package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// getting demo writen
// the struct are going to act as the database for right now
// the terminal is going ot act as the UI for now

// type book struct {
// 	title      string
// 	pageNumber int
// 	author     string
// }

// type movie struct {
// 	title string
// }

// type videoGame struct {
// 	title string
// }

func main() {
	// bookInfo := addBookInfo("testBook", 125, "Me")
	// movieInfo := addMovieInfo("testMovie")
	// videoGameInfo := addVideoGameInfo("testVideoGame")

	// println(bookInfo.title)
	// println(movieInfo.title)
	// println(videoGameInfo.title)

	initdatabase() // used for creating the database

	db, err := sql.Open("sqlite3", "example.db")
	checkError(err)
	defer db.Close()

	addBookInfo("testBook", 125, db, err)
	queryDB(db, err)

}

func initdatabase() {
	// Open (or create) a database
	db, err := sql.Open("sqlite3", "example.db")
	checkError(err)
	defer db.Close()

	// Create a table
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER
	);`

	_, err = db.Exec(createTableSQL)
	checkError(err)

	fmt.Println("Database and table created successfully!")

}

func queryDB(db *sql.DB, err error) {
	rows, err := db.Query("SELECT id, name, age FROM users;")
	checkError(err)
	defer rows.Close()

	fmt.Println("Users:")
	for rows.Next() {
		var id int
		var name string
		var age int

		err = rows.Scan(&id, &name, &age)
		checkError(err)

		fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func addBookInfo(title string, pageNumber int, db *sql.DB, err error) {
	// b := book{title: title}
	// b.pageNumber = pageNumber
	// b.author = author

	insertSQL := `INSERT INTO users (name, age) VALUES (?, ?)`
	_, err = db.Exec(insertSQL, title, pageNumber)
	checkError(err)

}

// func addMovieInfo(title string) *movie {
// 	m := movie{title: title}
// 	return &m
// }

// func addVideoGameInfo(title string) *videoGame {
// 	g := videoGame{title: title}
// 	return &g
// }

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
