// Package packetlogger logs all packets into a log file at
// "logs/Packet Logger/{region}_{UID}/{TIMESTAMP}.log".
// Warning, these can take up quite a lot of space over time and does not
// automatically rotate old logs.
package packetlogger

import (
	"bufio"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/kyoukaya/rhine/proxy"
	"github.com/kyoukaya/rhine/utils"

	"github.com/elazarl/goproxy"
)

const modName = "Packet Logger"

type rawPacketLoggerState struct {
	fileLogger *log.Logger
	buffer     *bufio.Writer
	*proxy.RhineModule
}

func (state *rawPacketLoggerState) handle(op string, data []byte, pktCtx *goproxy.ProxyCtx) []byte {
	go state.fileLogger.Printf("[%s] %s\n", op, string(data))
	return data
}

func (state *rawPacketLoggerState) Shutdown(bool) {
	state.Printf("Shutting down packetLogger for %d\n", state.UID)
	state.buffer.Flush()
}

func initFunc(mod *proxy.RhineModule) {
	dir := path.Join(utils.BinDir, "logs", modName, mod.Region, strconv.Itoa(mod.UID))
	err := os.MkdirAll(dir, 0755)
	utils.Check(err)
	now := time.Now()
	f, err := os.Create(path.Join(dir, now.Format("2006-01-02_15.04.05")+".log"))
	utils.Check(err)
	buffer := bufio.NewWriter(f)
	logger := log.New(buffer, "", log.Ltime)
	state := &rawPacketLoggerState{logger, buffer, mod}

	mod.OnShutdown(state.Shutdown)
	mod.Hook("*", 0, state.handle)
}

// Register hooks with dispatch
func init() {
	proxy.RegisterInitFunc(modName, initFunc)
}
