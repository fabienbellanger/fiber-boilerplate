package requests

// Pagination request
type Pagination struct {
	Page  string `query:"p"`
	Limit string `query:"l"`
	Sorts string `query:"s"`
}
