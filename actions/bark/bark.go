package bark

import (
	"context"
	"fmt"
)

type Bark struct{}

func New() *Bark {
	return &Bark{}
}

func (b *Bark) Name() string {
	return "bark"
}

func (b *Bark) Description() string {
	return "生成 Bark 加密通知"
}

func (b *Bark) Help() string {
	return helpText
}

func (b *Bark) Execute(
	ctx context.Context,
	args []string,
	input map[string]interface{},
) (
	map[string]interface{},
	error,
) {

	title, body, shouldPush := parseInput(
		args,
		input,
	)

	if body == "" {

		return nil,
			fmt.Errorf(
				"缺少推送内容！",
			)

	}

	payload := BarkPayload{

		Title: title,

		Body: body,
	}

	ciphertext, err := encrypt(
		payload,
	)

	if err != nil {

		return nil, err

	}

	result := map[string]interface{}{

		"ciphertext": ciphertext,
	}

	// --push 模式

	if shouldPush {

		config, err := LoadConfig()

		if err != nil {

			return nil, err

		}

		err = push(
			config,
			ciphertext,
		)

		if err != nil {

			return nil, err

		}

		result["pushed"] = true

	}

	return result, nil

}

func parseInput(
	args []string,
	input map[string]interface{},
) (
	string,
	string,
	bool,
) {

	var title string

	var body string

	var shouldPush bool

	// 解析命令行参数

	normalArgs := make(
		[]string,
		0,
	)

	for _, arg := range args {

		if arg == "--push" {

			shouldPush = true

			continue

		}

		normalArgs = append(
			normalArgs,
			arg,
		)

	}

	if len(normalArgs) >= 2 {

		title = normalArgs[0]

		body = normalArgs[1]

	} else if len(normalArgs) == 1 {

		body = normalArgs[0]

	}

	if title == "" {

		title, _ = input["title"].(string)

	}

	if body == "" {

		if value, ok := input["body"].(string); ok {

			body = value

		} else {

			body, _ = input["value"].(string)

		}

	}

	return title, body, shouldPush

}
