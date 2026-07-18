package bootstrap

import (
	"github.com/yangtudou/nb-action/actions"
	"github.com/yangtudou/nb-action/core"
)

func Load(runtime *core.Runtime) {

	actions.RegisterAll(
		runtime.Registry,
	)

}
