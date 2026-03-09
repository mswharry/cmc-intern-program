package model

type StatisticsResponse struct {
	Total int `json:"total"`
	ByType map[string]int `json:"by_type"`
	ByStatus map[string]int `json:"by_status"`
}
type CountResponse struct {
	Count int `json:"count"`
	Filters map[string]string `json:"filters"`
}