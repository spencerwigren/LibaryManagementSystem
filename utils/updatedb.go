package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func AddBookInfo(title string, pageNumber int, author string, tableType string, db *sql.DB) {

	insertSQL := `INSERT INTO books (title, pageNumber, author, type) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertSQL, title, pageNumber, author, tableType)
	checkError(err)
	println("Success")

}

func AddMovieInfo(movieTitle string, db *sql.DB) {

	insertSQL := `INSERT INTO movies (title) VALUES (?)`
	_, err := db.Exec(insertSQL, movieTitle)
	checkError(err)

}

func AddVideoGameInfo(gameTitle string, db *sql.DB) {

	insertSQL := `INSERT INTO videoGames (title) values (?)`
	_, err := db.Exec(insertSQL, gameTitle)
	checkError(err)
}

func AddUserInfo(name string, db *sql.DB) {

	insertSQL := `INSERT INTO users (name) 	VALUES (?)`
	_, err := db.Exec(insertSQL, name)
	checkError(err)
}
