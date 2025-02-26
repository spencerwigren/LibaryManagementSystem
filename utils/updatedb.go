package utils

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func AddBookInfo(title string, pageNumber int, author string, db *sql.DB, err error) {
	checkError(err)

	insertSQL := `INSERT INTO books (title, pageNumber, author) VALUES (?, ?, ?)`
	_, err = db.Exec(insertSQL, title, pageNumber, author)
	checkError(err)

}

func AddMovieInfo(title string, db *sql.DB, err error) {
	checkError(err)

	insertSQL := `INSERT INTO movies (title) VALUES (?)`
	_, err = db.Exec(insertSQL, title)
	checkError(err)

}

func AddVideoGameInfo(title string, db *sql.DB, err error) {
	checkError(err)

	insertSQL := `INSERT INTO videoGames (title) values (?)`
	_, err = db.Exec(insertSQL, title)
	checkError(err)
}

func AddUserInfo(name string, db *sql.DB, err error) {
	checkError(err)

	insertSQL := `INSERT INTO users (name) 	VALUES (?)`
	_, err = db.Exec(insertSQL, name)
	checkError(err)
}
