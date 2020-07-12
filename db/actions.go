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


func StartTimer(db *sql.DB, email, project, description string) error {
	userId, err := GetUserID(db, email)
	if err != nil {
		return err
	}

	start := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(`
		INSERT INTO timer_entries (project, description, startTime, userId)
		VALUES ($1, $2, $3, $4)
	`, project, description, start, userId)
	return err
}

func StopTimer(db *sql.DB, email string) error {
	userId, err := GetUserID(db, email)
	if err != nil {
		return err
	}

	stop := time.Now().Format("2006-01-02 15:04:05")

	_, err = db.Exec(`
		UPDATE timer_entries
		SET stopTime = $1
		WHERE stopTime IS NULL
			AND userId = $2
	`, stop, userId)

	return err
}

func GetHashedPassword(db *sql.DB, email string) (string, error) {
	res, err := db.Query(`
		SELECT hashed_passwords.hashedSaltedPassword
		FROM users
		RIGHT JOIN hashed_passwords
		ON users.id = hashed_passwords.id
		WHERE users.email = $1
	`, email)
	if err != nil {
		return "", err
	}
	defer res.Close()

	var hashedPassword string
	res.Next()
	res.Scan(&hashedPassword)

	return hashedPassword, nil
}

func CreateUser(db *sql.DB, email, hashedPassword string) error {
	existingUserRes, err := db.Query(`
		SELECT email
		FROM users
		WHERE email = $1
	`, email)
	if err != nil {
		return err;
	}
	defer existingUserRes.Close()

	if existingUserRes.Next() {
		return fmt.Errorf("User %s already exists.", email)
	}

	newUserRes, err := db.Query(`
		INSERT INTO users (email)
		VALUES ($1)
		RETURNING id
	`, email)
	if err != nil {
		return err;
	}
	defer newUserRes.Close()

	if !newUserRes.Next() {
		return fmt.Errorf("Failed to create user %s", email)
	}

	var userId string
	newUserRes.Scan(&userId)

	db.Exec(`
		INSERT INTO hashed_passwords (id, hashedSaltedPassword)
		VALUES ($1, $2)
	`, userId, hashedPassword)

	return nil
}

func UpdateUserPassword(db *sql.DB, email, newHashedPassword string) error {
	userId, err := GetUserID(db, email)
	if err != nil {
		return fmt.Errorf("user %s not found", email)
	}

	db.Exec(`
		UPDATE hashed_passwords
		set hashedSaltedPassword = $1
		WHERE id = $2
	`, newHashedPassword, userId)

	return nil
}


func GetUserID(db *sql.DB, email string) (int, error){
	res, err := db.Query(`
		SELECT id 
		FROM users
		WHERE email = $1
	`, email)
	if err != nil {
		return -1, err
	}
	defer res.Close()

	if !res.Next() {
		return -1, fmt.Errorf("user %s not found", email)
	}

	var id int
	res.Scan(&id)
	return id, nil
}