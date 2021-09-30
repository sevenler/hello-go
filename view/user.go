package view

import (
	"encoding/json"
	"git.in.zhihu.com/zhihu/hello/logic"
	"math/rand"
	"net/http"
)

type ResponseJson struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ResponseObject(code int, message string, data interface{})[]byte{
	rj := ResponseJson{
		Code:    code,
		Message: message,
		Data:    data,
	}
	js, err := json.Marshal(rj)
	if err != nil{
		panic(err)
	}

	return js
}

func BatchGetUser(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uo := logic.UserOperator{}
	users := uo.BatchGet(&ctx, []int{1, 2, 3, 4})
	w.Write(ResponseObject(0, "", users))
}

func GetUser(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uo := &logic.UserOperator{}
	user := uo.Get(&ctx, 3)
	w.Write(ResponseObject(10, "---", user))
}

func CreateUser(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uo := logic.UserOperator{}
	user := logic.User{
		Name:   "Daniel's HTTP",
		Gender: 100,
	}
	count := uo.Create(&ctx, &user)
	w.Write(ResponseObject(0, "", count))
}

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	uo := logic.UserOperator{}
	user := logic.User{
		ID: 1,
		Name:   "Daniel's HTTP Updated",
		Gender: rand.Int63n(100000),
	}
	count := uo.Update(&ctx, &user)
	w.Write(ResponseObject(0, "", count))
}