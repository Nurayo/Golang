package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	connStr := "host=localhost port=5433 user=nuray password=12345 dbname=practice4 sslmode=disable"

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("error connecting to DB:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Connected to PostgreSQL")

	fmt.Println("\n--- Adding new users ---")
	user1 := User{Name: "Gaukhar", Email: "gaukhar@kbtu.kz", Balance: 1000.0}
	if err := InsertUser(db, user1); err != nil {
		log.Println("InsertUser user1:", err)
	} else {
		fmt.Println("user1 added to DB")
	}

	user2 := User{Name: "Nuray", Email: "nuray@kbtu.kz", Balance: 500.0}
	if err := InsertUser(db, user2); err != nil {
		log.Println("InsertUser user2:", err)
	} else {
		fmt.Println("user2 added to DB")
	}

	fmt.Println("\n--- All users ---")
	users, err := GetAllUsers(db)
	if err != nil {
		log.Fatal("get users:", err)
	}
	for _, u := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s, Balance: %.2f\n", u.ID, u.Name, u.Email, u.Balance)
	}

	fmt.Println("\n--- Get user by ID (1) ---")
	user, err := GetUserByID(db, 1)
	if err != nil {
		log.Println("get by id:", err)
	} else {
		fmt.Printf("found: %s (Email: %s, Balance: %.2f)\n", user.Name, user.Email, user.Balance)
	}

	fmt.Println("\n--- Transfer: ID 1 → ID 2, amount: 200 ---")
	if err := TransferBalance(db, 1, 2, 200.0); err != nil {
		log.Println("transfer error:", err)
	} else {
		fmt.Println("transaction succeeded")
	}

	fmt.Println("\n--- Balances after transfer ---")
	users, _ = GetAllUsers(db)
	for _, u := range users {
		fmt.Printf("ID: %d, Name: %s, Balance: %.2f\n", u.ID, u.Name, u.Balance)
	}

	fmt.Println("\n--- Try to send more than available ---")
	if err := TransferBalance(db, 1, 2, 10000.0); err != nil {
		fmt.Println("expected transaction error:", err)
	}
}

func InsertUser(db *sqlx.DB, user User) error {
	const q = `
		INSERT INTO users (name, email, balance)
		VALUES (:name, :email, :balance)
	`
	_, err := db.NamedExec(q, user)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	const q = `SELECT id, name, email, balance FROM users ORDER BY id`
	if err := db.Select(&users, q); err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var u User
	const q = `SELECT id, name, email, balance FROM users WHERE id = $1`
	if err := db.Get(&u, q, id); err != nil {
		return User{}, err
	}
	return u, nil
}

func TransferBalance(db *sqlx.DB, fromID, toID int, amount float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot start transaction: %w", err)
	}

	var senderBalance float64
	if err := tx.Get(&senderBalance, `SELECT balance FROM users WHERE id = $1`, fromID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("sender not found: %w", err)
	}
	if senderBalance < amount {
		_ = tx.Rollback()
		return fmt.Errorf("insufficient funds: have %.2f need %.2f", senderBalance, amount)
	}

	var receiverExists bool
	if err := tx.Get(&receiverExists, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, toID); err != nil || !receiverExists {
		_ = tx.Rollback()
		return fmt.Errorf("receiver not found")
	}

	if _, err := tx.Exec(`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, fromID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("debit failed: %w", err)
	}

	if _, err := tx.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, toID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("credit failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}
