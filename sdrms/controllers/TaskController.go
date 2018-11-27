package controllers

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
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
	etcdHosts := beego.AppConfig.Strings("etcd::hosts")
	prefix := beego.AppConfig.String("etcd::prefix")
	serviceName := beego.AppConfig.String("etcd::task_server")
	client.InitConn(etcdHosts, prefix, serviceName)
}

type result struct {
	Err string
	Content interface{}
}

func (tc *TaskController) GetTasks() {
	b, _ := ioutil.ReadAll(tc.Ctx.Request.Body)
	param := make(map[string]*models.Task_QueryParam)
	err := json.Unmarshal(b, &param)
	r := new(result)
	if err != nil {
		r.Err = err.Error()
		r.Content = nil
		tc.Data["json"] = r
		tc.ServeJSON()
	}
	ts, err := client.QueryTasks(context.Background(), param["param"])
	if err != nil {
		r.Err = err.Error()
		r.Content = nil
		tc.Data["json"] = r
		tc.ServeJSON()
	} else {
		r.Content = ts
		r.Err = ""
		tc.Data["json"] = r
		tc.ServeJSON()
	}
}
