package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // pure Go, no GCC needed
)

func main() {
	fmt.Println("Welcome to ticket score engine!")

	// Open the SQLite database file
	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	//list tables in the database
	fmt.Println("Listing tables in database...")
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;`)
	if err != nil {
		log.Fatal("Query error:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Table:", tableName)
	}
}
