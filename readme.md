# rhine

rhine is a modular framework for intercepting, processing, and dispatching game traffic for Arknight's global, Japanese and Korean servers by means of a HTTPS proxy.
Another core functionality of rhine is to mirror the game state of the client by inspection of HTTPS traffic, allowing modules insight into the exact state and changes to the state of the client.
Optionally, rhine can also block requests to telemetry or ad domains frequently contacted when using an android emulator.

## Usage

While rhine is intended to be used as a framework on which developers can write their own programs, an example program is provided as [`cmd/example/rhine.go`](https://github.com/kyoukaya/rhine/blob/master/cmd/example/rhine.go) which initializes the `packetlogger` and `droplogger` modules so that developers can give it a spin.
Run `go build cmd/example/rhine.go && ./rhine.exe` to build and run the proxy server, and then direct your client to use it.
You will be required to install the generated root CA on your emulator/device so that rhine will be able to listen in on the HTTPS game traffic.

## Example Modules

The 2 provided example modules in this repository are pretty self explanatory, `packetlogger` logs the raw body of each game packet, while `droplogger` logs the drops from each battle.

Besides the modules provided in this repository, you can also try out:
- [ak-discordrpc](https://github.com/kyoukaya/ak-discordrpc) - a Discord rich presence client for Arknights.
- [angelina](https://github.com/kyoukaya/angelina) - websocket interface for rhine allowing packet and gamestate hooks.

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

## Background

A lot of the code for rhine came from [Hoxy](https://github.com/kyoukaya/hoxy), a previous attempt at this concept which didn't work out so well as it tried to marshal every single packet sent and received by the client, which caused many development problems.
The situation in Arknights is however slightly different, as most packets received by the client affect the client in a programmatic way.
Arknight's game servers send delta patches in the form of JSON objects to modify and delete aspects of the game state.
This lets us model mirror the client's game state a lot easier without perfectly marshalling all the packets.

## Contact
If you're interested in rhine and want to talk to me about it, message me on Discord @kyoukaya#6240, or join the [discord server](https://discord.gg/zXhE7vA)!
