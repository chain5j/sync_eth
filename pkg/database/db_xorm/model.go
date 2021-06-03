// description: sync_eth
//
// @author: xwc1125
// @date: 2020/10/05
package db_xorm

import (
	"fmt"
)

const (
	ABLE int64 = iota
	UNABLE
	DELETE
	UNACTIVE
)

// Xorm
type XormModel struct {
	Id         int64  `xorm:"pk autoincr INT(11) notnull" json:"id" form:"id"`
	CreateTime int64  `xorm:"bigint(15) " json:"createTime"`
	CreateBy   string `xorm:"varchar(55)" json:"createBy"`
	UpdateTime int64  `xorm:"bigint(15) " json:"updateTime"`
	UpdateBy   string `xorm:"varchar(55)" json:"updateBy"`
	Remark     string `xorm:"varchar(500)" json:"remark"`
	Status     int64  `xorm:"int(8) notnull default(0)" json:"status"`
}

type XormModel1 struct {
	Id         int64  `xorm:"pk autoincr INT(11) notnull" json:"id" form:"id"`
	CreateTime int64  `xorm:"bigint(15) " json:"createTime"`
	UpdateTime int64  `xorm:"bigint(15) " json:"updateTime"`
	Remark     string `xorm:"varchar(500)" json:"remark"`
	Status     int64  `xorm:"int(8) notnull default(0)" json:"status"`
}

type XormModel2 struct {
	Id         int64 `xorm:"pk autoincr INT(11) notnull" json:"id" form:"id"`
	CreateTime int64 `xorm:"bigint(15) " json:"createTime"`
	UpdateTime int64 `xorm:"bigint(15) " json:"updateTime"`
	Status     int64 `xorm:"int(8) notnull default(0)" json:"status"`
}

type XormModelId struct {
	Id int64 `xorm:"pk autoincr INT(11) notnull" json:"id" form:"id"` // ID
}

func (m *XormModel) GetById(master MysqlConfig, id int64, bean interface{}) error {
	e := MasterEngine(master)
	_, err := e.Get(bean)
	if err != nil {
		return fmt.Errorf("model GetById: %w", err)
	}
	return nil
}

type PageWhereOrder struct {
	Order string
	Where string
	Value []interface{}
}

type OrderByCol struct {
	Column string
	Asc    bool
}

type PageResult struct {
	Page    *Paging     `json:"page"`
	Results interface{} `json:"results"`
}

type CursorResult struct {
	Results interface{} `json:"results"`
	Cursor  string      `json:"cursor"`
}

type Paging struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func (p *Paging) Offset() int {
	offset := 0
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Limit
	}
	return offset
}

func (p *Paging) TotalPage() int {
	if p.Total == 0 || p.Limit == 0 {
		return 0
	}
	totalPage := p.Total / p.Limit
	if p.Total%p.Limit > 0 {
		totalPage = totalPage + 1
	}
	return totalPage
}
