package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDb() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Erreur ouverture DB", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Erreur de connexion DB", err)
	}

	fmt.Println("Connexion a la db réussite!")

	return db
}
