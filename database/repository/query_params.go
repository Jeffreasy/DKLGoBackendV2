package repository

// QueryParams bevat parameters voor database queries zoals paginering, sortering en filtering
type QueryParams struct {
	Page      int                    `json:"page"`
	PageSize  int                    `json:"page_size"`
	SortField string                 `json:"sort_field"`
	SortOrder string                 `json:"sort_order"`
	Filters   map[string]interface{} `json:"filters"`
	Search    string                 `json:"search"`
}

// NewQueryParams maakt een nieuwe QueryParams instantie met standaardwaarden
func NewQueryParams() *QueryParams {
	return &QueryParams{
		Page:      1,
		PageSize:  10,
		SortField: "created_at",
		SortOrder: "desc",
		Filters:   make(map[string]interface{}),
	}
}

// WithPage stelt de pagina in
func (q *QueryParams) WithPage(page int) *QueryParams {
	q.Page = page
	return q
}

// WithPageSize stelt de paginagrootte in
func (q *QueryParams) WithPageSize(pageSize int) *QueryParams {
	q.PageSize = pageSize
	return q
}

// WithSort stelt het sorteerveld en de sorteerrichting in
func (q *QueryParams) WithSort(field, order string) *QueryParams {
	q.SortField = field
	q.SortOrder = order
	return q
}

// WithFilter voegt een filter toe
func (q *QueryParams) WithFilter(key string, value interface{}) *QueryParams {
	q.Filters[key] = value
	return q
}

// WithSearch stelt de zoekterm in
func (q *QueryParams) WithSearch(search string) *QueryParams {
	q.Search = search
	return q
}

// GetOffset berekent de offset voor paginering
func (q *QueryParams) GetOffset() int {
	return (q.Page - 1) * q.PageSize
}

// GetLimit geeft de limiet voor paginering
func (q *QueryParams) GetLimit() int {
	return q.PageSize
}
