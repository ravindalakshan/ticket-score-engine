package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"ticket-score-engine/internal"
	"ticket-score-engine/internal/ui"

	_ "modernc.org/sqlite" // TODO: check GCC needed option is needed or not
)

func getYearDateRange(year, startMonth, endMonth int) (time.Time, time.Time) {
	start := time.Date(year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)

	// Calculate end of end month
	nextMonth := time.Month(endMonth) + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	end := time.Date(nextYear, nextMonth, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)

	return start, end
}

func main() {
	fmt.Println("Welcome to ticket score engine!")

	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	// ======= Category scores for past week
	start := time.Now().AddDate(0, 0, -7) // 7 days ago
	end := time.Now()

	scores, err := internal.GetCategoryScores(db, start, end)
	if err != nil {
		log.Fatalf("Failed to get scores: %v", err)
	}

	fmt.Println("Category Scores (Past Week):")
	for _, s := range scores {
		fmt.Printf("%s | %s | %.2f%% (%d ratings)\n", s.CategoryName, s.Date, s.Score, s.RatingCount)
	}

	// ========= Category score for given year stand and end
	fmt.Println("Category Scores (when date range exceeds):")
	rangeStart, rangeEnd := getYearDateRange(2020, 1, 2)

	weeklyScores, err := internal.GetCategoryScores(db, rangeStart, rangeEnd)
	if err != nil {
		log.Fatalf("Failed to get yearly scores: %v", err)
	}

	// Determine if weekly aggregation was used
	duration := rangeEnd.Sub(rangeStart)
	weekly := duration > 30*24*time.Hour

	// Transform and display
	uiData := ui.TransformForUI(weeklyScores, weekly)
	ui.PrintUITable(uiData, weekly)
}
