package gcurd

import "errors"

type (
	WhereValue struct {
		Name  string `json:"name"`
		Op    Op     `json:"op"`
		Value any    `json:"value"`
	}

	KeyValue struct {
		Key   string `json:"key"`
		Value any    `json:"value"`
	}

	OrderBy struct {
		Direction string `json:"direction"`
		Filed     string `json:"filed"`
	}

	Pagination struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}

	Request struct {
		Where      []*WhereValue
		OrderBy    OrderBy
		Pagination Pagination
	}
)

const (
	DefaultPage     = 0
	DefaultPageSize = 10
)

// SimplePage calculate "from", "to" without total_counts
// "from" index start from 1
func (p *Pagination) SimplePage() (from int, to int) {
	if p.Page == 0 || p.PageSize == 0 {
		p.Page, p.PageSize = 1, DefaultPageSize
	}
	from = (p.Page-1)*p.PageSize + 1
	to = from + p.PageSize - 1
	return
}

// CalPage calculate "from", "to" with total_counts
// index start from 1
func (p *Pagination) CalPage(total int) (from int, to int) {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}

	if total == 0 || total < p.PageSize*(p.Page-1) {
		return
	}
	if total <= p.PageSize {
		return 1, total
	}
	from = (p.Page-1)*p.PageSize + 1
	if (total - from + 1) < p.PageSize {
		return from, total
	}
	return from, from + p.PageSize - 1
}

// VagueOffsetLimit calculate "offset", "limit" without total_counts
func (p *Pagination) VagueOffsetLimit() (offset int, limit int) {
	from, to := p.SimplePage()
	if to == 0 || from == 0 {
		return 0, 0
	}
	return from - 1, to - from + 1
}

// OffsetLimit calculate "offset" and "start" with total_counts
func (p *Pagination) OffsetLimit(total int) (offset int, limit int) {
	from, to := p.CalPage(total)
	if to == 0 || from == 0 {
		return 0, 0
	}
	return from - 1, to - from + 1
}

func (p *Pagination) Verify() error {
	if p.Page < 0 {
		return errors.New("page error")
	} else if p.Page == 0 {
		p.Page = DefaultPage
	}
	if p.PageSize < 0 {
		return errors.New("page size error")
	} else if p.PageSize == 0 {
		p.PageSize = DefaultPageSize
	}
	return nil
}
