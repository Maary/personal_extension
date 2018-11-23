package controllers

import "fmt"

type TaskController struct {
	BaseController
}

type BackendMessage struct {
	Code uint8
	Message string
}

func (tc *TaskController) Prepare() {
	tc.BaseController.Prepare()
}

func (tc *TaskController) ReceiveData() {
	result := tc.GetString("data")
	fmt.Println(result) //TODO
	bm := new(BackendMessage)
	bm.Code = 200
	tc.Data["json"] = bm
	tc.ServeJSON()
}