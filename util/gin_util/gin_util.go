package gin_util

type PaginationParams struct {
	Page  int `form:"page" json:"page" validate:"gte=0"`
	Limit int `form:"limit" json:"limit" validate:"gte=0"`
}
