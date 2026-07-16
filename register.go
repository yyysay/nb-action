package main

import (
	"github.com/yangtudou/nb-action/actions"
)

func RegisterActions(
	server string,
	deviceKey string,
) {

	Register(
		actions.NewBark(
			server,
			deviceKey,
		),
		actions.NewTest(),
		actions.NewPassword(),
	)
}
