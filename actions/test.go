package actions

import (
	"context"
)

type Test struct{}

func NewTest() *Test {
	return &Test{}
}

func (t *Test) Name() string {
	return "test"
}

func (t *Test) Execute(
	ctx context.Context,
	args []string, // 👈 新增：对齐 Action 接口规范[cite: 4]
	input map[string]interface{},
) (
	map[string]interface{},
	error,
) {

	return map[string]interface{}{
		"received": input,
		"args":     args, // 👈 顺便把参数也吐出来，方便你调试
		"status":   "ok",
	}, nil
}
