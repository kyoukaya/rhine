package proxy

type rhineModule struct {
	name     string
	initFunc ModuleInitFunc
}

var (
	modules []*rhineModule
)

// ShutdownCb will be called when the proxy is shutting down or when a user reconnects.
type ShutdownCb func(shuttingDown bool)

// ModuleInitFunc will be executed when a user authenticates with the server to get
// initialized packethooks and the shutdown callback for a module.
type ModuleInitFunc func(d *Dispatch) ([]*PacketHook, ShutdownCb)

// RegisterMod adds a rhineModule that will be have its hook and shutdown generators run
// when a user authenticates with the game servers.
func RegisterMod(name string, initFunc ModuleInitFunc) {
	modules = append(modules, &rhineModule{
		name:     name,
		initFunc: initFunc,
	})
}
