package logic

import (
	"context"
	"git.in.zhihu.com/zhihu/hello/model"
	sq "github.com/Masterminds/squirrel"
)

type UserOperatorInterface interface {
	BatchGet(ctx *context.Context, ids []int) []*User
	Get(ctx *context.Context, id int) *User
	Create(ctx *context.Context, user *User) bool
	Update(ctx *context.Context, user *User) bool
}

type UserOperator struct {
	operator *model.Operator
}

func NewUserOperator(operator *model.Operator) *UserOperator {
	return &UserOperator{
		operator: operator,
	}
}

func (uo *UserOperator) BatchGet(ctx *context.Context, ids []int) []*User{
	ouo := *uo.operator
	args := []interface{}{sq.Eq{"id": ids}}
	var v model.Table = model.ORMUser{}
	q, err := ouo.Query(ctx, &v, args)
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
	ouo := *uo.operator
	args := []interface{}{sq.Eq{"id": id}}
	var v model.Table = model.ORMUser{}
	q, err := ouo.Query(ctx, &v, args)
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
	ouo := *uo.operator
	u := Reverse(user)
	var v model.Table = *u
	c, err := ouo.Create(ctx, &v)
	if err != nil{
		panic(err)
	}
	return c == 1
}

func (uo *UserOperator) Update(ctx *context.Context, user *User) bool{
	ouo := *uo.operator
	u := Reverse(user)
	var v model.Table = *u
	c, err := ouo.Update(ctx, &v)
	if err != nil{
		panic(err)
	}
	return c == 1
}