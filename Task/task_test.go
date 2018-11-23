package Task

import "testing"

func TestDataContainer_Send(t *testing.T) {
	a := make(map[string]interface{})
	a["1"] = 1
	a["2"] = 2
	send(a, "http://127.0.0.1:8080/task/input")
}


