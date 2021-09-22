package model

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	b "github.com/orca-zhang/borm"
	"time"
)

type ORMUser struct {
	ID     int `borm:"id"`
	Name   string `borm:"name"`
	Gender int64  `borm:"gender"`
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


type ORMUserOperator struct {}

func (uo *ORMUserOperator) TableName() string{
	return "user"
}

func (uo *ORMUserOperator) Query(ctx *context.Context, filterArgs ...b.BormItem) []*ORMUser{
	t := GetORMTable(ctx, uo.TableName())
	ret := make([]*ORMUser, 0)
	t.Select(&ret, filterArgs...)
	return ret
}

func (uo *ORMUserOperator) Create(ctx *context.Context, user *ORMUser) int64{
	t := GetORMTable(ctx, uo.TableName())
	n, err := t.Insert(user)
	if err != nil{
		panic(err)
	}
	return int64(n)
}

func (uo *ORMUserOperator) Update(ctx *context.Context, user *ORMUser) int64{
	t := GetORMTable(ctx, uo.TableName())
	n, err := t.Update(user)
	if err != nil{
		panic(err)
	}
	return int64(n)
}