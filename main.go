package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	scujwc "github.com/mohuishou/scujwc-go"
)

func main() {
	// http.Handle("/", http.FileServer(http.Dir("/tmp/static/")))
	http.HandleFunc("/gpa", gpa)
	http.HandleFunc("/gpa/all", gpaAll)
	http.HandleFunc("/gpa/not-pass", gpaNotPass)

	port := "6777"

	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
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
	password := v.Get("password")
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
