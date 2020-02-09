package gamedata

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/utils"
	"github.com/kyoukaya/rhine/utils/gamedata/itemtable"
	"github.com/kyoukaya/rhine/utils/gamedata/stagetable"
	"github.com/tidwall/gjson"
)

const (
	excelPathFmt = "%s/data/%s/gamedata/excel/%s.json"
	apiBaseURL   = "https://api.github.com/repos/Kengxxiao/ArknightsGameData/"
	rawBaseURL   = "https://raw.githubusercontent.com/Kengxxiao/ArknightsGameData/master/"
)

var (
	fileList = []string{
		"%s/gamedata/excel/stage_table.json",
		"%s/gamedata/excel/item_table.json",
	}
	updateChecked = false
	// Locked on program init
	fileMutex sync.Mutex
)

func getLatestCommit() (string, error) {
	resp, err := http.Get(apiBaseURL + "commits")
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	s := gjson.GetBytes(b, "0.sha").String()
	return s, nil
}

// Caller is assumed to be holding the fileMutex lock
func isLocalDataCurrent(currentHash string) bool {
	if updated {
		return true
	}
	bindir := utils.BinDir
	verPath := bindir + "/data/.version"
	_, err := os.Stat(verPath)
	if os.IsNotExist(err) {
		return false
	}
	f, err := os.Open(verPath)
	utils.Check(err)
	b := make([]byte, 40)
	_, err = f.Read(b)
	utils.Check(err)
	return string(b) == currentHash
}

func getDataFile(dataPath string, done chan error) {
	var err error
	fileName := utils.BinDir + "/data/" + dataPath
	dirpath := path.Dir(fileName)
	err = os.MkdirAll(dirpath, 0755)
	if err != nil {
		done <- err
		return
	}
	f, err := os.Create(fileName)
	if err != nil {
		done <- err
		return
	}
	resp, err := http.Get(rawBaseURL + dataPath)
	if err != nil {
		done <- err
		return
	}
	if resp.StatusCode != 200 {
		done <- errors.New(resp.Status)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		done <- err
		return
	}
	_, err = f.Write(b)
	done <- err
}

func updateGameData(l log.Logger) {
	var err error
	defer fileMutex.Unlock()
	if updateChecked {
		return
	}
	lastCommit, err := getLatestCommit()
	if err != nil {
		l.Warn(err)
	}
	// Check cached data commit
	if isLocalDataCurrent(lastCommit) {
		return
	}
	l.Info("gamedata is outdated, updating...")
	done := make(chan error)
	doneCount := 0
	for _, region := range regionMap {
		for _, sFormat := range fileList {
			go getDataFile(fmt.Sprintf(sFormat, region), done)
			doneCount++
		}
	}
	for i := 0; i < doneCount; i++ {
		err = <-done
		if err != nil {
			l.Warn(err)
		}
	}
	// write .version file
	f, err := os.Create(utils.BinDir + "/data/.version")
	if err != nil {
		l.Warn(err)
		return
	}
	_, err = f.Write([]byte(lastCommit))
	if err != nil {
		l.Warn(err)
		return
	}
	updateChecked = true
	l.Info("gamedata updated")
}

func loadExcelJSON(region, table string) []byte {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	f, err := os.Open(fmt.Sprintf(
		excelPathFmt, utils.BinDir, regionMap[region], table))
	utils.Check(err)
	b, err := ioutil.ReadAll(f)
	utils.Check(err)
	return b
}

func (d *GameData) loadStageTable() {
	b := loadExcelJSON(d.region, "stage_table")
	stageTable, err := stagetable.Unmarshal(b)
	utils.Check(err)
	state.stageTableMap[d.region] = &stageTable
}

func (d *GameData) loadItemTable() {
	b := loadExcelJSON(d.region, "item_table")
	itemTable, err := itemtable.Unmarshal(b)
	utils.Check(err)
	state.itemTableMap[d.region] = &itemTable
}
