package dto

type ProgressResponse struct {
	Completed  int64   `json:"completed_lessons"`
	Total      int64   `json:"total_lessons"`
	Percentage float64 `json:"percentage"`
}
