package utils

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func SearchTables(db *sql.DB, search string) ([]interface{}, error) {
	// TODO: may need to find a way to search outside of title, like page number or author

	tableNames, err := fetchTableName(db)
	if err != nil {
		log.Printf("In SearchTables: %s\n", err)
	}

	/*
		Debugging Tables
		temp, err := db.Query("PRAGMA table_info(books);")
		if err != nil {
			return nil, err
		}
		defer temp.Close()

		log.Println("Table Structure:")
		for temp.Next() {
			var cid int
			var name string
			var ctype string
			var notnull int
			var dflt_value sql.NullString
			var pk int

			err := temp.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
			if err != nil {
				return nil, err
			}

			log.Printf("Column ID: %d, Name: %s, Type: %s, Not Null: %d, Default: %v, Primary Key: %d\n",
				cid, name, ctype, notnull, dflt_value.String, pk)
		}

		if err = temp.Err(); err != nil {
			return nil, err
		}
	*/

	// [1:] is to skip over the first index of the tableName
	// tableNames: [sqlite_sequence users books movies videoGames]
	// sqlite_sequence in tableNames is not one of the tables
	for i, tableName := range tableNames[1:] {
		// TODO: The title is temp, will need to find a way to seach all fields
		if !checkTableColumn(db, tableName, "title") {
			log.Println("Not In Table", tableName)
			continue
		}

		// Query for searching the table
		query := fmt.Sprintf("SELECT * FROM %s WHERE title LIKE ?", tableName)

		// Debugging
		log.Printf("Number of loops: %d", i)
		log.Printf("Query Test: %s", query)
		log.Printf("Len of tableName, %d", len(tableName))

		// making search a wildcard
		searchWildCard := "%" + search + "%"
		rows, err := db.Query(query, searchWildCard)
		if err != nil {
			// To Till if an error happened
			log.Printf("IN FOR LOP SEACHING DB TABLES: %d in db %s", err, tableName)
			// return nil, err
			continue

		}
		defer rows.Close()

		// Searching rows in table
		for rows.Next() {
			valuePtrs, err := processing(rows)
			if err != nil {
				// Tell if an error happened
				log.Printf("IN SEARCH TABLES rows.Next() processing rows: %s\n", err)
				return nil, err
			}
			//TODO: make this less hardcoded
			// Value 0 is the id number
			// Value 1 is the title
			// Value 2 is the page number
			// Value 3 is the Author

			// converting valuePtrs[1] - title - into string for comparison
			strValue, ok := (*(valuePtrs[1].(*interface{}))).(string)
			if !ok {
				log.Println("Failed to convert valuePtrs[1] to str")
			} else {
				log.Printf("Converted value: %s", strValue)
			}

			// comparing the value with the user search input
			if strValue == search {
				log.Printf("Found Search: [%d]: %s", i, tableName)
				return valuePtrs, nil
			} else {
				// if not found, sent to log for debugging
				log.Printf("SEACH VALUE: %s", search)
				log.Printf("ValuePtrs: %s", valuePtrs[1])
				for i, ptr := range valuePtrs {
					log.Printf("ValuePtrs[%d]: %v", i, *ptr.(*interface{}))
				}
			}
		}
	}

	return nil, err
}

func checkTableColumn(db *sql.DB, tableName, columnName string) bool {
	// Checking table to see if the table as the columnName in it
	// If not checked err will close program

	// query to get table
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error Getting Table Info: %s, %d", tableName, err)
		return false
	}
	defer rows.Close()

	// Checking all rows for correct row name
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue sql.NullString

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			log.Println("Error scanning table info:", tableName, err)
			return false
		}

		// Found correct name
		if name == columnName {
			return true
		}
	}
	return false
}

func processing(row *sql.Rows) ([]interface{}, error) {
	// Checking rows
	columns, err := row.Columns()
	if err != nil {
		log.Printf("IN PROCESSING Couldn't get columns: %s\n", err)
		return nil, err
	}

	// getting the values of the columns
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Row processing
	if err := row.Scan(valuePtrs...); err != nil {
		log.Printf("IN ROW Processing: %s\n", err)
		return nil, err
	}

	return valuePtrs, nil
}

func fetchTableName(db *sql.DB) ([]string, error) {
	// Getting the Tables name in the DB
	// Idea from here:
	// Along for processing
	// https://search.brave.com/search?q=how+to+search+through+all+db+tables+in+go+in+sqlite&source=web&summary=1&conversation=209bf105ecdc6630ee6a42
	var names []string
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Printf("Cannot Get Table Name: %s\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Cannot Fetch Table Name: %s\n", err)
			return nil, err
		}
		names = append(names, name)
	}

	return names, nil

}

/*
Debugging whats in db in the termial
*/
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
