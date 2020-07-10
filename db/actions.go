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
		INSERT INTO timer_entries (project, description, startTime)
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
		SET stopTime = $1
		WHERE stopTime IS NULL
	`, stop)

	return err
}

func GetHashedPassword(email string) (string, error) {
	db, err := Connect()
	if err != nil {
		return "", err
	}
	defer db.Close()

	res, err := db.Query(`
		SELECT hashed_passwords.hashedSaltedPassword
		FROM users
		RIGHT JOIN hashed_passwords
		ON users.id = hashed_passwords.id
		WHERE users.email = $1
	`, email)

	var hashedPassword string
	res.Next()
	res.Scan(&hashedPassword)

	return hashedPassword, nil
}

func CreateUser(email, hashedPassword string) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Query(`
		SELECT email
		FROM users
		WHERE email = $1
	`, email)
	if err != nil {
		return err;
	}

	if res.Next() {
		return fmt.Errorf("User %s already exists.", email)
	}

	res, err = db.Query(`
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id
	`, email)
	if err != nil {
		return err;
	}
	if !res.Next() {
		return fmt.Errorf("Failed to create user %s", email)
	}

	var userId string
	res.Scan(&userId)
	fmt.Println(userId)

	db.Exec(`
		INSERT INTO hashed_passwords (id, hashedSaltedPassword)
		VALUES ($1, $2)
	`, userId, hashedPassword)

	return nil
}

// func QueryRows() (*sql.Rows, error) {
// 	db, err := Connect()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	res, err := db.Query(`SELECT * FROM timer_entries`)
// 	return res, err
// }