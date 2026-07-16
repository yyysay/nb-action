package pwd

import (
	"context"
	"fmt"
)

type Password struct{}

func New() *Password {
	return &Password{}
}

func (p *Password) Name() string {
	return "pwd"
}

func (p *Password) Description() string {
	return "生成随机字符串或密码"
}

func (p *Password) Help() string {
	return helpText
}

func (p *Password) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (map[string]interface{}, error) {

	if len(args) == 0 {
		return nil, fmt.Errorf("missing subcommand")
	}

	switch args[0] {

	case "wg-keypair":
		return generateWGKeypair()

	case "rand":
		return generateRand(args[1:])

	default:
		return nil, fmt.Errorf(
			"unknown subcommand: %s",
			args[0],
		)
	}
}
