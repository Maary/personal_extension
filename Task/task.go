package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type DataContainer struct {
	Data chan interface{}
	Length uint8
	Slice time.Duration
	SendMethod sendMethod
	Url string //TODO
}

type sendMethod uint8
const (
	SEND_PATCH = sendMethod(0)
	SEND_PER = sendMethod(1)
)

func (dc *DataContainer) New(length uint8, slice time.Duration) *DataContainer {
	if length <= 0 {
		panic("Data container store data and send to backend, the length must greater than 0")
	}

	if slice >= time.Duration(2 * time.Minute) {
		panic("Data container store data and send to backend, the time slice must not greater than 2m")
	}

	return &DataContainer{
		Data: make(chan interface{}, length),
		Length: length,
		Slice: slice,
	}
}

func(dc *DataContainer) Push(data interface{}) {
	dc.Data <- data
}

//TODO
func (dc *DataContainer) Send(method sendMethod) {
	dataSet := make([]interface{}, 0)
	if uint8(len(dc.Data)) < dc.Length {
		return
	}
	for {
		for i := uint8(0); i < dc.Length; i++ {
			switch method {
			case SEND_PER:
				per := <-dc.Data
				send(per, dc.Url)
			case SEND_PATCH:
				per := <-dc.Data
				dataSet = append(dataSet, per)
			}
			if len(dataSet) > 0 {
				send(dataSet, dc.Url)
			}
		}
	}
}

func send(data interface{}, targetUrl string) bool {
	param := new(url.Values)
	dataByte, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return false
	}
	param.Set("data", string(dataByte))
	rsp, err := http.PostForm(targetUrl, *param)
	if err != nil {
		log.Println(err)
		return false
	}
	if rsp.StatusCode == 200 {
		return true
	}
	log.Println(ioutil.ReadAll(rsp.Body))
	return false
}
