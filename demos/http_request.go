package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

var URL = "https://www.zhihu.com/api/v4/roundtables?limit=10&offset=10"

type ResponseJson struct {
	Data []struct {
		Name   string `json:"name"`
		Banner string `json:"banner"`
		Logo   string `json:"logo"`
	} `json:"data"`
}

func MakeHttpReqeust(ctx context.Context, url string) (int, *ResponseJson){
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil{
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}

	data := ResponseJson{}
	err = json.Unmarshal(body, &data)
	if err != nil{
		panic(err)
	}
	return resp.StatusCode, &data
}

func MakeOne(ctx context.Context, index int, wg *sync.WaitGroup){
	wg.Add(1)
	defer wg.Done()

	code, resp := MakeHttpReqeust(ctx, URL)
	json, err := json.Marshal(resp)
	if err != nil{
		panic(err)
	}
	fmt.Sprintf("%s", string(json))
	fmt.Printf("Response for index %d code %d \n", index, code)
}

func main(){
	var wg sync.WaitGroup

	// https://tour.golang.org/concurrency/1
	ctx := context.Background()
	for i := 0; i <= 10; i++{
		go MakeOne(ctx, i, &wg)
	}

	MakeOne(ctx, 11, &wg)
	wg.Wait()
}

// 知识点
//使用 go 关键字可以非常方便地发起一个 goroutine
//在脚本中，如果不实用 WaitGroup 这样的关键字，
//脚本并不会等待所有的 goruntine 全部执行完成
//那么，如果是一个 runserver 式的常驻进程呢？
// http 进程如何管理多个 goruntine 呢？