package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Request structure
type RequestData struct {
	Amount    float64 `json:"amount"`
	Rate      float64 `json:"rate"`
	StartDate string  `json:"start_date"`
}

// Function to calculate interest
func calculateSimpleInterest(amount, rate float64, months int) float64 {
	return amount + (amount * rate * float64(months) / 100)
}

func interestHandler(c *gin.Context) {
	var requestData RequestData

	// Validate JSON input
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Parse date in DD.MM.YYYY format
	startDate, err := time.Parse("02.01.2006", requestData.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (use DD.MM.YYYY)"})
		return
	}

	today := time.Now()

	// Validate that date is not in the future
	if startDate.After(today) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is in the future. Please enter a valid past or present date."})
		return
	}

	// Validate that date is not before 2018
	if startDate.Year() < 2018 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a date after 2018."})
		return
	}

	// Calculate the number of months properly
	yearDiff := today.Year() - startDate.Year()
	monthDiff := int(today.Month()) - int(startDate.Month())

	totalMonths := (yearDiff * 12) + monthDiff

	// If the day of the month is earlier in the current month, subtract one month
	if today.Day() < startDate.Day() {
		totalMonths--
	}

	// If the difference is less than a month, set it as one month
	if totalMonths < 1 {
		totalMonths = 1
	}

	finalAmount := calculateSimpleInterest(requestData.Amount, requestData.Rate, totalMonths)

	c.JSON(http.StatusOK, gin.H{
		"months":      totalMonths,
		"finalAmount": finalAmount,
	})
}

func main() {
	r := gin.Default()

	// Enable CORS for frontend
	r.Use(cors.Default())

	// Serve index.html
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API Endpoint
	r.POST("/calculate", interestHandler)

	// Get PORT from environment variable for deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run server
	r.Run(":" + port)
}
