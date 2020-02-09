/*
Package rhine provides a modular framework to interact with the game traffic of
Arknights.

Modules can only be registered at compile time with imports.
Each module should register itself by calling proxy.RegisterMod at program start up, either in the init() function or otherwise.
The initFunc provided will be called when a user authenticates with the game server to set up an instance of the module for that user.
While the ShutdownCb return value allows a module to clean up after itself when shutting down gracefully.
Here's an example:

	package yourmodule

	import (
		"github.com/elazarl/goproxy"
		"github.com/kyoukaya/rhine/proxy"
	)

	const modName = "yourModule"

	type modState struct{}

	func (modState) handle(op string, data []byte, pktCtx *goproxy.ProxyCtx) []byte {
		return data
	}

	func (modState) cleanUp(shuttingDown bool) {}

	func initFunc(d *proxy.Dispatch) ([]*proxy.PacketHook, proxy.ShutdownCb) {
		mod := modState{}
		return []*proxy.PacketHook{
			proxy.NewPacketHook(modName, "S/quest/battleStart", 0, mod.handle),
		}, mod.cleanUp
	}

	func init() {
		proxy.RegisterMod(modName, initFunc)
	}
*/
package rhine
