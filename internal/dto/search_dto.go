package dto

type SearchPostDTO struct {
	Query    string `json:"query" form:"query" query:"query" validate:"required"`
	Page     int    `json:"page" form:"page" query:"page" validate:"required"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size" validate:"required"`
}
