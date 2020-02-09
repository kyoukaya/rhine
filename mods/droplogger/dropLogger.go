// Package droplogger logs all map drops by appending them into a log file at
// "logs/Drop Logger/{region}_{UID}.log" and prints them to the attached logger.
package droplogger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	rhLog "github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/proxy"
	"github.com/kyoukaya/rhine/utils"
	"github.com/kyoukaya/rhine/utils/gamedata"
	"github.com/kyoukaya/rhine/utils/gamedata/itemtable"
	"github.com/kyoukaya/rhine/utils/gamedata/stagetable"

	"github.com/elazarl/goproxy"
	"github.com/tidwall/gjson"
)

const modName = "Drop Logger"

type modState struct {
	d           *proxy.Dispatch
	fileLogger  *log.Logger
	currStage   string
	stageStartT *time.Time
	isPractice  bool
	gd          *gamedata.GameData
	itemTable   *itemtable.ItemTable
	stageTable  *stagetable.StageTable
	rhLog.Logger
	mutex sync.Mutex
}

type logEntry struct {
	Ts      time.Time `json:"ts"`
	Rating  bool      `json:"ra"` // is 3star
	Rewards []reward  `json:"re"`
}

func (mod *modState) logDrops(rewards []reward, is3Star bool) {
	entry := logEntry{
		Ts:      time.Now(),
		Rating:  is3Star,
		Rewards: rewards,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		mod.Warn(err)
	}
	mod.fileLogger.Println(string(b))
}

func (mod *modState) battleFinishHandler(_ string, data []byte, _ *goproxy.ProxyCtx) []byte {
	mod.mutex.Lock()
	go mod.battleFinish(data)
	return data
}

func (mod *modState) battleFinish(data []byte) {
	defer mod.mutex.Unlock()
	if raw := gjson.GetBytes(data, "alert").Raw; raw != "" && raw != "[]" {
		mod.Warnf("Unexpected value when '[]' expected: %s", raw)
	}
	battle, err := unmarshalBattleFinish(data)
	if err != nil {
		mod.Warn(err)
	}
	if mod.isPractice || battle.ExpScale == 0 {
		return
	}
	var rewards []reward
	for _, r := range [][]reward{battle.Rewards, battle.UnusualRewards, battle.AdditionalRewards, battle.FurnitureRewards} {
		for _, x := range r {
			// Don't log gold or items with 0 count
			if x.Count > 0 && x.ID != "4001" {
				rewards = append(rewards, x)
			}
		}
	}
	sbuilder := strings.Builder{}
	for _, reward := range rewards {
		if reward.Count > 0 {
			sbuilder.WriteString("\"")
			sbuilder.WriteString(mod.itemTable.Items[reward.ID].Name)
			sbuilder.WriteString("\"x")
			sbuilder.WriteString(strconv.Itoa(int(reward.Count)))
			sbuilder.WriteString(" ")
		}
	}
	mod.logDrops(rewards, battle.ExpScale == 1.2)
	dropStr := strings.TrimRight(sbuilder.String(), " ")
	mod.Infof("Stage %s completed in %ds, drops: %s",
		mod.stageTable.Stages[mod.currStage].Code,
		int(time.Since(*mod.stageStartT).Seconds()),
		dropStr,
	)
}

func (mod *modState) battleStartRoutine(data []byte) {
	defer mod.mutex.Unlock()
	if mod.stageTable == nil || mod.itemTable == nil {
		mod.stageTable = mod.gd.GetStageInfo()
		mod.itemTable = mod.gd.GetItemInfo()
	}
	mod.isPractice = gjson.GetBytes(data, "usePracticeTicket").Bool()
	mod.currStage = gjson.GetBytes(data, "stageId").String()
}

func (mod *modState) battleStartHandler(op string, data []byte, ctx *goproxy.ProxyCtx) []byte {
	reqCtx := proxy.GetRequestContext(ctx)
	t := time.Now()
	mod.stageStartT = &t
	mod.mutex.Lock()
	go mod.battleStartRoutine(reqCtx.RequestData)
	return data
}

func initFunc(d *proxy.Dispatch) ([]*proxy.PacketHook, proxy.ShutdownCb) {
	dir := fmt.Sprintf("%s/logs/%s/", utils.BinDir, modName)
	err := os.MkdirAll(dir, 0755)
	utils.Check(err)
	f, err := os.OpenFile(fmt.Sprintf("%s%s_%s.log", dir, d.Region, strconv.Itoa(d.UID)),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	utils.Check(err)
	fileLogger := log.New(f, "", 0)
	gd, err := gamedata.New(d.Region, d.Logger)
	utils.Check(err)
	mod := modState{
		d:          d,
		fileLogger: fileLogger,
		gd:         gd,
		Logger:     d.Logger,
	}
	return []*proxy.PacketHook{
		proxy.NewPacketHook(modName, "S/quest/battleFinish", 0, mod.battleFinishHandler),
		proxy.NewPacketHook(modName, "S/quest/battleStart", 0, mod.battleStartHandler),
	}, func(bool) {}
}

func init() {
	proxy.RegisterMod(modName, initFunc)
}
