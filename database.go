package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3" // ไดรเวอร์ SQLite
)

// ตัวแปร Global สำหรับการเชื่อมต่อฐานข้อมูล
var DB *sql.DB

// InitDB: เชื่อมต่อและเตรียมฐานข้อมูล
func InitDB() error {
	var err error

	// 1. เชื่อมต่อฐานข้อมูล (ไฟล์ bot_trades.db จะถูกสร้างขึ้นถ้าไม่มี)
	DB, err = sql.Open("sqlite3", "./bot_trades.db")
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// 2. สร้างตารางถ้ายังไม่มี
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS trades (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATETIME,
		asset TEXT,
		operation TEXT, 
		amount_thb REAL,
		coin_amount REAL,
		price REAL,
		mode TEXT,         -- 'PRODUCTION' หรือ 'DRY_RUN'
		deviation REAL,
		log_message TEXT
	);
	`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating trades table: %w", err)
	}

	fmt.Println("✅ Database (SQLite) initialized successfully.")
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
