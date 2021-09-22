package logic

import (
	"context"
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
	ouo := model.ORMUserOperator{}
	userModels := ouo.Query(ctx, borm.Where(borm.In("id", ids)))
	results := make([]*User, 0)
	for _, um := range userModels{
		u := Format(um)
		results = append(results, u)
	}
	return results
}

func (uo *UserOperator) Get(ctx *context.Context, id int) *User{
	ouo := model.ORMUserOperator{}
	userModels := ouo.Query(ctx, borm.Where(borm.Eq("id", id)))
	for _, um := range userModels{
		u := Format(um)
		return u
	}
	return nil
}

func (uo *UserOperator) Create(ctx *context.Context, user *User) bool{
	ouo := model.ORMUserOperator{}
	u := Reverse(user)
	return ouo.Create(ctx, u) == 1
}

func (uo *UserOperator) Update(ctx *context.Context, user *User) bool{
	ouo := model.ORMUserOperator{}
	u := Reverse(user)
	return ouo.Update(ctx, u) == 1
}