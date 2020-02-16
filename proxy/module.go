package proxy

import "github.com/kyoukaya/rhine/proxy/gamestate"

type RhineModule struct {
	Region    string
	UID       int
	GameState *gamestate.GameState

	name string
	*dispatch
	shutdownCB  ShutdownCb
	hooks       []*PacketHook
	initialized bool
}

type initFunc struct {
	name string
	fun  ModuleInitFunc
}

var (
	modules []initFunc
)

// ShutdownCb will be called when the proxy is shutting down or when a user reconnects.
type ShutdownCb func(shuttingDown bool)

// ModuleInitFunc will be executed when a user authenticates with the server to get
// initialized packethooks and the shutdown callback for a module.
type ModuleInitFunc func(*RhineModule)

// RegisterInitFunc adds a rhineModule that will be have its hook and shutdown generators run
// when a user authenticates with the game servers.
func RegisterInitFunc(name string, fun ModuleInitFunc) {
	modules = append(modules, initFunc{name: name, fun: fun})
}

type Hooker interface {
	Unhook()
}

func (m *RhineModule) Hook(target string, priority int, handler PacketHandler) Hooker {
	hook := &PacketHook{target, priority, handler, m}
	if m.initialized {
		m.Warnf("Failed to add hook %#v, module already initialized.", hook)
		return nil
	}
	m.hooks = append(m.hooks, hook)
	m.dispatch.insertHook(hook)
	return hook
}

func (m *RhineModule) OnShutdown(cb ShutdownCb) {
	m.shutdownCB = cb
}
