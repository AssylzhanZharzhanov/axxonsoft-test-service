package domain

// Default values.
const (
	DefaultPageNumber     = 0 // page starts from "0"
	DefaultPageSize       = 25
	DefaultIncludeDeleted = false
)

// PageRequest - request the page.
type PageRequest struct {
	Offset int // page offset (starts from 0)
	Size   int // page size
}

type Total int
