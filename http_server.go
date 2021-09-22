package main

import (
	"fmt"
	"git.in.zhihu.com/zhihu/hello/view"
	"io/ioutil"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	//url, _ := req.URL.Query()["url"]
	url := []string{
		"http://www.baidu.com",
	}

	resp, err := http.Get(url[0])
	defer resp.Body.Close()

	if err != nil{
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/user/batch/get", view.BatchGetUser)
	http.HandleFunc("/user/get", view.GetUser)
	http.HandleFunc("/user/create", view.CreateUser)
	http.ListenAndServe(":8090", nil)
}
