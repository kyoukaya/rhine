// Example is a binary that loads the droplogger and packetlogger for testing
// and development of Rhine and or mods for it.
package main

import (
	"log"

	_ "github.com/kyoukaya/rhine/mods/droplogger"
	_ "github.com/kyoukaya/rhine/mods/packetlogger"

	"github.com/kyoukaya/rhine/proxy"
)

var env string

func main() {
	logFlags := log.Llongfile | log.Ltime
	if env == "release" {
		logFlags = log.Lshortfile | log.Ltime
	}
	options := &proxy.Options{
		LoggerFlags:      logFlags,
		EnableHostFilter: true, // Enables the host filter
		HostFilter:       nil,  // Defaults to the built-in host filter
		Verbose:          true, // Prints additional VERB messages to the console
	}
	rhine := proxy.NewProxy(options)
	rhine.Start()
}
