package model

import (
	"context"
	"database/sql"
	"fmt"
	b "github.com/orca-zhang/borm"
	"reflect"
	"runtime/debug"
	"time"
)

type Operator interface {
	Query(ctx *context.Context, filterArgs ...b.BormItem) reflect.Value
	Create(ctx *context.Context, user interface{}) int64
	Update(ctx *context.Context, user interface{}) int64
}



var ORMConnectionString = "root:@/butterfly"

func GetORMConnection(ctx *context.Context, connectString *string) *sql.DB{
	db, err := sql.Open("mysql", *connectString)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func GetORMTable(ctx *context.Context, tableName string) *b.BormTable{
	conn := GetORMConnection(ctx, &ORMConnectionString)
	t := b.Table(conn, tableName)
	return t
}



type OperatorImpl struct {
	TableName string
	ModelStruct interface{}
}

func NewOperator(ctx *context.Context, model interface{}, tableName string) Operator{
	return OperatorImpl{
		TableName:   tableName,
		ModelStruct: model,
	}
}

func (uo OperatorImpl) Query(ctx *context.Context, filterArgs ...b.BormItem) reflect.Value{
	t := GetORMTable(ctx, uo.TableName)
	where := b.Where(b.Eq("id", 4))

	array := make([]*ORMUser, 0)
	_, err := t.Select(&array, where)
	fmt.Println(filterArgs)
	if err != nil{
		debug.PrintStack()
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(fmt.Sprintf("==1== %v %v %v", reflect.TypeOf(array), reflect.ValueOf(array), array))

	class := reflect.TypeOf(uo.ModelStruct)
	reflectArray := reflect.MakeSlice(reflect.SliceOf(class), 0, 2).Interface()
	objArray := reflectArray.([]*ORMUser)
	_, err = t.Select(&objArray, where)
	fmt.Println(fmt.Sprintf("==2== %v %v %v", reflect.TypeOf(objArray), reflect.ValueOf(objArray), objArray))


	for _, row := range objArray{
		fmt.Println(fmt.Sprintf("000 Name: %s", row.Name))
	}



	fmt.Println("==========----=========")
	fmt.Println(where)
	fmt.Println(where.Type())
	fmt.Println("==========----=========")
	if err != nil{
		fmt.Println("===================")
		debug.PrintStack()
		fmt.Println("===================")
		panic(err)
	}


	//class := reflect.TypeOf(&ORMUser{})
	//reflectArray := reflect.MakeSlice(reflect.SliceOf(class), 0, 0)


	//reflectArray := reflect.ValueOf(array)


	return reflect.ValueOf(reflectArray)
}

func (uo OperatorImpl) Create(ctx *context.Context, user interface{}) int64{
	t := GetORMTable(ctx, uo.TableName)
	n, err := t.Insert(user)
	if err != nil{
		panic(err)
	}
	return int64(n)
}

func (uo OperatorImpl) Update(ctx *context.Context, user interface{}) int64{
	t := GetORMTable(ctx, uo.TableName)
	n, err := t.Update(user)
	if err != nil{
		panic(err)
	}
	return int64(n)
}

func nameToStruct(ctx *context.Context, tableName string) interface{}{
	// ??
	structs := map[string]interface{}{
		"user": ORMUser{},
	}
	class := structs[tableName]
	return class
}
