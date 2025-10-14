package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	dsn := "postgres://user:password@localhost:5430/mydatabase?sslmode=disable"
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalln("Failed to connect:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Connected to PostgreSQL")

	user := User{Name: "Alice", Email: "alice@example.com", Balance: 100.0}
	err = InsertUser(db, user)
	if err != nil {
		log.Println("Insert error:", err)
	}

	users, err := GetAllUsers(db)
	if err != nil {
		log.Println("Select error:", err)
	}
	fmt.Println("All users:", users)

	
	err = TransferBalance(db, 1, 2, 25.0)
	if err != nil {
		log.Println("Transaction failed:", err)
	} else {
		fmt.Println("Transaction complete")
	}
}

func InsertUser(db *sqlx.DB, user User) error {
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`
	_, err := db.NamedExec(query, user)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	return users, err
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func TransferBalance(db *sqlx.DB, fromID int, toID int, amount float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return err
	}

	var balance float64
	err = tx.Get(&balance, "SELECT balance FROM users WHERE id=$1", fromID)
	if err != nil {
		return err
	}
	if balance < 0 {
		return fmt.Errorf("insufficient funds for user %d", fromID)
	}

	return tx.Commit()
}
