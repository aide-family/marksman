package shared

type PageRequest struct {
	Page     int32
	PageSize int32
	total    int64
}

func NewPageRequest(page, pageSize int32) *PageRequest {
	return &PageRequest{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *PageRequest) Offset() int {
	return int((p.Page - 1) * p.PageSize)
}

func (p *PageRequest) Limit() int {
	return int(p.PageSize)
}

func (p *PageRequest) WithTotal(total int64) *PageRequest {
	p.total = total
	return p
}

func (p *PageRequest) Total() int64 {
	return p.total
}

type Page[T any] struct {
	Items    []T
	Total    int64
	Page     int32
	PageSize int32
}

func NewPage[T any](req *PageRequest, items []T) *Page[T] {
	return &Page[T]{
		Items:    items,
		Total:    req.total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
}

