package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"reflect"
)

type Table interface {
	PrimaryKey() map[string]interface{}
	TableName() string
}

func BuildInsertSQL(t interface{}) (string, []interface{}, error){
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	tableName := t.(Table).TableName()

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
			return "", nil, fmt.Errorf("field %v not set json tag", field.Name)
		}
		columns = append(columns, name)
		vf := tv.Field(i)
		values = append(values, vf.Interface())
	}

	sqlBuilder := sq.Insert(tableName).Columns(columns...).Values(values...)
	return sqlBuilder.ToSql()
}

func BuildUpdateSQL(t interface{}) (string, []interface{}, error){
	values := make(map[string]interface{}, 0)
	primary := t.(Table).PrimaryKey()
	tableName := t.(Table).TableName()

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
			return "", nil, fmt.Errorf("field %v not set json tag", field.Name)
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
