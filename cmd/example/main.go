// Example is a binary that loads the droplogger and packetlogger for testing
// and development of Rhine and or mods for it.
package main

import (
	"flag"
	"log"

	_ "github.com/kyoukaya/rhine/mods/droplogger"
	_ "github.com/kyoukaya/rhine/mods/packetlogger"

	"github.com/kyoukaya/rhine/proxy"
)

var env string

func main() {
	logPath := flag.String("log-path", "logs/proxy.log", "file to output the log to")
	silent := flag.Bool("silent", false, "don't print anything to stdout")
	filter := flag.Bool("filter", false, "enable the host filter")
	verbose := flag.Bool("v", false, "print Rhine verbose messages")
	verboseGoProxy := flag.Bool("v-goproxy", false, "print verbose goproxy messages")
	host := flag.String("host", ":8080", "hostname:port")
	flag.Parse()

	logFlags := log.Llongfile | log.Ltime
	if env == "release" {
		logFlags = log.Lshortfile | log.Ltime
	}
	options := &proxy.Options{
		LogPath:          *logPath,
		LogDisableStdOut: *silent,
		EnableHostFilter: *filter,
		LoggerFlags:      logFlags,
		Verbose:          *verbose,
		VerboseGoProxy:   *verboseGoProxy,
		Address:          *host,
	}
	rhine := proxy.NewProxy(options)
	rhine.Start()
}
