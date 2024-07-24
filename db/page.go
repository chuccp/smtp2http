package db

type Data interface {
	*SMTP | any
}

type Page[T Data] struct {
	Total int `json:"total"`
	List  []T `json:"list"`
}

func ToPage[T Data](total int, list []T) *Page[T] {
	return &Page[T]{Total: total, List: list}
}
