package utils

type PaginationQuery struct {
	Page     *int    `form:"page"`
	PageSize *int    `form:"page_size"`
	Sort     *string `form:"sort"`
	Order    *string `form:"order"`
	Type     *string `form:"type"`
	Search   *string `form:"search"`
}
