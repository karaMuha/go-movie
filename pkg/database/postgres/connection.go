package postgres

import (
	"database/sql"
	"log"
	"time"
)

func initDatabase(dbDriver string, dbConnection string) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbConnection)
	if err != nil {
		log.Printf("Error while initializing database connection %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Error while validation database connection: %v", err)
		return nil, err
	}

	return db, nil
}

func ConnectToDb(dbDriver string, dbConnectionString string) (*sql.DB, error) {
	var count int

	for {
		dbHandler, err := initDatabase(dbDriver, dbConnectionString)

		if err == nil {
			return dbHandler, nil
		}

		log.Println("Postgres container not yet ready...")
		count++
		log.Println(count)

		if count > 10 {
			log.Println("Too many retries")
			return nil, err
		}

		log.Println("Backing off for five seconds...")
		time.Sleep(5 * time.Second)
	}
}
