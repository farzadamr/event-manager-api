package filter

import (
	"math"

	"github.com/farzadamr/event-manager-api/common"
)

func Paginate[TInput any, TOutout any](totalRows int64, items *[]TInput, pageNumber int, pageSize int64) (*PagedList[TOutout], error) {
	var rItems []TOutout

	rItems, err := common.TypeConverter[[]TOutout](items)
	if err != nil {
		return nil, err
	}
	return NewPagedList(&rItems, totalRows, pageNumber, pageSize), err
}

func NewPagedList[T any](items *[]T, count int64, pageNumber int, pageSize int64) *PagedList[T] {
	pl := &PagedList[T]{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalRows:  count,
		Items:      items,
	}
	pl.TotalPages = int(math.Ceil(float64(count) / float64(pageSize)))
	pl.HasNextPage = pl.PageNumber < pl.TotalPages
	pl.HasPreviosPage = pl.PageNumber > 1

	return pl
}

type PagedList[T any] struct {
	PageNumber     int   `json:"pageNumber"`
	PageSize       int64 `json:"pageSize"`
	TotalRows      int64 `json:"totalRows"`
	TotalPages     int   `json:"totalPage"`
	HasPreviosPage bool  `json:"hasPreviousPage"`
	HasNextPage    bool  `json:"hasNextPage"`
	Items          *[]T  `json:"items"`
}

type PaginationInput struct {
	PageSize   int `json:"pageSize"`
	PageNumber int `json:"pageNumber"`
}

func (p *PaginationInput) GetOffset() int {
	return (p.GetPageNumber() - 1) * p.GetPageSize()
}

func (p *PaginationInput) GetPageSize() int {
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	return p.PageSize
}

func (p *PaginationInput) GetPageNumber() int {
	if p.PageNumber == 0 {
		p.PageNumber = 1
	}
	return p.PageNumber
}
