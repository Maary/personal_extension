package models

import (
	"github.com/astaxie/beego/orm"
)

type Result struct {
	Id int
	Name string `json:"name"`
	Value string `json:"value"`
}

func (r *Result) TableName() string {
	return TableName("results")
}

func Insert(result *Result) error {
	om := orm.NewOrm()
	if _, err := om.Insert(result); err != nil {
		return err
	} else {
		return nil
	}
}
