package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Result structure to send calculated results to the template
type Result struct {
	MonthsPassed int
	TotalDue     float64
	Error        string
}

// Function to parse flexible date formats
func parseFlexibleDate(inputDate string) (time.Time, error) {
	// List of possible date formats
	formats := []string{
		"2.1.06", "02.01.2006", // 3.2.21, 03.02.2021
		"2-1-06", "02-01-2006", // 3-2-21, 03-02-2021
		"2/1/06", "02/01/2006", // 3/2/21, 03/02/2021
		"2 1 06", "02 01 2006", // 3 2 21, 03 02 2021
		"2.1.06", "02.01.2006", // 12.2.22, 12.02.2022, etc
		"02.02.2022", "12.02.2022", // Add more variations
		"2.2.2022", "12.2.2022", // More formats like 12.2.22
		"2.2.22", "12.02.22", // Another common format
	}

	inputDate = strings.TrimSpace(inputDate)

	for _, format := range formats {
		parsedDate, err := time.Parse(format, inputDate)
		if err == nil {
			return parsedDate, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format. Please use formats like '3.2.21' or '03-02-2021'")
}

// Function to calculate the months passed (including 1 extra day as a new month)
func calculateMonthsPassed(receiptDate time.Time) int {
	// Get current date
	now := time.Now()

	// Calculate the months difference
	months := (now.Year()-receiptDate.Year())*12 + int(now.Month()) - int(receiptDate.Month())

	// If there's an extra day after the full month, count it as a new month
	if now.Day() > receiptDate.Day() {
		months++
	}

	return months
}

// Handler for the interest calculator
func interestCalculatorHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := Result{}

	if r.Method == http.MethodPost {
		// Get the receipt date and amount from the form
		receiptDateInput := r.FormValue("receipt_date")
		amountInput := r.FormValue("amount")

		// Parse the date
		receiptDate, err := parseFlexibleDate(receiptDateInput)
		if err != nil {
			data.Error = err.Error()
			tmpl.Execute(w, data)
			return
		}

		// Parse the amount
		amount, err := strconv.ParseFloat(amountInput, 64)
		if err != nil {
			data.Error = "Invalid amount. Please enter a valid number."
			tmpl.Execute(w, data)
			return
		}

		// Calculate months passed using the new logic
		monthsPassed := calculateMonthsPassed(receiptDate)

		// Simple Interest: 3% per month
		interestRate := 0.03
		totalInterest := float64(monthsPassed) * interestRate * amount

		// Calculate total due with simple interest
		totalDue := amount + totalInterest

		// Populate the results
		data.MonthsPassed = monthsPassed
		data.TotalDue = totalDue // No rounding required for simple interest
	}

	// Render the template with the calculated data or error
	tmpl.Execute(w, data)
}

func main() {
	// Serve static files (HTML templates)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	// Handle the main page
	http.HandleFunc("/", interestCalculatorHandler)

	// Start the server
	fmt.Println("Server is running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
