package registry_sync

type RegistrySync struct{}

func New() *RegistrySync {
	return &RegistrySync{}
}

func (r *RegistrySync) Name() string {
	return "registry-sync"
}

func (r *RegistrySync) Description() string {
	return "Docker 镜像仓库同步"
}

func (r *RegistrySync) Help() string {
	return helpText
}
