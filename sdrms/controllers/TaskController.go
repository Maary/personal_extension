package controllers

import (
	"encoding/json"
	"log"
	"personal_extension/sdrms/models"
)

type TaskController struct {
	BaseController
}

type BackendMessage struct {
	Code int
	Message string
}

func (tc *TaskController) Prepare() {
	tc.BaseController.Prepare()
}

func (tc *TaskController) ReceiveData() {
	result := tc.GetString("data")
	bm := new(BackendMessage)
	r := new(models.Result)
	if err := json.Unmarshal([]byte(result), r); err != nil {
		log.Println(err)
		bm.Message = err.Error()
		bm.Code = 400
		tc.Data["json"] = bm
		tc.ServeJSON()
	} else {
		if err := models.Insert(r); err != nil {
			bm.Code = 400
			bm.Message = err.Error()
			tc.Data["json"] = bm
			tc.ServeJSON()
		} else {
			bm.Code = 200
			tc.Data["json"] = bm
			tc.ServeJSON()
		}
	}
}