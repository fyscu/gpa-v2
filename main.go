package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/goless/config"
	"github.com/mohuishou/scujwc-go"

	"encoding/json"

	"fmt"

	"errors"
	"net/url"
	"strconv"
)

func login(w http.ResponseWriter, r *http.Request) (scujwc.Jwc, error) {
	var s scujwc.Jwc
	if r.Method != "POST" {
		return s, errors.New("请求方式错误")
	}
	result, _ := ioutil.ReadAll(r.Body)
	v, err := url.ParseQuery(string(result))
	if err != nil {
		return s, err
	}
	uid, err := strconv.Atoi(v.Get("uid"))
	if err != nil {
		return s, err
	}
	password := v.Get("password")
	err = s.Init(uid, password)
	if err != nil {
		return s, err
	}
	return s, nil
}

func gpa(w http.ResponseWriter, r *http.Request) {
	s, err := login(w, r)
	if err != nil {
		errorRetrun(w, err, nil)
		return
	}
	grade, err := s.GPA()
	if err != nil {
		errorRetrun(w, err, nil)
		return
	}
	successReturn(w, "本学期成绩获取成功", grade)
}

func gpaAll(w http.ResponseWriter, r *http.Request) {
	s, err := login(w, r)
	if err != nil {
		errorRetrun(w, err, nil)
		return
	}
	grade, err := s.GPAAll()
	if err != nil {
		errorRetrun(w, err, nil)
		return
	}
	successReturn(w, "所有成绩获取成功", grade)
}

func gpaNotPass(w http.ResponseWriter, r *http.Request) {
	s, err := login(w, r)
	if err != nil {
		errorRetrun(w, err, nil)
		return
	}
	grade, err := s.GPANotPass()
	if err != nil {
		errorRetrun(w, err, nil)
		return
	}
	successReturn(w, "不及格成绩获取成功", grade)
}

func main() {
	http.HandleFunc("/gpa", gpa)
	http.HandleFunc("/gpa/all", gpaAll)
	http.HandleFunc("/gpa/not-pass", gpaNotPass)

	//读取配置文件，设置端口
	config := config.New("config.json")
	ports := config.Get("port")
	port, ok := ports.(string)
	if ok == false {
		log.Fatal("端口配置读取出错，请使用\"port\":\"123\" 形式表示端口号 ")
	}

	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

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
	config := config.New("config.json")
	headers := config.Get("headers")
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
	if status == 1 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	w.Write(dataByte)
	return
}
