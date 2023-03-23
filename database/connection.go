package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Db *sql.DB

func InitDB() {

	password := os.Getenv("PASSWORD")
	if password == "" {
		log.Fatal("PASSWORD environment variable not set")
	}

	db, err := sql.Open("postgres", fmt.Sprintf("user=postgres password=%s host=db.xozmizrgnjdmofvmjvvs.supabase.co port=5432 dbname=postgres", password))

	if err != nil {
		panic(err)
	}

	log.Printf("connected to database")

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	Db = db
}

func CloseDB() error {
	return Db.Close()
}
