package model

import (
	"context"
	sqlCon "database/sql"
)

type Operator interface {
	Query(ctx *context.Context, st interface{}, queryArgs ...interface{}) (interface{}, error)
	Create(ctx *context.Context, st interface{}) (int64, error)
	Update(ctx *context.Context, st interface{}) (int64, error)
}


type OperatorImpl struct {}

func NewOperator() Operator{
	return OperatorImpl{}
}

func GetConnection(ctx *context.Context) (*sqlCon.DB, error){
	return sqlCon.Open("mysql","root:@/butterfly?parseTime=true")
}

func (uo OperatorImpl) Query(ctx *context.Context, st interface{}, queryArgs ...interface{}) (interface{}, error){
	sql, args, err := BuildSelectSQL(st, queryArgs...)
	if err != nil{
		return nil, err
	}

	con, err := GetConnection(ctx)
	defer con.Close()
	if err != nil {
		return nil, err
	}

	rows, err := con.Query(*sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	return Scan(st, rows)
}

func (uo OperatorImpl) Create(ctx *context.Context, st interface{}) (int64, error){
	sql, args, err := BuildInsertSQL(st)
	if err != nil{
		return 0, err
	}

	con, err := GetConnection(ctx)
	defer con.Close()
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

func (uo OperatorImpl) Update(ctx *context.Context, st interface{}) (int64, error){
	sql, args, err := BuildUpdateSQL(st)
	if err != nil{
		return 0, err
	}

	con, err := GetConnection(ctx)
	defer con.Close()
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
