package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"personal_extension/Task/models"
	"personal_extension/Task/task_service"
	"personal_extension/lib/service"
)

var (
	client task_rpc_config.TaskClient
)

func InitConn(hosts []string, prefix string, serviceName string) {
	if client == nil {
		if len(hosts) > 0 && prefix != "" {
			err := service.InitEtcd(hosts, prefix)
			if err != nil {
				log.Println(err)
			}
		}
		conn, err := service.Dial(serviceName)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(hosts, prefix, serviceName)
		client = task_rpc_config.NewTaskClient(conn)
	}
}

func QueryTasks(ctx context.Context, params *models.Task_QueryParam) (ts []*models.Task, err error) {
	spar := new(task_rpc_config.Params)
	jsonB, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	spar.ParamsStr = string(jsonB)
	rsp, err := client.QueryTasks(context.Background(), spar)
	if err != nil {
		return nil, err
	}
	ts = make([]*models.Task, 0)
	if err := json.Unmarshal([]byte(rsp.Content), &ts); err != nil {
		return nil, err
	}
	return ts, nil
}

func InsertTasks(ctx context.Context, ts []*models.Task) (status string, err error) {
	st := new(task_rpc_config.Tasks)
	if len(ts) > 0 {
		if jsonB, err := json.Marshal(ts); err != nil {
			return "failed", err
		} else {
			st.Content = string(jsonB)
		}
	}
	rsp, err := client.InsertTasks(context.Background(), st)
	if err != nil {
		return "fields", err
	}
	return "success", errors.New(rsp.ErrMessage) //TODO: the error should be nil
}
