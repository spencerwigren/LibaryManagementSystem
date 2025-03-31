package utils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func AddBookInfo(title string, pageNumber int, author string, db *sql.DB) {

	insertSQL := `INSERT INTO books (title, pageNumber, author) VALUES (?, ?, ?)`
	_, err := db.Exec(insertSQL, title, pageNumber, author)
	if err != nil {
		log.Fatal(err)
	}
	println("Success")

}

func AddMovieInfo(movieTitle string, db *sql.DB) {

	insertSQL := `INSERT INTO movies (title) VALUES (?)`
	_, err := db.Exec(insertSQL, movieTitle)
	if err != nil {
		log.Fatal(err)
	}

}

func AddVideoGameInfo(gameTitle string, db *sql.DB) {

	insertSQL := `INSERT INTO videoGames (title) values (?)`
	_, err := db.Exec(insertSQL, gameTitle)
	if err != nil {
		log.Fatal(err)
	}
}

func AddUserInfo(name string, db *sql.DB) {

	insertSQL := `INSERT INTO users (name) 	VALUES (?)`
	_, err := db.Exec(insertSQL, name)
	if err != nil {
		log.Fatal(err)
	}
}
