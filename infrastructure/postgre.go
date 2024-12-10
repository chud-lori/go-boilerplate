package infrastructure

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgreDB() *sql.DB {

	var (
		host     = os.Getenv("PSQL_HOST")
		port     = 5432
		user     = os.Getenv("PSQL_USER")
		password = os.Getenv("PSQL_PASSWORD")
		dbname   = os.Getenv("DB_NAME")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}
