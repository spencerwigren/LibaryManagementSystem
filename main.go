package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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

	// Used for debugging db
	directory()
	userCommandExe(db)
	utils.Query(db)

	// Starting App
	tui.App(db)
	log.Println("------------------------------------------------")
}

func userCommandExe(db *sql.DB) {
	check, userCommand := userCommandInput()

	if check {
		switch userCommand {
		// case 1:
		// 	utils.AddBookInfo("testBook", 125, "Admin", db)
		// 	// utils.AddBookInfo(db)
		// case 2:
		// 	utils.AddMovieInfo(db)
		// case 3:
		// 	utils.AddVideoGameInfo(db)
		case 4:
			// utils.AddUserInfo(db)
			utils.Query(db)
		}
	} else {
		println("Input not valid")
	}
}

func userCommandInput() (bool, int64) {
	// This is for the commandline part
	// REMOVE AFTER TUI is built
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

func directory() {
	// this func is for showing commands in the terminal
	// TODO: will replace with a faq once TUI uis build

	fmt.Println(`
	========================
	Commands for App

	Add Book: 1
	Add Movie: 2
	Add VideoGame: 3

	Quary All Items: 4

	Quit: 5
	========================
	`)
}
