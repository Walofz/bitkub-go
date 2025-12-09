package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3" // ไดรเวอร์ SQLite
)

// ตัวแปร Global สำหรับการเชื่อมต่อฐานข้อมูล
var DB *sql.DB

func InitDB() error {
	dbPath := "./db/bot_trades.db"

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating database directory: %w", err)
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS trades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		asset TEXT,
		operation TEXT, 
		amount_thb REAL,
		coin_amount REAL,
		price REAL,
		mode TEXT,
		deviation REAL,
		log_message TEXT
	);
	`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating trades table: %w", err)
	}

	fmt.Println("✅ Database initialized at:", dbPath)
	return nil
}

func LogTrade(
	asset string,
	operation string,
	amountTHB float64,
	coinAmount float64,
	price float64,
	mode string,
	deviation float64,
	logMessage string,
) {
	if DB == nil {
		fmt.Println("❌ Error: Database connection is nil. Cannot log trade.")
		return
	}

	insertSQL := `INSERT INTO trades (timestamp, asset, operation, amount_thb, coin_amount, price, mode, deviation, log_message) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(
		insertSQL,
		time.Now(),
		asset,
		operation,
		amountTHB,
		coinAmount,
		price,
		mode,
		deviation,
		logMessage,
	)

	if err != nil {
		fmt.Printf("❌ Error saving trade to DB: %v\n", err)
	}
}
