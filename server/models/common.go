package models

type PaginationResponse struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`      // 总记录数
	Page       int   `json:"page"`       // 当前页
	PageSize   int   `json:"pageSize"`   // 每页大小
	TotalPages int   `json:"totalPages"` // 总页数
}
