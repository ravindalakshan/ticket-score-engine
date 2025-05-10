package domain

// TicketCategoryScore represents aggregated category score per ticket
type TicketCategoryScore struct {
	TicketID     int
	CategoryName string
	Score        float64
}
