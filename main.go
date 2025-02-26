package main

import (
	"database/sql"
	"fmt"
	"log"

	"Libarymanagementsystem/utils"

	_ "github.com/mattn/go-sqlite3"
)

// getting demo writen
// the terminal is going ot act as the UI for now

func main() {
	utils.InitDatabase()

	db, err := sql.Open("sqlite3", "projectdb.db")
	checkError(err)
	defer db.Close()

	utils.AddUserInfo("Admin", db, err)
	utils.AddBookInfo("testBook", 125, "Admin", db, err)
	utils.AddMovieInfo("testMoive", db, err)
	utils.AddVideoGameInfo("testVideoGame", db, err)
	queryDB(db, err)

}

func queryDB(db *sql.DB, err error) {
	checkError(err)
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

		fmt.Printf("Book ID: %d, Title: %s, Page Number %d, Author %s\n", bookId, title, pageNumber, author)
	}

	if err = bookRows.Err(); err != nil {
		log.Fatal(err)
	}

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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
