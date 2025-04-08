package utils

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func SearchTables(db *sql.DB, search string) ([]interface{}, string, error) {
	tableNames, err := fetchTableName(db)
	if err != nil {
		log.Printf("In SearchTables: %s\n", err)
	}

	// none, err := debuggingTable(db)
	// log.Printf("debuggingTable Output: %v, %v", none, err)

	// [1:] is to skip over the first index of the tableName
	// tableNames: [sqlite_sequence, users, books, movies, videoGames]
	// sqlite_sequence in tableNames is not one of the tables
	for i, tableName := range tableNames[1:] {
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
				return nil, "", err
			}

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
				//Returning True Values
				return valuePtrs, tableName, nil
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

	return nil, "", err
}

func checkTableColumn(db *sql.DB, tableName, columnName string) bool {
	// Checking table to see if the table as the columnName in it

	// query to get table
	log.Println("Table Name:", tableName)
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
		// Name here is the column title
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

	// going through the tables names
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Printf("Cannot Fetch Table Name: %s\n", err)
			return nil, err
		}
		// Saving to list
		names = append(names, name)
	}

	return names, nil

}

func QueryAllEntry(db *sql.DB) []any {
	tableNames, err := fetchTableName(db)
	if err != nil {
		log.Printf("In SearchTables: %s\n", err)
	}

	var rowsEntries []any

	for i, tableName := range tableNames[1:] {
		query := fmt.Sprintf("SELECT * FROM %s", tableName)

		rows, err := db.Query(query)
		if err != nil {
			log.Printf("CAN NOT QUERY DB TABLE: %s", tableName)
			continue
		}
		defer rows.Close()

		for rows.Next() {
			valuePtrs, err := processing(rows)
			if err != nil {
				log.Printf("IN queryAllEntry rows.Next() processing rows: %s\n", err)
				return nil
			}
			log.Printf("TABLE ROW ENTRIES[%d]: %s", i, valuePtrs)
			rowsEntries = append(rowsEntries, valuePtrs)

		}

	}

	return rowsEntries
}

// func debuggingTable(db *sql.DB) (any, error) {

// 	temp, err := db.Query("PRAGMA table_info(books);")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer temp.Close()

// 	log.Println("Table Structure:")
// 	for temp.Next() {
// 		var cid int
// 		var name string
// 		var ctype string
// 		var notnull int
// 		var dflt_value sql.NullString
// 		var pk int

// 		err := temp.Scan(&cid, &name, &ctype, &notnull, &dflt_value, &pk)
// 		if err != nil {
// 			return nil, err
// 		}

// 		log.Printf("Column ID: %d, Name: %s, Type: %s, Not Null: %d, Default: %v, Primary Key: %d\n",
// 			cid, name, ctype, notnull, dflt_value.String, pk)
// 	}

// 	if err = temp.Err(); err != nil {
// 		return nil, err
// 	}

// 	return nil, err
// }

func QueryMostRecent(db *sql.DB) *sql.Rows {

	tableNames, err := fetchTableName(db)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return nil
	}

	if len(tableNames) == 0 {
		log.Println("ERROR: No Tables in tableNames")
		return nil
	}

	// tableNames[2:] skipping the first and second enty, not a table name using here
	// Getting all table names with the time column in it
	var queries []string
	// log.Println("TABLE NAME", tableNames[2:])
	for _, tableName := range tableNames[2:] {
		// checking if the table contains the column time
		// BY seeing if the tableName is empty
		if tableName == "" {
			continue
		}

		queries = append(queries, fmt.Sprintf("SELECT title, time FROM %s", tableName))
	}

	// Creating search query
	query := fmt.Sprintf(`
		SELECT * FROM (
			%s
		)
		ORDER BY time DESC
		LIMIT 5;`, strings.Join(queries, " UNION ALL "))

	rows, err := db.Query(query)
	if err != nil {
		log.Println("ERROR", err)
		return nil
	}

	return rows
}

func CheckExisting(db *sql.DB, existingEntry string) bool {
	/*
		Checks to see if an entry is already in DB
		returns true if entry is in db
		returns false if entry is not in db
	*/

	tableNames, err := fetchTableName(db)
	if err != nil {
		log.Printf("ERROR: %s", err)
	} else if len(tableNames) == 0 {
		log.Println("ERROR: NO Tables in TableNames")
	}

	// tableNames[2:] skipping first and seconds tables entry
	var queries []string
	for _, tableName := range tableNames[2:] {
		// Checking if tablName is empty
		if tableName == "" {
			continue
		}

		queries = append(queries, fmt.Sprintf("SELECT title FROM %s", tableName))
	}

	// Creating search query
	query := fmt.Sprintf(`
		SELECT * FROM (
			%s
		)
		ORDER BY title DESC;`, strings.Join(queries, " UNION ALL "))

	// Getting data from db
	rows, err := db.Query(query)
	if err != nil {
		log.Println("ERROR", err)
		// Err Escape
		return true
	}
	defer rows.Close()

	var title string

	// Checking if entry is in db
	for rows.Next() {
		if err := rows.Scan(&title); err != nil {
			log.Println("ERROR", err)
			continue
		}

		log.Printf("TITLE: %s", title)

		if title == existingEntry {
			// Found similar
			return true
		}

	}

	// Not Found - good
	return false
}

func CheckUserExisting(db *sql.DB, existingEntry string) bool {

	querySelect := "SELECT name FROM users"

	rows, err := db.Query(querySelect)
	if err != nil {
		log.Println("ERROR", err)
		return true
	}
	defer rows.Close()

	var name string

	for rows.Next() {
		if err := rows.Scan(&name); err != nil {
			log.Println("ERROR", err)
			continue
		}

		if name == existingEntry {
			return true
		}
	}

	return false
}
