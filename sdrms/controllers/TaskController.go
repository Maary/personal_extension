package controllers

import (
	"context"
	"crawler.center/lib/system"
	"encoding/json"
	"personal_extension/Task/client"
	"personal_extension/Task/models"
)

type TaskController struct {
	BaseController
}

type BackendMessage struct {
	Code int
	Message string
}

//TODO: add auth
func (tc *TaskController) Prepare() {
	tc.BaseController.Prepare()
	etcdHosts := system.AppConfig.Strings("etcd::hosts")
	prefix := system.AppConfig.String("etcd::prefix")
	serviceName := system.AppConfig.String("etcd::task_server")
	client.InitConn(etcdHosts, prefix, serviceName)
}

type result struct {
	Err string
	Content interface{}
}

func (tc *TaskController) GetTasks() {
	paramStr := tc.GetString("param")
	param := new(models.Task_QueryParam)
	err := json.Unmarshal([]byte(paramStr), param)
	r := new(result)
	if err != nil {
		r.Err = err.Error()
		r.Content = nil
		tc.Data["json"] = r
		tc.ServeJSON()
	}
	ts, err := client.QueryTasks(context.Background(), param)
	if err != nil {
		r.Err = err.Error()
		r.Content = nil
		tc.Data["json"] = r
		tc.ServeJSON()
	} else {
		r.Content = ts
		r.Err = err.Error()
		tc.Data["json"] = r
		tc.ServeJSON()
	}
}
