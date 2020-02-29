# Rhine

Rhine is a modular framework for intercepting, processing, and dispatching game traffic for Arknight's global, Japanese and Korean servers by means of a HTTPS proxy.
Optionally, Rhine can also block requests to telemetry or ad domains frequently contacted when using an android emulator.

## Usage

While Rhine is intended to be used as a framework on which developers can write their own programs, an example program is provided as [`cmd/example/rhine.go`](https://github.com/kyoukaya/rhine/blob/master/cmd/example/rhine.go) which initializes the `packetlogger` and `droplogger` modules so that developers can give it a spin.
Execute `go run cmd/example/rhine.go` to start the proxy server, and then direct your client to use it. You will be required to install a generated root CA on your emulator/device so that Rhine will be able to listen in on the HTTPS game traffic.

## Modules

The 2 provided example modules in this repository are pretty self explanatory, `packetlogger` logs the raw body of each game packet, while `droplogger` logs the drops from each battle.
Besides the modules provided in this repository, you can also try out the [ak-discordrpc](https://github.com/kyoukaya/ak-discordrpc) module.

## Development

Modules can only be registered at compile time with imports.
Each module should register itself by calling `proxy.RegisterMod(modName, initFunc)` at program start up, either in the `init()` function or otherwise.
The `initFunc` provided will be called when a user authenticates with the game server to set up an instance of the module for that user.
While the `ShutdownCb` return value allows a module to clean up after itself when shutting down gracefully.

Here's an example:

```go
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

func initFunc(mod *proxy.RhineModule) {
	state := modState{}

	mod.OnShutdown(state.cleanUp)
	mod.Hook("S/quest/battleStart", 0, mod.handle)
}

func init() {
	proxy.RegisterMod(modName, initFunc)
}
```

## Background and Future Plans

A lot of the code for Rhine came from [Hoxy](https://github.com/kyoukaya/hoxy), a previous attempt at this concept which didn't work out so well as it tried to marshal every single packet sent and received by the client, which caused many development problems. The situation in Arknights is however slightly different, as most packets received by the client affect the client in a very programmatic way, it essentially sends delta patches in the form of JSON objects. This lets us model mirror the client's game state a lot easier without perfectly marshalling all the packets.

Mirroring the client's game state is exactly the goal of the gamestate module, though at the moment it only unmarshalls the initial synchronization packet and stores it. A good amount of reflection will be required to apply the delta patches sent by the server, but it's high on my priority. Once, completed, gamestate would enable a KancolleViewer-like application which is always something I wanted to do as well.

Higher on my priority though, is porting [ArkPlanner](https://github.com/ycremar/ArkPlanner) to go, combined with the data available to Rhine, mods will make one's farm a lot easier than inputting it into the website.

If you're interested in Rhine and want to talk to me about it, message me on discord @ kyoukaya#6240. Might make a server if there's any interest.
