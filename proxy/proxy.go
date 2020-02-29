package proxy

import (
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/proxy/filters"
	"github.com/kyoukaya/rhine/utils"

	"github.com/elazarl/goproxy"
)

var (
	// regionMap maps the TLD of the host string to their constant regional
	// representation: ["GL", "JP"]
	regionMap = map[string]string{
		"global": "GL",
		"jp":     "JP",
	}
	// onStartCbs will be called when proxy.Run() is called.
	onStartCbs []func(log.Logger)
)

const (
	certPath = "cert.pem"
	keyPath  = "key.pem"
)

// Options optionally changes the behavior of the proxy
type Options struct {
	Logger      log.Logger // Defaults to log.Log if not specified
	LoggerFlags int        // Flags to pass to the standard logger, if a custom logger is not specified
	// LogPath defaults to "logs/proxy.log", setting it to "/dev/null", even on Windows,
	// will make the logger not output a file.
	LogPath          string
	LogDisableStdOut bool           // Should stdout output be DISABLED for the default logger
	EnableHostFilter bool           // Filters out packets from certain hosts if they match HostFilter
	HostFilter       *regexp.Regexp // Custom regexp filter for filtering packets, defaults to the block list in proxy/filters.go
	Verbose          bool           // log more Rhine information
	VerboseGoProxy   bool           // log every GoProxy request to stdout
	Address          string         // proxy listen address, defaults to ":8080"
	DisableCertStore bool           // Disables the built in certstore, reduces memory usage but increases HTTP latency and CPU usage.
}

// Proxy contains the internal state relevant to the proxy
type Proxy struct {
	mutex      *sync.Mutex
	server     *goproxy.ProxyHttpServer
	hostFilter *regexp.Regexp
	options    *Options
	// dispatches contains a mapping of a user's UID and region in string form
	// to the user's Dispatch.
	dispatches map[string]*dispatch
	log.Logger
}

// RegisterInitFunc adds a rhineModule that will be initialized when a user authenticates
// with the Arknights server.
func RegisterInitFunc(name string, fun ModuleInitFunc) {
	modules = append(modules, initFunc{name: name, fun: fun})
}

// OnStart registers a function to be called back when the proxy is initialized, i.e.,
// when the proxy server is ready, not when an Arknights user is connected. The
// Logger interface provided will be the proxy's logger.
func OnStart(cb func(log.Logger)) {
	onStartCbs = append(onStartCbs, cb)
}

// NewProxy returns a new initialized Dispatch
func NewProxy(options *Options) *Proxy {
	logger := options.Logger
	if logger == nil {
		logger = log.New(!options.LogDisableStdOut, options.Verbose, options.LogPath, options.LoggerFlags)
	}
	if options.Address == "" {
		options.Address = ":8080"
	}
	var proxyFilter *regexp.Regexp = nil
	if options.EnableHostFilter {
		if options.HostFilter == nil {
			proxyFilter = filters.HostFilter
		} else {
			proxyFilter = options.HostFilter
		}
	}

	server := goproxy.NewProxyHttpServer()
	if !options.DisableCertStore {
		server.CertStore = newCertStore(logger)
	}

	server.Logger = printfFunc(logShim(logger))
	server.Verbose = options.VerboseGoProxy
	proxy := &Proxy{
		mutex:      &sync.Mutex{},
		server:     server,
		options:    options,
		Logger:     logger,
		dispatches: make(map[string]*dispatch),
		hostFilter: proxyFilter,
	}
	server.OnRequest().DoFunc(proxy.HandleReq)
	server.OnResponse().DoFunc(proxy.HandleResp)

	_, certStatErr := os.Stat(utils.BinDir + certPath)
	_, keyStatErr := os.Stat(utils.BinDir + keyPath)
	// Generate CA if it doesn't exist
	if os.IsNotExist(certStatErr) || os.IsNotExist(keyStatErr) {
		proxy.Printf("Generating CA...")
		if err := utils.GenerateCA(certPath, keyPath); err != nil {
			proxy.Warnln(err)
			panic(err)
		}
		proxy.Printf("Copy and register the created 'cert.pem' with your client.")
	}
	if err := utils.LoadCA(certPath, keyPath); err != nil {
		proxy.Warnln(err)
		panic(err)
	}
	server.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(proxy.httpsHandler))
	return proxy
}

// Interface shim for goproxy.Logger
func logShim(logger log.Logger) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		logger.Printf(format, v...)
	}
}

type printfFunc func(format string, v ...interface{})

func (f printfFunc) Printf(format string, v ...interface{}) {
	f("[goproxy] "+format, v...)
}

// HTTPSHandler to allow HTTPS connections to pass through the proxy without being
// MITM'd.
func (p *Proxy) httpsHandler(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	if p.hostFilter != nil && p.hostFilter.MatchString(host) {
		p.Verbosef("==== Rejecting %v", host)
		return goproxy.RejectConnect, host
	}
	return goproxy.MitmConnect, host
}

// Start starts the proxy. This is blocking and does not return.
func (p *Proxy) Start() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// Catch sigint/sigterm and cleanly exit
	go func() {
		<-sigs
		p.Printf("Shutting down.\n")
		p.Flush()
		p.Shutdown()
		os.Exit(0)
	}()

	ipstring := utils.GetOutboundIP()
	addrSplit := strings.Split(p.options.Address, ":")
	if len(addrSplit) == 2 {
		ipstring += ":" + addrSplit[1]
	}

	for _, cb := range onStartCbs {
		cb(p.Logger)
	}

	p.Printf("proxy server listening on %s", ipstring)
	err := http.ListenAndServe(p.options.Address, p.server)
	p.Warnln(err)
	panic(err)
}

// Shutdown calls Shutdown on all modules for all users.
func (p *Proxy) Shutdown() {
	for _, dispatch := range p.dispatches {
		for _, mod := range dispatch.modules {
			if mod.shutdownCB != nil {
				mod.shutdownCB(true)
			}
		}
	}
}

// getUser returns a Dispatch for the specified UID
func (p *Proxy) getUser(UID, region string) *dispatch {
	rUID := region + "_" + UID
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.dispatches[rUID]
}

// addUser records a user's information indexed by their UID, if a record belonging to
// the specified UID already exists, its hooks will be shutdown and the record will be overwritten.
func (p *Proxy) addUser(UID, region string) *dispatch {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	rUID := region + "_" + UID

	if dispatch, exists := p.dispatches[rUID]; exists {
		p.Printf("%s reconnecting. Shutting down mods.", rUID)
		for _, module := range dispatch.modules {
			if module.shutdownCB != nil {
				module.shutdownCB(true)
			}
		}
	} else {
		p.Printf("User %s logged in", rUID)
	}

	UIDint, err := strconv.Atoi(UID)
	utils.Check(err)
	d := &dispatch{
		mutex:  &sync.Mutex{},
		uid:    UIDint,
		region: region,
		hooks:  make(map[string][]*PacketHook),
		Logger: p.Logger,
	}
	d.initMods(modules)
	p.dispatches[rUID] = d
	return d
}
