package model

import (
	"context"
	sqlCon "database/sql"
)

type Option struct {
	// 是否使用事务，如果使用事务，则会到 context 中获取事务
	ShouldUseTransaction bool
}
type SetOption func(option *Option)

func HoldOption(options ...SetOption) *Option {
	op := &Option{
		ShouldUseTransaction: true,
	}
	for _, s := range options {
		s(op)
	}
	return op
}

type Operator interface {
	Query(ctx *context.Context, st interface{}, queryArgs []interface{}, options ...SetOption) (interface{}, error)
	Create(ctx *context.Context, st interface{}, options ...SetOption) (int64, error)
	Update(ctx *context.Context, st interface{}, options ...SetOption) (int64, error)
}


type OperatorImpl struct {}

func NewOperator() Operator{
	return OperatorImpl{}
}

func GetConnection(ctx *context.Context) (*sqlCon.DB, error){
	return sqlCon.Open("mysql","root:@/butterfly?parseTime=true")
}

func (uo OperatorImpl) Query(ctx *context.Context, st interface{}, queryArgs []interface{}, options ...SetOption) (interface{}, error){
	_ := HoldOption(options...)
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

func (uo OperatorImpl) Create(ctx *context.Context, st interface{}, options ...SetOption) (int64, error){
	_ := HoldOption(options...)
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

func (uo OperatorImpl) Update(ctx *context.Context, st interface{}, options ...SetOption) (int64, error){
	_ := HoldOption(options...)
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
