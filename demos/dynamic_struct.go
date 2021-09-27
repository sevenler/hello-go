package main

import (
	"git.in.zhihu.com/zhihu/hello/utils"
	"reflect"
)

func SetStructValue(){
	type stu struct { N int `json:"n"`}

	t := reflect.TypeOf(stu{})
	ts := reflect.New(t).Elem()
	// New 或者 Elem 返回的对象 Value 可能是任意的数据类型，起本质是一个指针
	// 通过 Kind() 判断其类型
	if ts.Kind() == reflect.Struct{
		field := ts.FieldByName("N")
		// FieldByName 没有找到会是一个空的 field
		// field 必须是 addressable
		if field.IsValid() && field.CanSet(){
			field.SetInt(30)
		}
	}


	utils.PrintJson("%v \n", reflect.ValueOf(stu{100}))
}

func main(){
	SetStructValue()
}