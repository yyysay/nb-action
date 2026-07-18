package actions

import (
	"github.com/yangtudou/nb-action/actions/bark"
	"github.com/yangtudou/nb-action/actions/pwd"
	"github.com/yangtudou/nb-action/actions/registry_sync"
	"github.com/yangtudou/nb-action/actions/test"
	"github.com/yangtudou/nb-action/core"
)

func RegisterAll(
	registry *core.Registry,
) {

	registry.Register(
		bark.New(),
		test.New(),
		pwd.New(),
		registry_sync.New(),
	)
}
