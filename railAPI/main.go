package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/atanda0x/Metro-Rail-Api/dbutils"
)

func main() {
	db, err := sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("Driver creation failed!!!")
	}

	// Create tables
	dbutils.Initialise(db)
}
