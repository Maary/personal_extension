package Task

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"log"
	"personal_extension/lib/misc"
	"time"
)

type Task struct {
	URL string
	Tag
}

type Result struct {
	Data   string
	Status uint8
	Tag
}

type Tag struct {
	UUID string
	Type string
}

type Task_QueryParam struct {
	Limit     int
	Offset    int
	Order     string
	AscOrNo   bool
	Condition map[string]struct { //TODO
		ExOrNo bool
		Value interface{}
	}
}

//TODO
func (t *Task)updateTask(fields []string, newTasks *Task) map[string]interface{} {
	upParam := make(map[string]interface{})
	for _, f := range fields {
		switch f {
		case "URL":
			upParam["url"] = newTasks.URL
		}
	}
	upParam["updated"] = time.Now()
	return upParam
}

func (t *Task) TableName() string {
	return TableName("tasks")
}

func InsertTasks(ts []*Task) (count int64, err error) {
	return orm.NewOrm().InsertMulti(len(ts), ts)
}

func QueryTasks(param *Task_QueryParam) (ts []*Task, err error) {
	ts = make([]*Task, 0)
	query := orm.NewOrm().QueryTable(new(Task))
	order := "id"
	limit := 1000
	offset := 0
	if param.Order != "" {
		order = param.Order
	}
	if param.AscOrNo == false {
		order = fmt.Sprintf("-%s", order)
	}
	if param.Limit > 1000 {
		log.Println("query limit may too big")
	} else if param.Limit <= 0 {
		panic("query limit can not set the value that is less than 1")
	} else {
		limit = param.Limit
	}
	if param.Offset < 0 {
		panic("query offset can not set the value that is less than 0")
	} else {
		offset = param.Offset
	}
	query = query.Limit(limit).Offset(offset).OrderBy(order)
	if param.Condition != nil {
		for field, option := range param.Condition {
			if option.ExOrNo {
				query = query.Exclude(field, option.Value)
			} else {
				query = query.Filter(field, option.Value)
			}
		}
	}
	_, err = query.All(&ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func UpdateTasks(newTasks map[int]*Task) (count int64, err error) {
	om := orm.NewOrm()
	query := om.QueryTable(new(Task))
	for id, newTask := range newTasks {
		storedTask := new(Task)
		if err := query.Filter("id", id).One(storedTask); err != nil {
			return 0, err
		}
		ok, fields := misc.StructFac(newTasks, storedTask, "id", "created", "updated")
		if ok {
			upPrams := storedTask.updateTask(fields, newTask)
			_, err := query.Filter("id", id).Update(upPrams)
			if err != nil {
				return 0, err
			}
			count++
		}
	}
	return count, nil
}
