package model

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
)

type ORMUser struct {
	ID     int `borm:"id"`
	Name   string `borm:"name"`
	Gender int64  `borm:"gender"`
}

func NewUserOperator(ctx *context.Context) Operator{
	operator := NewOperator(ctx, &ORMUser{}, "user")
	return operator
}