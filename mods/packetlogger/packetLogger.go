// Package packetlogger logs all packets into a log file at
// "logs/Packet Logger/{region}_{UID}/{TIMESTAMP}.log".
// Warning, these can take up quite a lot of space over time and does not
// automatically rotate old logs.
package packetlogger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	rhLog "github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/proxy"
	"github.com/kyoukaya/rhine/utils"

	"github.com/elazarl/goproxy"
)

const modName = "Packet Logger"

type rawPacketLoggerState struct {
	fileLogger *log.Logger
	buffer     *bufio.Writer
	d          *proxy.Dispatch
	rhLog.Logger
}

func (state *rawPacketLoggerState) handle(op string, data []byte, pktCtx *goproxy.ProxyCtx) []byte {
	go state.fileLogger.Printf("[%s] %s\n", op, string(data))
	return data
}

func (state *rawPacketLoggerState) Shutdown(bool) {
	state.Printf("Shutting down packetLogger for %d\n", state.d.UID)
	state.buffer.Flush()
}

// Register hooks with dispatch
func init() {
	initFunc := func(d *proxy.Dispatch) ([]*proxy.PacketHook, proxy.ShutdownCb) {
		dir := fmt.Sprintf("%s/logs/%s/%s_%s/", utils.BinDir, modName, d.Region, strconv.Itoa(d.UID))
		err := os.MkdirAll(dir, 0755)
		utils.Check(err)
		now := time.Now()
		f, err := os.Create(fmt.Sprintf("%s%s.log", dir, now.Format("2006-01-02_15.04.05")))
		utils.Check(err)
		buffer := bufio.NewWriter(f)
		logger := log.New(buffer, "", log.Ltime)

		mod := &rawPacketLoggerState{logger, buffer, d, d.Logger}
		hooks := []*proxy.PacketHook{
			proxy.NewPacketHook(modName, "*", 0, mod.handle),
		}
		return hooks, mod.Shutdown
	}
	proxy.RegisterMod(modName, initFunc)
}
