package main

import "context"

type Action interface {
	Name() string

	Execute(
		ctx context.Context,
		args []string, // 👈 新增这一行，用来承载游离的参数
		input map[string]interface{},
	) (
		map[string]interface{},
		error,
	)
}
