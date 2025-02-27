package utils

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func Query(db *sql.DB) {

	quaryDirectory()
	userQueryCommandExe(db)

}

func queryBooks(db *sql.DB) {
	// This will get all books
	// TODO: have user pick which books when GUI is build

	bookRows, err := db.Query("SELECT bookId, title, pageNumber, author FROM books;")
	checkError(err)
	defer bookRows.Close()

	fmt.Println("\nBooks:")
	for bookRows.Next() {
		var bookId int
		var title string
		var pageNumber int
		var author string

		err = bookRows.Scan(&bookId, &title, &pageNumber, &author)
		checkError(err)

		fmt.Printf("Book ID: %d, Title: %s, Page Number: %d, Author: %s\n", bookId, title, pageNumber, author)
	}

	if err = bookRows.Err(); err != nil {
		log.Fatal(err)
	}

}

func queryMovies(db *sql.DB) {
	// This will get all Movies
	// TODO: have user pick which movie when GUI is build
	movieRows, err := db.Query("SELECT movieId, title FROM movies")
	checkError(err)
	defer movieRows.Close()

	fmt.Println("\nMovies:")
	for movieRows.Next() {
		var movieId int
		var title string

		err = movieRows.Scan(&movieId, &title)
		checkError(err)

		fmt.Printf("Movie ID: %d, Title: %s\n", movieId, title)

	}

	if err = movieRows.Err(); err != nil {
		log.Fatal(err)
	}

}

func queryGames(db *sql.DB) {
	// This will get all Video Games
	// TODO: have user pick which Video Games when GUI is build
	videoGameRows, err := db.Query("SELECT videoGameId, title FROM videoGames")
	checkError(err)
	defer videoGameRows.Close()

	fmt.Println("\nVideo Games:")
	for videoGameRows.Next() {
		var videoGameId int
		var title string

		err = videoGameRows.Scan(&videoGameId, &title)
		checkError(err)

		fmt.Printf("Video Games ID: %d, Title: %s\n", videoGameId, title)
	}

}

func queryUser(db *sql.DB) {
	// This will get all Users
	// TODO: have user pick which Users when GUI is build
	userRows, err := db.Query("SELECT id, name FROM users;")
	checkError(err)
	defer userRows.Close()

	fmt.Println("\nUsers:")
	for userRows.Next() {
		var id int
		var name string

		err = userRows.Scan(&id, &name)
		checkError(err)

		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	if err = userRows.Err(); err != nil {
		log.Fatal(err)
	}
}

func queryAll(db *sql.DB) {
	// This will get all Info
	queryBooks(db)
	queryMovies(db)
	queryGames(db)
	queryUser(db)
}

// This is temp for termial use only along with userQueryCommandInput
func userQueryCommandExe(db *sql.DB) {
	check, userCommand := userQueryCommandInput()

	if check {
		switch userCommand {
		case 1:
			queryBooks(db)
		case 2:
			queryMovies(db)
		case 3:
			queryGames(db)
		case 4:
			queryUser(db)
		case 5:
			queryAll(db)
		}
	} else {
		println("Input not valid")
	}
}

func userQueryCommandInput() (bool, int64) {
	// This is for the commandline part
	// REMOVE AFTER GUI is built
	var command string
	commandList := [5]int64{1, 2, 3, 4, 5}
	fmt.Print("> ")
	fmt.Scanln(&command)

	input, err := strconv.ParseInt(command, 10, 64) // Base 10, 64-bit integer
	if err != nil {
		println("Not a valid input")
	} else {
		for _, value := range commandList {
			if value == input {
				return true, input
			}
		}
	}
	return false, input
}

func quaryDirectory() {
	fmt.Println(`
	========================
	Commands for Quary Data
	========================

	See Books:     1
	See Movie:     2
	See VideoGame: 3
	See Users:     4

	Quary All:     5
	========================
	`)
}
