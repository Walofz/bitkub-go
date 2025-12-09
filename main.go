package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	if err := InitDB(); err != nil {
		fmt.Printf("Fatal error during DB initialization: %v\n", err)
		return
	}
	defer DB.Close()
	go StartBotLoop()
	r := gin.Default()

	r.Static("/static", "./web")
	r.LoadHTMLFiles("web/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/api/status", func(c *gin.Context) {
		summary := CalculatePortfolio()
		ConfigMutex.RLock()
		mode := "PRODUCTION"
		if IsDryRun {
			mode = "DRY_RUN"
		}
		ConfigMutex.RUnlock()

		c.JSON(http.StatusOK, gin.H{
			"status":      "Running",
			"mode":        mode,
			"last_run":    time.Now().Format("15:04:05"),
			"eth_price":   RoundFloat(LatestEthPrice, 2),
			"total_value": RoundFloat(summary.TotalValue, 2),
			"roi":         RoundFloat(summary.ROI, 2),
			"portfolio":   summary.Portfolio,
		})
	})

	r.GET("/api/history", func(c *gin.Context) {
		trades, err := GetProductionTrades(50)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"trades": trades,
		})
	})

	r.POST("/api/mode/:mode", func(c *gin.Context) {
		newMode := c.Param("mode")
		ConfigMutex.Lock()
		switch newMode {
		case "dry":
			IsDryRun = true
		case "prod":
			IsDryRun = false
		}
		ConfigMutex.Unlock()

		c.Redirect(http.StatusFound, "/api/status")
	})

	r.Run(":8888")
}
