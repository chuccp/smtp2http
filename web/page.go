package web

type Page struct {
	PageNo   int
	PageSize int
}
type PageAble struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}
