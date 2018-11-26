package client

import (
	"context"
	"encoding/json"
	"fmt"
	"personal_extension/Task/task_service"
	"testing"
	"time"
)

func InitConfigForTest() {
	InitConn([]string{"192.168.99.254:2379"}, "deepdraw", "task_server_service")
}

func TestInsertTasks(t *testing.T) {
	InitConfigForTest()
	in := new(task_rpc_config.Tasks)
	ts := map[string]interface{} {"Id": 2, "Page": 10, "Offset": 100, "Created": time.Now(), "Updated": time.Now(), "Type": "test", "UUID": "testuuid"}
	tss := []map[string]interface{} {ts}
	if jsonB, err := json.Marshal(tss); err != nil {
		panic(err)
	} else {
		in.Content = string(jsonB)
	}
	rsp, err := client.InsertTasks(context.Background(), in)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.ErrMessage)
}

type Task_QueryParam struct {
	Limit     int
	Offset    int
	Order     string
	AscOrNo   bool
	Condition map[string]struct { //TODO
		ExOrNo bool
		Value  interface{}
	}
}

func TestQueryTasks(t *testing.T) {
	InitConfigForTest()
	param := new(Task_QueryParam)
	param.Limit = 10
	jsonB, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	in := new(task_rpc_config.Params)
	in.ParamsStr = string(jsonB)
	rsp, err := client.QueryTasks(context.Background(), in)
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Content)
}