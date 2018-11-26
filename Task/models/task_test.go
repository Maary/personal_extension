package models

import (
	"testing"
	"time"
)

func TestDataContainer_httpSend(t *testing.T) {
	a := make(map[string]interface{})
	a["name"] = "wsz"
	a["vale"] = "100"
	send(a, "http://127.0.0.1:8080/task/input")
}

func TestDataContainer_Send(t *testing.T) {
	dua := time.Duration(1 * time.Minute)
	Container := New(10, dua)
	Container.SetUrl("http://127.0.0.1:8080/task/input")
	for i := 0; i < 10; i++ {
		Container.Push(map[string]interface{}{"1": 1, "2": 2})
	}
	Container.Send(SEND_PER)
}
