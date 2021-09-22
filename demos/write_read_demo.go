package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Gender int64  `json:"gender"`
}

var ConnectionString = "root:@/butterfly"

func GetConnection(ctx *context.Context, connectString *string) *sql.DB{
	db, err := sql.Open("mysql", *connectString)
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

func Write(ctx *context.Context, user *User) int64{
	db := GetConnection(ctx, &ConnectionString)
	sql := fmt.Sprintf("INSERT INTO user (name, gender) VALUES ('%s', %d)", user.Name, user.Gender)
	ret, err := db.Exec(sql)
	if err != nil{
		panic(err)
	}
	updateCount, err := ret.RowsAffected()
	if err != nil{
		panic(err)
	}
	return updateCount
}

func Read(ctx *context.Context) []*User{
	db := GetConnection(ctx, &ConnectionString)
	sql := fmt.Sprintf("SELECT id, name, gender FROM user WHERE gender in (0, 1)")
	rows, err := db.Query(sql)
	if err != nil{
		panic(err)
	}

	users := make([]*User, 0)
	for {
		hasNext := rows.Next()
		if hasNext == false{
			break
		}
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Gender)
		if err != nil{
			panic(err)
		}
		users = append(users, user)
	}
	return users
}


func main(){
	// Demo about sample write and read
	ctx := context.Background()

	// Write Data to db
	//Write(&ctx, &User{
	//	Name:   "Daniel",
	//	Gender: 0,
	//})
	//Write(&ctx, &User{
	//	Name:   "Miss Yang",
	//	Gender: 2,
	//})

	// Read Data from db
	ret := Read(&ctx)
	json, err := json.Marshal(ret)
	if err != nil{
		panic(err)
	}
	fmt.Printf("Users: %v", string(json))
}