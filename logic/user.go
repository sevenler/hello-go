package logic

import (
	"context"
	"git.in.zhihu.com/zhihu/hello/model"
	sq "github.com/Masterminds/squirrel"
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


type UserOperator struct {}

func (uo *UserOperator) BatchGet(ctx *context.Context, ids []int) []*User{
	ouo := model.NewOperator()
	q, err := ouo.Query(ctx, &model.ORMUser{}, sq.Eq{"id": ids})
	if err != nil{
		panic(err)
	}
	userModels := q.([]interface{})
	results := make([]*User, 0)

	for _, um := range userModels{
		p := um.(model.ORMUser)
		u := Format(&p)
		results = append(results, u)
	}
	return results
}

func (uo *UserOperator) Get(ctx *context.Context, id int) *User{
	ouo := model.NewOperator()
	q, err := ouo.Query(ctx, &model.ORMUser{}, sq.Eq{"id": id})
	if err != nil{
		panic(err)
	}
	for _, um := range q.([]interface{}){
		p := um.(model.ORMUser)
		u := Format(&p)
		return u
	}
	return nil
}

func (uo *UserOperator) Create(ctx *context.Context, user *User) bool{
	ouo := model.NewOperator()
	u := Reverse(user)
	c, err := ouo.Create(ctx, u)
	if err != nil{
		panic(err)
	}
	return c == 1
}

func (uo *UserOperator) Update(ctx *context.Context, user *User) bool{
	ouo := model.NewOperator()
	u := Reverse(user)
	c, err := ouo.Update(ctx, u)
	if err != nil{
		panic(err)
	}
	return c == 1
}