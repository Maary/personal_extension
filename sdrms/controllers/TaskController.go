package controllers

type TaskController struct {
	BaseController
}

type BackendMessage struct {
	Code int
	Message string
}

//TODO
func (tc *TaskController) Prepare() {
	tc.BaseController.Prepare()
}


