package main

import (
	"database/sql"
	"log"
	"os"

	"Libarymanagementsystem/tui"
	"Libarymanagementsystem/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	utils.InitDatabase()

	// Setting up the log
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetFlags(log.Lshortfile)
	log.SetOutput(file)

	// Setting up the db
	db, err := sql.Open("sqlite3", "projectdb.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Starting App
	tui.App(db)
	log.Println("------------------------------------------------")
}
