package entity

type Pagination struct {
	Page    uint32
	PerPage uint32
	Keyword string
	Total   uint32
}
