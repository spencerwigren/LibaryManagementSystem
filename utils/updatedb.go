package utils

import (
	"bufio"
	"database/sql"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func AddBookInfo(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

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

func AddMovieInfo(db *sql.DB) {

	reader := bufio.NewReader(os.Stdin)
	print("Input Title Name: ")
	movieTitle, _ := reader.ReadString('\n')
	movieTitle = strings.TrimSpace(movieTitle)

	insertSQL := `INSERT INTO movies (title) VALUES (?)`
	_, err := db.Exec(insertSQL, movieTitle)
	checkError(err)

}

func AddVideoGameInfo(db *sql.DB) {

	reader := bufio.NewReader(os.Stdin)

	print("Input Video Game Title: ")
	gameTitle, _ := reader.ReadString('\n')
	gameTitle = strings.TrimSpace(gameTitle)

	insertSQL := `INSERT INTO videoGames (title) values (?)`
	_, err := db.Exec(insertSQL, gameTitle)
	checkError(err)
}

func AddUserInfo(db *sql.DB) {

	reader := bufio.NewReader(os.Stdin)

	print("Input User Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	insertSQL := `INSERT INTO users (name) 	VALUES (?)`
	_, err := db.Exec(insertSQL, name)
	checkError(err)
}
