package mock_logic

import (
	"context"
	"git.in.zhihu.com/zhihu/hello/logic"
	"git.in.zhihu.com/zhihu/hello/model"
	mock_model "git.in.zhihu.com/zhihu/hello/model/mock"
	"testing"
)

// go test -run TestBatchGet -v
func TestBatchGet(t *testing.T) {
	ctx := context.Background()
	u := model.NewOperator()
	uo := logic.NewUserOperator(&u)
	users := uo.BatchGet(&ctx, []int{1, 2, 3, 4})
	// db 查询返回0
	if len(users) == 0{
		t.Errorf("Batch Get want be 4, but got %d", len(users))
	}

	ctrl, mockUserOperator := mock_model.MockUserOperator(t, &ctx)
	defer ctrl.Finish()
	op := logic.NewUserOperator(&mockUserOperator)
	users = op.BatchGet(&ctx, []int{1, 2, 3, 4})
	// mock 查询返回1
	if len(users) != 1{
		t.Errorf("Batch Get want be 4, but got %d", len(users))
	}
}