package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RequestData struct {
	Amount    float64 `json:"amount"`
	Rate      float64 `json:"rate"`
	StartDate string  `json:"start_date"`
}

func calculateSimpleInterest(amount, rate float64, months int) float64 {
	return amount + (amount * rate * float64(months) / 100)
}

func interestHandler(c *gin.Context) {
	var requestData RequestData

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

	// If the entered date is in the future, return an error
	if startDate.After(today) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is in the future. Please enter a valid past or present date."})
		return
	}

	// Ensure the date entered is not older than 2018
	if startDate.Year() < 2018 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a date after 2018."})
		return
	}

	// Calculate the number of months
	months := int(today.Sub(startDate).Hours() / 730) // Approximate months

	// If the difference is less than 30 days, count it as one month
	if today.Sub(startDate).Hours()/24 < 30 {
		months = 1
	}

	finalAmount := calculateSimpleInterest(requestData.Amount, requestData.Rate, months)

	c.JSON(http.StatusOK, gin.H{
		"months":      months,
		"finalAmount": finalAmount,
	})
}

func main() {
	r := gin.Default()

	// Enable CORS to allow frontend requests
	r.Use(cors.Default())

	// Load HTML templates from the "templates" folder
	r.LoadHTMLGlob("templates/*")

	// Serve index.html when visiting "/"
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API Endpoint
	r.POST("/calculate", interestHandler)

	// Start the server on port 8080
	r.Run(":8080")
}
