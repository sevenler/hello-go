package mock_model

import (
	"context"
	"git.in.zhihu.com/zhihu/hello/model"
	"github.com/golang/mock/gomock"
	"testing"
)

func MockUserOperator(t *testing.T, ctx *context.Context) (*gomock.Controller, model.Operator) {
	ctrl := gomock.NewController(t)
	m := NewMockOperator(ctrl)
	u := model.ORMUser{
		ID:     1,
		Name:   "Hello",
		Gender: 0,
	}
	ret := []interface{}{u}
	m.EXPECT().
		Query(gomock.Eq(ctx), gomock.Any(), gomock.Any()).
		Return(ret, nil)
	return ctrl, m
}
