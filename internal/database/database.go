package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./webui.db")
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	err = createTables()
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	userTable := `CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        name TEXT,
        email TEXT UNIQUE,
        password TEXT,
        balance REAL,
        created_at DATETIME,
        updated_at DATETIME
    )`

	productTable := `CREATE TABLE IF NOT EXISTS products (
        id TEXT PRIMARY KEY,
        megnevezes TEXT,
        parameterek TEXT,
        price REAL,
        stock INTEGER
    )`

	orderTable := `CREATE TABLE IF NOT EXISTS orders (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        product_id TEXT,
        total_price REAL,
        status TEXT,
        created_at DATETIME,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (product_id) REFERENCES products(id)
    )`

	sessionTable := `CREATE TABLE IF NOT EXISTS sessions (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        expires_at DATETIME,
        FOREIGN KEY (user_id) REFERENCES users(id)
    )`

	tables := []string{userTable, productTable, orderTable, sessionTable}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
