package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"git.in.zhihu.com/zhihu/hello/model"
	"github.com/orca-zhang/borm"
)

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender int64  `json:"gender"`
}

func Format(user *model.ORMUser) *User{
	u := User{
		ID:     user.ID,
		Name:   user.Name,
		Gender: user.Gender,
	}
	return &u
}

func Reverse(user *User) *model.ORMUser{
	u := model.ORMUser{
		ID:     user.ID,
		Name:   user.Name,
		Gender: user.Gender,
	}
	return &u
}


type UserOperator struct {

}

func (uo *UserOperator) BatchGet(ctx *context.Context, ids []int) []*User{
	ouo := model.NewUserOperator(ctx)
	q := ouo.Query(ctx, borm.Where(borm.In("id", ids)))
	userModels := q.Interface().([]*model.ORMUser)
	results := make([]*User, 0)
	PrintJson("Batch get %v", ids)
	for _, um := range userModels{
		PrintJson("Batch get %v", um)
		u := Format(um)
		results = append(results, u)
	}
	PrintJson("Batch get %v", results)
	return results
}

func (uo *UserOperator) Get(ctx *context.Context, id int) *User{
	ouo := model.NewUserOperator(ctx)
	userModels := ouo.Query(ctx, borm.Where(borm.Eq("id",
		id))).Interface().([]*model.ORMUser)
	for _, um := range userModels{
		u := Format(um)
		return u
	}
	return nil
}

func (uo *UserOperator) Create(ctx *context.Context, user *User) bool{
	ouo := model.NewUserOperator(ctx)
	u := Reverse(user)
	return ouo.Create(ctx, u) == 1
}

func (uo *UserOperator) Update(ctx *context.Context, user *User) bool{
	ouo := model.NewUserOperator(ctx)
	u := Reverse(user)
	return ouo.Update(ctx, u) == 1
}

func PrintJson(format string, args ...interface{}){
	argJsons := make([]interface{}, 0)
	for i := 0; i<len(args); i++{
		s, e := json.Marshal(args[i])
		if e != nil{
			fmt.Printf("error format %v", e)
		}
		argJsons = append(argJsons, string(s))
	}
	fmt.Printf(format, argJsons...)
}