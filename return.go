package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type jsonData struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func successReturn(w http.ResponseWriter, msg string, data interface{}) {
	jsonReturn(w, 1, msg, data)
}

func errorRetrun(w http.ResponseWriter, e error, data interface{}) {
	jsonReturn(w, 0, e.Error(), data)
}

func jsonReturn(w http.ResponseWriter, status int, msg string, data interface{}) {

	//读取配置文件，设置header
	headers := conf.Get("headers")
	header, ok := headers.(map[string]interface{})
	if ok == true {
		for k, v := range header {
			val, ok := v.(string)
			if ok == false {
				log.Println("header配置读取出错")
			}
			w.Header().Set(k, val)
		}
	} else {
		log.Println("header配置读取出错")

	}

	w.Header().Set("Content-Type", "application/json")
	d := jsonData{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
	dataByte, err := json.Marshal(d)
	if err != nil {
		fmt.Fprint(w, err)
	}
	w.Write(dataByte)
	return
}
