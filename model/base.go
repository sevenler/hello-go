package model

import (
	sqlCon "database/sql"
	sq "github.com/Masterminds/squirrel"

	"context"
	"fmt"
	"reflect"
)

type Operator interface {
	Query(ctx *context.Context, queryArgs ...interface{}) (interface{}, error)
	Create(ctx *context.Context, user interface{}) (int64, error)
	Update(ctx *context.Context, user interface{}) (int64, error)
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

func GetConnection(ctx *context.Context) (*sqlCon.DB, error){
	return sqlCon.Open("mysql","root:@/butterfly?parseTime=true")
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
	return Scan(ctx, uo.ModelStruct, *sql, args...)
}

func (uo OperatorImpl) Create(ctx *context.Context, user interface{}) (int64, error){
	sql, args, err := BuildInsertSQL(user)
	if err != nil{
		return 0, err
	}

	con, err := GetConnection(ctx)
	if err != nil{
		panic(err)
	}

	rows, err := con.Exec(sql, args...)
	count, err := rows.RowsAffected()
	if err != nil{
		panic(err)
	}
	return count, nil
}


func (uo OperatorImpl) Update(ctx *context.Context, user interface{}) (int64, error){
	sql, args, err := BuildUpdateSQL(user)
	if err != nil{
		return 0, err
	}

	con, err := GetConnection(ctx)
	if err != nil{
		panic(err)
	}

	rows, err := con.Exec(sql, args...)
	count, err := rows.RowsAffected()
	if err != nil{
		panic(err)
	}
	return count, nil
}
