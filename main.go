package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	scujwc "github.com/mohuishou/scujwc-go"
)

func main() {
	port := flag.String("p", "6777", "请输入需要绑定的端口号")
	dist := flag.String("d", "dist/", "请输入前端目录路径")
	flag.Parse()

	log.Printf("即将监听端口: %s", *port)

	http.Handle("/", http.FileServer(http.Dir(*dist)))
	http.HandleFunc("/gpa", gpaAll)
	http.HandleFunc("/gpa/all", gpaAll)
	http.HandleFunc("/gpa/not-pass", gpaNotPass)

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", *port), nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

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
	password := url.QueryEscape(v.Get("password"))
	s, err = scujwc.NewJwc(uid, password)
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
