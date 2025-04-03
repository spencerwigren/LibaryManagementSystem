package utils

import (
	"database/sql"
	"log"
	"time"
)

func UpdateEntryBook(db *sql.DB, titleUpdate string, pageNum int, authUpdate string, orgTitle string) {

	updateQuery := "UPDATE books SET title = ?, pageNumber = ?, author = ?, time = ? WHERE title = ?"
	row, err := db.Exec(updateQuery, titleUpdate, pageNum, authUpdate, time.Now(), orgTitle)
	if err != nil {
		log.Println("ERROR", err)
		return
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UPDATE %d rows(s)\n", rowsAffected)

}

func UpdateEntryMovie(db *sql.DB, titleUpdate string, ratingUpdate string, yearConverted int, orgTitle string) {
	updateQuery := "UPDATE movie SET title = ?, rating = ?, releaseYear = ?, time = ? WHERE title = ?"
	row, err := db.Exec(updateQuery, titleUpdate, ratingUpdate, yearConverted, time.Now(), orgTitle)
	if err != nil {
		log.Println("ERROR", err)
		return
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UPDATE %d rows(s)\n", rowsAffected)

}

func UpdateEntryGames(db *sql.DB, titleUpdate string, ratingUpdate string, yearConverted int, orgTitle string) {
	updateQuery := "UPDATE videoGames SET title = ?, rating = ?, releaseYear = ?, time = ? WHERE title = ?"
	row, err := db.Exec(updateQuery, titleUpdate, ratingUpdate, yearConverted, time.Now(), orgTitle)
	if err != nil {
		log.Println("ERROR", err)
		return
	}

	rowsAffected, err := row.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("UPDATE %d rows(s)\n", rowsAffected)

}
