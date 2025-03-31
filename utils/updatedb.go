package utils

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func AddBookInfo(title string, pageNumber int, author string, db *sql.DB) {

	dateTime := time.Now()
	log.Printf("DATETIME: %v", dateTime)

	insertSQL := `INSERT INTO books (title, pageNumber, author, time) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertSQL, title, pageNumber, author, dateTime)
	if err != nil {
		log.Fatal(err)
	}
}

func AddMovieInfo(movieTitle string, db *sql.DB) {

	insertSQL := `INSERT INTO movies (title, time) VALUES (?, ?)`
	_, err := db.Exec(insertSQL, movieTitle, time.Now())
	if err != nil {
		log.Fatal(err)
	}

}

func AddVideoGameInfo(gameTitle string, db *sql.DB) {

	insertSQL := `INSERT INTO videoGames (title, time) values (?, ?)`
	_, err := db.Exec(insertSQL, gameTitle, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}

func AddUserInfo(name string, db *sql.DB) {

	insertSQL := `INSERT INTO users (name, time) VALUES (?, ?)`
	_, err := db.Exec(insertSQL, name, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}
