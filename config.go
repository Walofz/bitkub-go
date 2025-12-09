package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var (
	APIKey            string
	APISecret         string
	IsDryRun          bool
	InitialInvestment float64
	Threshold         float64
	MinEthAmount      float64

	TargetAssets = map[string]float64{
		"THB": 50.0,
		"ETH": 50.0,
	}
)

var ConfigMutex sync.RWMutex

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("❌ Warning: No .env file found or failed to load. Using system environment/defaults.")
	}

	APIKey = os.Getenv("BITKUB_API_KEY")
	APISecret = os.Getenv("BITKUB_API_SECRET")

	if APIKey == "" || APISecret == "" {
		fmt.Println("❌ CRITICAL: BITKUB_API_KEY or SECRET is missing. Trading will fail.")
	}
	IsDryRun, _ = strconv.ParseBool(os.Getenv("BOT_IS_DRY_RUN"))

	if val, err := strconv.ParseFloat(os.Getenv("BOT_INITIAL_INVESTMENT"), 64); err == nil {
		InitialInvestment = val
	}

	if val, err := strconv.ParseFloat(os.Getenv("BOT_THRESHOLD"), 64); err == nil {
		Threshold = val
	}

	if val, err := strconv.ParseFloat(os.Getenv("BOT_MIN_ETH_AMOUNT"), 64); err == nil {
		MinEthAmount = val
	}

	fmt.Printf("✅ Config loaded. Mode: %s, Initial Inv: %.2f THB\n",
		func() string {
			if IsDryRun {
				return "DRY_RUN"
			} else {
				return "PRODUCTION"
			}
		}(),
		InitialInvestment)
}
