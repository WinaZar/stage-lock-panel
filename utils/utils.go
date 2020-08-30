package utils

// LockData presents passwords for lock stage
type LockData struct {
	Code     string `json:"code,omitempty" valid:"alphanum,length(0|8)"`
	LockedBy string `json:"locked_by" valid:"length(1|100)"`
	Comment  string `json:"comment,omitempty" valid:"length(0|500)"`
}

// PaginationData presents pagination info
type PaginationData struct {
	Page       int `query:"page" json:"page" valid:"range(1|1000)"`
	PerPage    int `query:"per-page" json:"per_page" valid:"range(1|100)"`
	TotalItems int `json:"total,omitempty"`
}

// SortData presents sort column and sort desc
type SortData struct {
	SortBy    string `query:"sort-by" valid:"in(created_at|locked_by|action)"`
	SortOrder string `query:"sort-order" valid:"in(asc|desc)"`
}

// StandartJSONResponse implement standart json-response
type StandartJSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
