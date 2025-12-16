package models

type PaginationMetadata struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Count      int64 `json:"count"`
	TotalPages int64 `json:"totalPages"`
}
