package common

type PageCondition[T any] struct {
	Page      Page `json:"page"`
	Condition T    `json:"condition"`
}

type PageResult[T any] struct {
	Records []T `json:"records"`
	Total   int `json:"total"`
}

type Page struct {
	Current int `json:"current"`
	Size    int `json:"size"`
}

type CountResult struct {
	Count int `json:"count" db:"count"`
}

type GlobalContidion struct {
	Condition string `json:"condition"`
}
