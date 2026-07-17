package pwd

import (
	"context"
	"fmt"
)

type Action struct{}

func New() *Action {
	return &Action{}
}

func (a *Action) Name() string {
	return "pwd"
}

func (a *Action) Description() string {
	return "生成随机字符串或密码"
}

func (a *Action) Help() string {
	return helpText
}

func (a *Action) Execute(
	_ context.Context,
	args []string,
	_ map[string]interface{},
) (map[string]interface{}, error) {

	if len(args) == 0 {
		return nil, fmt.Errorf(
			"missing subcommand (supported: rand, wg-keypair)",
		)
	}

	switch args[0] {

	case "rand":
		return generateRand(args[1:])

	case "wg-keypair":
		return generateWGKeypair()

	default:
		return nil, fmt.Errorf(
			"unknown subcommand: %s",
			args[0],
		)
	}
}
