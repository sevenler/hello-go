package model

import (
	_ "github.com/go-sql-driver/mysql"
)

type ORMUser struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender int64  `json:"gender"`
}

func (om *ORMUser)PrimaryKey() map[string]interface{}{
	return map[string]interface{}{
		"id": om.ID,
	}
}

func (om *ORMUser)TableName() string{
	return "user"
}