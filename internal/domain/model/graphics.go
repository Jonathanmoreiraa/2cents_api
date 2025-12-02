package model

type GraphicMonthTotal struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

type CategoryCount struct {
	Category float64 `json:"category" gorm:"column:category"`
	Name     string  `json:"name" gorm:"column:nome"`
}

type LastMetrics struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
