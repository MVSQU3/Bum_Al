package db

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDb() *sql.DB {
	connStr := "user=postgres password=postgres dbname=album_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur ouverture DB", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erreur de connexion DB", err)
	}

	fmt.Println("Connexion a la db réussite!")

	return db
}
