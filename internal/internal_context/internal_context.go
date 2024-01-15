package internal_context

type InternalContextIface interface {
	GetEnvList() []string
	GetEnv(envName string) string
	GetCurrentDir() string
}

type InternalContext struct{}

func (internalcontext *InternalContext) GetEnvList() []string {
	panic("not implemented") // TODO: Implement
}

func (internalcontext *InternalContext) GetEnv(envName string) string {
	panic("not implemented") // TODO: Implement
}

func (internalcontext *InternalContext) GetCurrentDir() string {
	panic("not implemented") // TODO: Implement
}
