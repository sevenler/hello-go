package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	b "github.com/orca-zhang/borm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type ORMUser struct {
	ID     int `borm:"id"`
	Name   string `borm:"name"`
	Gender int64  `borm:"gender"`
}

func (*ORMUser) TableName() string{
	return "user"
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

func ORMWrite(ctx *context.Context, user ORMUser) int64{
	t := GetORMTable(ctx, user.TableName())
	n, err := t.Insert(&user)
	if err != nil{
		panic(err)
	}
	return int64(n)
}

func ORMRead(ctx *context.Context) []*ORMUser{
	user := &ORMUser{}
	t := GetORMTable(ctx, user.TableName())
	ret := make([]*ORMUser, 0)
	t.Select(&ret, b.Where(b.Eq("gender", 10)))
	return ret
}

func main(){
	// Demo about sample write and read
	ctx := context.Background()

	// Write Data to db
	ORMWrite(&ctx, ORMUser{
		Name:   "Hello",
		Gender: 10,
	})
	ORMWrite(&ctx, ORMUser{
		Name:   "Hello",
		Gender: 10,
	})

	// Read Data from db
	users := ORMRead(&ctx)
	json, err := json.Marshal(users)
	if err != nil{
		panic(err)
	}
	fmt.Printf("Users: %v", string(json))
}