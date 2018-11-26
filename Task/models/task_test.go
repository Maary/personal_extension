package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"testing"
	"time"
	"personal_extension/lib/system"
	"personal_extension/lib/database"
)

func initSys() {
	config := system.Config{AppName: "", RunMode: system.DEV, ServerName: "", Filename: "app.conf", Dir: "../conf", AppConfigProvider: "ini"}
	system.InitSystem(&config)
	database.InitDatabase()
	Register()
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		panic(err)
	}
}

func TestInsertTasks(t *testing.T) {
	initSys()
	task := &Task{
		Id:1,
		Page:1,
		Offset:10,
		Content:"test content",
		Created:time.Now(),
		Type: "test",
		UUID: "1234uuid",
		Updated:time.Now(),
	}
	c, err := InsertTasks([]*Task {task})
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
}

func TestQueryTasks(t *testing.T) {
	initSys()
	param := new(Task_QueryParam)
	param.Limit = 100
	param.Condition = map[string]struct { //TODO
		ExOrNo bool
		Value  interface{}
	}{"offset": {ExOrNo: false,
		Value: 10,
	},
	}
	ts,err := QueryTasks(param)
	if err != nil {
		panic(err)
	}
	fmt.Println(ts)
}
