package core

import (
	"fmt"

	"github.com/dominikbraun/graph"
)

var (
	moduleSeen     []string
	moduleSeenCall []string
	CalledEvent    bool
	g              = graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())
)

type FuncCb func() error

// IModule defines the interface for a module that can manage providers and dependencies.
type IModule interface {
	GetName() string
	Get(provider string) IProvider
	AddProvider(p IProvider)
	AddModule(dep IModule)
	Init() error
	getProvider(provider string) IProvider
	On(event string, cb FuncCb)
	callEvent(event string) error
	Call(event string) error
	List() []string
	Graph() graph.Graph[string, string]
}

// Module is a struct that implements the IModule interface, managing providers and dependencies.
type Module struct {
	name     string
	module   map[string]IModule
	provider map[string]IProvider
	events   map[string]FuncCb
}

// NewModule creates a new instance of Module with the given name.
func NewModule(name string) *Module {
	return &Module{
		name:     name,
		module:   make(map[string]IModule),
		provider: make(map[string]IProvider),
		events:   make(map[string]FuncCb),
	}
}

// GetName returns the name of the module.
func (m *Module) GetName() string {
	return m.name
}

// Get retrieves a provider by its name, searching within the module and its dependencies.
func (m *Module) Get(provider string) IProvider {
	if p := m.getProvider(provider); p != nil {
		return p
	}

	for _, d := range m.module {
		if p := d.getProvider(provider); p != nil {
			return p
		}
	}

	panic(fmt.Sprintf("Provider %s not found", provider))
}

// getProvider is a helper method that retrieves a provider by its name within the current module.
func (m *Module) getProvider(provider string) IProvider {
	return m.provider[provider]
}

// AddProvider adds a provider to the module and assigns the module to the provider.
func (m *Module) AddProvider(p IProvider) {
	m.provider[p.GetName()] = p
	p.AssignModule(m)
}

// AddModule adds a dependency module to the current module.
func (m *Module) AddModule(dep IModule) {
	m.module[dep.GetName()] = dep
}

// Init initializes the module and its dependencies, ensuring each module is initialized only once.
func (m *Module) Init() error {
	for _, d := range m.module {
		if inArray(d.GetName(), moduleSeen) {
			continue
		}

		moduleSeen = append(moduleSeen, d.GetName())

		if err := d.Init(); err != nil {
			return err
		}
	}

	fmt.Printf("==> Initializing %s\n", m.GetName())

	for _, p := range m.provider {
		if err := p.OnInit(); err != nil {
			return err
		}
	}

	return nil
}

// On registers a callback function for a specific event.
func (m *Module) On(event string, cb FuncCb) {
	m.events[event] = cb
}

func (m *Module) callEvent(event string) error {
	cb, ok := m.events[event]
	if !ok {
		return nil
	}
	CalledEvent = true
	return cb()
}

// Call triggers the callback function associated with a specific event.
func (m *Module) Call(event string) error {
	for _, d := range m.module {
		if inArray(d.GetName(), moduleSeenCall) {
			continue
		}

		moduleSeenCall = append(moduleSeenCall, d.GetName())

		if err := d.Call(event); err != nil {
			return err
		}
	}

	if err := m.callEvent(event); err != nil {
		return err
	}

	return nil
}

// List returns a list of event names registered in the module and its dependencies.
func (m *Module) List() []string {
	var list []string

	for _, d := range m.module {
		if inArray(d.GetName(), moduleSeenCall) {
			continue
		}

		moduleSeenCall = append(moduleSeenCall, d.GetName())
		list = append(list, d.List()...)
	}

	for event := range m.events {
		line := "- " + event + " (" + m.GetName() + ")"
		list = append(list, line)
	}

	return list
}

func (m *Module) Graph() graph.Graph[string, string] {
	g.AddVertex(
		m.GetName(),
		graph.VertexAttribute("colorscheme", "blues3"),
		graph.VertexAttribute("style", "filled"),
		graph.VertexAttribute("color", "2"),
		graph.VertexAttribute("fillcolor", "1"),
		graph.VertexAttribute("shape", "box"),
		graph.VertexAttribute("group", "0"),
	)

	for _, d := range m.module {
		d.Graph()
		g.AddEdge(d.GetName(), m.GetName())
	}

	return g
}
