package models

const (
	DefaultPageNum        = 1
	DefaultPageSize       = 10
	DefaultTechnologyType = 0
	DefaultSortType       = 0
)

type SearchArgs struct {
	Search   string
	PageNum  int
	PageSize int
}
