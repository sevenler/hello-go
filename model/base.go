package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"

	"context"
	"database/sql"
	b "github.com/orca-zhang/borm"
	"reflect"
	"time"
)

type Operator interface {
	Query(ctx *context.Context, queryArgs ...interface{}) (interface{}, error)
	Create(ctx *context.Context, user interface{}) (int64, error)
	Update(ctx *context.Context, user interface{}) (int64, error)
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

func (uo OperatorImpl) buildSQL(ctx *context.Context, queryArgs ...interface{}) (*string, []interface{}, error){
	t := reflect.TypeOf(uo.ModelStruct)
	if t.Kind() == reflect.Ptr{
		t = t.Elem()
	}
	names := make([]string, 0)
	for i := 0; i< t.NumField();i++{
		field := t.Field(i)
		v := field.Tag.Get("json")
		if v == ""{
			return nil, nil, fmt.Errorf("field %v not set json tag", field.Name)
		}
		names = append(names, v)
	}

	sqlBuilder := sq.Select(names...).From(uo.TableName)
	if len(queryArgs) > 0{
		sqlBuilder = sqlBuilder.Where(queryArgs[0], queryArgs[1:]...)
	}
	sql, args, err := sqlBuilder.ToSql()
	if err != nil{
		return nil, nil, err
	}
	return &sql, args, nil
}

func (uo OperatorImpl) Query(ctx *context.Context, queryArgs ...interface{}) (interface{}, error){
	sql, args, err := uo.buildSQL(ctx, queryArgs...)
	if err != nil{
		return nil, err
	}
	return Scan(uo.ModelStruct, *sql, args...)
}

func (uo OperatorImpl) Create(ctx *context.Context, user interface{}) (int64, error){
	t := GetORMTable(ctx, uo.TableName)
	n, err := t.Insert(user)
	if err != nil{
		panic(err)
	}
	return int64(n), nil
}

func (uo OperatorImpl) Update(ctx *context.Context, user interface{}) (int64, error){
	t := GetORMTable(ctx, uo.TableName)
	n, err := t.Update(user)
	if err != nil{
		panic(err)
	}
	return int64(n), nil
}
