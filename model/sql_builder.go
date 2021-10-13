package model

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"reflect"
)

func BuildSelectSQL(tb *Table, queryArgs ...interface{}) (*string, []interface{}, error){
	tv := *tb
	t := reflect.TypeOf(tv)
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

	tableName := tv.TableName()
	sqlBuilder := sq.Select(names...).From(tableName)
	if len(queryArgs) > 0{
		sqlBuilder = sqlBuilder.Where(queryArgs[0], queryArgs[1:]...)
	}
	sql, args, err := sqlBuilder.ToSql()
	if err != nil{
		return nil, nil, err
	}
	return &sql, args, nil
}

func BuildInsertSQL(tb *Table) (string, []interface{}, error){
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	t := *tb
	tableName := t.TableName()

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

func BuildUpdateSQL(tb *Table) (string, []interface{}, error){
	values := make(map[string]interface{}, 0)
	t := *tb
	primary := t.PrimaryKey()
	tableName := t.TableName()

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
