package logic

import "git.in.zhihu.com/zhihu/hello/model"

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
