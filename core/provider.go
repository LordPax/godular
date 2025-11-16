package core

// IProvider defines the interface for a provider with initialization and module assignment capabilities.
type IProvider interface {
	GetName() string
	GetModule() IModule
	AssignModule(m IModule)
	OnInit() error
}

// Provider is a struct that implements the IProvider interface.
type Provider struct {
	name   string
	module IModule
}

// NewProvider creates a new instance of Provider with the given name.
func NewProvider(name string) *Provider {
	return &Provider{
		name: name,
	}
}

// GetName returns the name of the provider.
func (p *Provider) GetName() string {
	return p.name
}

// GetModule returns the module to which the provider belongs.
func (p *Provider) GetModule() IModule {
	return p.module
}

// AssignModule assigns a module to the provider.
func (p *Provider) AssignModule(m IModule) {
	p.module = m
}

// OnInit is a lifecycle hook that is called during module initialization.
func (p *Provider) OnInit() error {
	return nil
}
