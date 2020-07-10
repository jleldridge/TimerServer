package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
  host     = "localhost"
  port     = 5432
  dbname   = "GoTimerApp"
)

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(`host=%s port=%d dbname=%s sslmode=disable`,
		host, port, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}


func StartTimer(project, description string) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	start := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(start)

	_, err = db.Exec(`
		INSERT INTO timer_entries (project, description, start)
		VALUES ($1, $2, $3)
	`, project, description, start)
	return err
}

func StopTimer() error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stop := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(stop)

	// this is probably too slow in the long run
	_, err = db.Exec(`
		UPDATE timer_entries
		SET stop = $1
		WHERE stop IS NULL
	`, stop)

	return err
}

func QueryRows() (*sql.Rows, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res, err := db.Query(`SELECT * FROM timer_entries`)
	return res, err
}