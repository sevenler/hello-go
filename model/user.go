package model

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
)

type ORMUser struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender int64  `json:"gender"`
}

func NewUserOperator(ctx *context.Context) Operator{
	operator := NewOperator(ctx, &ORMUser{}, "user")
	return operator
}