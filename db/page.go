package db

import (
	"math"
)

//PageObj 分页数据
type PageObj struct {
	Page     int   `json:"page"`
	Rows     int   `json:"rows"`
	Total    int64 `json:"total"`
	PageSize int   `json:"allpage"`
}

//setTotal 设置分页
func (p *PageObj) SetTotal(total int64) {
	if p.Rows < 1 { //如果未设置分页,自动设置成20
		p.Rows = 20
	}
	p.Total = total
	p.PageSize = int(math.Ceil(float64(total) / float64(p.Rows)))
}
