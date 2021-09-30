package main

import (
	"fmt"
	"git.in.zhihu.com/zhihu/hello/utils"
	sq "github.com/Masterminds/squirrel"
	"reflect"
	"time"
)

func BuildInsertSQL(t interface{}, tableName string) (string, interface{}, error){
	columns := make([]string, 0)
	values := make([]interface{}, 0)

	tp := reflect.TypeOf(t)
	if tp.Kind() == reflect.Ptr{
		tp = tp.Elem()
	}
	tv := reflect.ValueOf(t)
	if tv.Kind() == reflect.Ptr{
		tv = tv.Elem()
	}

	for i := 0; i < tp.NumField();i++{
		field := tp.Field(i)
		name := field.Tag.Get("json")
		if name == ""{
			return "", "", fmt.Errorf("field %v not set json tag", field.Name)
		}
		columns = append(columns, name)
		vf := tv.Field(i)
		values = append(values, vf.Interface())
	}

	sqlBuilder := sq.Insert(tableName).Columns(columns...).Values(values...)
	return sqlBuilder.ToSql()
}

type Primary interface {
	PrimaryKey() map[string]interface{}
}

func BuildUpdateSQL(t interface{}, tableName string) (string, interface{}, error){
	values := make(map[string]interface{}, 0)
	primary := t.(Primary).PrimaryKey()

	tp := reflect.TypeOf(t)
	if tp.Kind() == reflect.Ptr{
		tp = tp.Elem()
	}
	tv := reflect.ValueOf(t)
	if tv.Kind() == reflect.Ptr{
		tv = tv.Elem()
	}

	for i := 0; i < tp.NumField();i++{
		field := tp.Field(i)
		name := field.Tag.Get("json")
		if name == ""{
			return "", "", fmt.Errorf("field %v not set json tag", field.Name)
		}
		vf := tv.Field(i)
		values[name] = vf.Interface()

		primaryKey := field.Tag.Get("primary")
		if primaryKey != ""{
			primary[name] = vf.Interface()
		}
	}



	sqlBuilder := sq.Update(tableName)
	for column, value := range values{
		sqlBuilder = sqlBuilder.Set(column, value)
	}
	noPrimary := true
	for column, key := range primary{
		sqlBuilder = sqlBuilder.Where(sq.Eq{column: key})
		noPrimary = false
	}
	if noPrimary{
		return "", nil, fmt.Errorf("not set primary key")
	}
	return sqlBuilder.ToSql()
}

type DmUser struct {
	ID      *int       `json:"id",primary:"1"`
	Name    *string    `json:"name"`
	Gender  *int64     `json:"gender"`
	Created *time.Time `json:"created"`
	Time    *time.Time `json:"time"`
}

func (du DmUser)PrimaryKey() map[string]interface{}{
	return map[string]interface{}{
		"id": du.ID,
	}
}

func main() {
	u := DmUser{
		ID:      utils.IntPtr(1),
		Name:    utils.StringPtr("hello"),
		Gender:  utils.Int64Ptr(10),
	}

	sql, args, err := BuildInsertSQL(u, "user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("Insert SQL:%v  Args:%v \n", sql, args)

	sql, args, err = BuildInsertSQL(&u, "user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("Insert SQL Ptr:%v  Args:%v \n", sql, args)


	sql, args, err = BuildUpdateSQL(u, "user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("Update SQL:%v  Args:%v \n", sql, args)

	sql, args, err = BuildUpdateSQL(&u, "user")
	if err != nil{
		panic(err)
	}
	utils.PrintJson("Update SQL Ptr:%v  Args:%v \n", sql, args)
}