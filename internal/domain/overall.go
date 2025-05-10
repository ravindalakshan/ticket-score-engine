package domain

// OverallScoreResult represents the result of an overall score calculation
type OverallScoreResult struct {
	Score       float64
	RatingCount int
}

// PeriodComparisonResult represents the comparison between two time periods
type PeriodComparisonResult struct {
	PercentageChange float64
	CurrentScore     float64
	PreviousScore    float64
	CurrentCount     int
	PreviousCount    int
}
