package test

import (
	"context"
)

type Test struct{}

func New() *Test {
	return &Test{}
}

func (t *Test) Name() string {
	return "test"
}

func (t *Test) Description() string {
	return "测试 Action，用于验证管道和执行流程"
}

func (t *Test) Help() string {
	return helpText
}

func (t *Test) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (
	map[string]interface{},
	error,
) {

	val := ""

	if len(args) > 0 {
		val = args[0]
	}

	return map[string]interface{}{
		"value":          val,
		"received_count": len(input),
		"args":           args,
		"status":         "ok",
	}, nil
}
