package ui

import (
	"fmt"
	"sort"
	"strings"
	"ticket-score-engine/internal/scoring"
)

// CategoryScoreSummary represents aggregated scores for UI presentation
type CategoryScoreSummary struct {
	CategoryName string
	RatingCount  int
	DateScores   map[string]float64 // Date/Week -> Score
	FinalScore   float64
}

// TransformForUI converts raw scores to UI-friendly format
func TransformForUI(scores []scoring.CategoryScore, weekly bool) []CategoryScoreSummary {
	// Group by category
	categoryMap := make(map[string]*CategoryScoreSummary)
	allDates := make(map[string]bool)

	for _, score := range scores {
		if _, exists := categoryMap[score.CategoryName]; !exists {
			categoryMap[score.CategoryName] = &CategoryScoreSummary{
				CategoryName: score.CategoryName,
				DateScores:   make(map[string]float64),
			}
		}

		summary := categoryMap[score.CategoryName]
		summary.RatingCount += score.RatingCount
		summary.DateScores[score.Date] = score.Score
		allDates[score.Date] = true

		// Calculate weighted final score
		summary.FinalScore = (summary.FinalScore*float64(summary.RatingCount-score.RatingCount) +
			score.Score*float64(score.RatingCount)) / float64(summary.RatingCount)
	}

	// Convert map to slice
	var results []CategoryScoreSummary
	for _, summary := range categoryMap {
		results = append(results, *summary)
	}

	return results
}

// PrintUITable displays scores in a formatted table
func PrintUITable(data []CategoryScoreSummary, weekly bool) {
	// First collect all unique dates/weeks
	var allDates []string
	dateSet := make(map[string]bool)

	for _, item := range data {
		for date := range item.DateScores {
			if !dateSet[date] {
				dateSet[date] = true
				allDates = append(allDates, date)
			}
		}
	}
	sort.Strings(allDates)

	// Print header
	printHeader(allDates, weekly)

	// Print separator
	printSeparator(allDates)

	// Print data
	printDataRows(data, allDates)
}

// Helper functions for table printing
func printHeader(allDates []string, weekly bool) {
	fmt.Printf("| %-20s | %-8s", "Category", "Ratings")
	for _, date := range allDates {
		if weekly {
			parts := strings.Split(date, "-")  // Split year-week
			fmt.Printf(" | Week %s", parts[1]) // Just show week number
		} else {
			fmt.Printf(" | %s", date)
		}
	}
	fmt.Printf(" | Score |\n")
}

func printSeparator(allDates []string) {
	fmt.Printf("|%s|%s", strings.Repeat("-", 22), strings.Repeat("-", 10))
	for range allDates {
		fmt.Printf("|%s", strings.Repeat("-", 10))
	}
	fmt.Printf("|%s|\n", strings.Repeat("-", 8))
}

func printDataRows(data []CategoryScoreSummary, allDates []string) {
	for _, item := range data {
		fmt.Printf("| %-20s | %-8d", item.CategoryName, item.RatingCount)
		for _, date := range allDates {
			if score, exists := item.DateScores[date]; exists {
				fmt.Printf(" | %-6.1f%%", score)
			} else {
				fmt.Printf(" | %-6s", "N/A")
			}
		}
		fmt.Printf(" | %-5.1f%% |\n", item.FinalScore)
	}
}
