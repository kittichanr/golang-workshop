package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func createTable(db *sql.DB) {
	query := `CREATE TABLE users (
		id INT AUTO_INCREMENT,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME,
		PRIMARY KEY (id)
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func insert(db *sql.DB) {
	var username string
	var password string
	fmt.Scan(&username)
	fmt.Scan(&password)

	createAt := time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createAt)
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}

func query(db *sql.DB) {
	var (
		id         int
		coursename string
		price      float64
		instructor string
	)

	for {
		var inputID int
		fmt.Scan(&inputID)

		query := "SELECT id, coursename, price, instructor FROM onlinecourse WHERE id = ?"
		if err := db.QueryRow(query, inputID).Scan(&id, &coursename, &price, &instructor); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, coursename, price, instructor)
	}
}

func delete(db *sql.DB) {
	var deleteId string
	fmt.Scan(&deleteId)

	result, err := db.Exec(`DELETE FROM users WHERE id = ?`, deleteId)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func main() {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		fmt.Println("failed to connect")
	} else {
		fmt.Println("connect successfully")
	}

	delete(db)
}
