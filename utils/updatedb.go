package utils

import (
	"bufio"
	"database/sql"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// func AddBookInfo(title string, pageNumber int, author string, db *sql.DB, err error) {
// 	checkError(err)

// 	insertSQL := `INSERT INTO books (title, pageNumber, author) VALUES (?, ?, ?)`
// 	_, err = db.Exec(insertSQL, title, pageNumber, author)
// 	checkError(err)

// }

func AddBookInfo(db *sql.DB, err error) {
	checkError(err)

	reader := bufio.NewReader(os.Stdin)
	// var title string
	// var pageNumber int
	// var author string

	print("Input Title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	// TODO will need to make sure to check if an int
	print("Input Page Number: ")
	pageNumberString, _ := reader.ReadString('\n')
	pageNumberString = strings.TrimSpace(pageNumberString)
	pageNumber, err := strconv.Atoi(pageNumberString)
	checkError(err)

	print("Input Author Name: ")
	author, _ := reader.ReadString('\n')
	author = strings.TrimSpace(author)

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
