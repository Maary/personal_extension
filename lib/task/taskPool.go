package task

type taskPool interface {
	GenerateTask(interface{}) bool
	PullTask() interface{}
	DeleteTask(int) bool
}

