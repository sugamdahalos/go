package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	dsn := "gouser:gopassword@tcp(localhost:3306)/goapp?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}
	defer db.Close()

	// Set connection pool parameters
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 3)

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database: ", err)
	}

	fmt.Println("Successfully connected to MySQL!")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}

	fmt.Println("Table 'users' created or already exists")

	result, err := db.Exec("INSERT INTO users (username, email) VALUES (?, ?)", "testuser", "test@example.com")
	if err != nil {
		log.Fatal("Failed to insert data: ", err)
	}

	id, _ := result.LastInsertId()
	fmt.Printf("Inserted user with ID: %d\n", id)

	var (
		username  string
		email     string
		createdAt time.Time
	)

	err = db.QueryRow("SELECT username, email, created_at FROM users WHERE id = ?", id).Scan(&username, &email, &createdAt)
	if err != nil {
		log.Fatal("Failed to query data: ", err)
	}

	fmt.Printf("User: %s, Email: %s, Created: %s\n", username, email, createdAt)
}
