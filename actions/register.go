package actions

import (
	"github.com/yangtudou/nb-action/actions/pwd"
	"github.com/yangtudou/nb-action/actions/test"
	"github.com/yangtudou/nb-action/core"
)

func RegisterAll(
	registry *core.Registry,
	server string,
	deviceKey string,
) {

	registry.Register(
		NewBark(
			server,
			deviceKey,
		),
		test.New(),
		pwd.New(),
		&RegistrySync{},
	)
}
