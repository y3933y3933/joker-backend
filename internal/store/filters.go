package store

type Filters struct {
	Page     int
	PageSize int
	SortBy   string
}

type Metadata struct {
	TotalCount  int `json:"totalCount,omitzero"`
	CurrentPage int `json:"currentPage,omitzero"`
	PageSize    int `json:"pageSize,omitzero"`
	FirstPage   int `json:"firstPage,omitzero"`
	LastPage    int `json:"lastPage,omitzero"`
}

func CalculateMetadata(totalCount, page, pageSize int) Metadata {
	if totalCount == 0 {
		return Metadata{}
	}

	return Metadata{
		TotalCount:  totalCount,
		PageSize:    pageSize,
		FirstPage:   1,
		LastPage:    (totalCount + pageSize - 1) / pageSize,
		CurrentPage: page,
	}
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
