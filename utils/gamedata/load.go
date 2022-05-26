package gamedata

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"sort"
	"sync"

	"github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/utils"
	"github.com/kyoukaya/rhine/utils/gamedata/itemtable"
	"github.com/kyoukaya/rhine/utils/gamedata/stagetable"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
)

const (
	rawBaseURL     = "https://raw.githubusercontent.com/Kengxxiao/ArknightsGameData/master/"
	maxConnections = 4
)

var (
	fileList = []string{
		"%s/gamedata/excel/stage_table.json",
		"%s/gamedata/excel/item_table.json",
		"%s/gamedata/excel/character_table.json",
		"%s/gamedata/excel/gacha_table.json",
	}
	updateChecked = false
	// fileMutex is locked on program init
	fileMutex            sync.Mutex
	versionFileDelimiter = []byte(" ")
)

// updateGameData sends parallel GET requests to each individual file in fileList for all regions in regionMap.
// The number of parallel connections is limited by maxConnections which defaults to 4,
// ETags are sent with the GET requests where possible, which are read from the
// data/.version file, to minimize network traffic.
func updateGameData(l log.Logger) {
	var err error
	defer fileMutex.Unlock()
	if updateChecked {
		return
	}
	l.Println("Updating game data...")
	// Load .version file
	verMap := loadVersionFile()
	nFiles := len(regionMap) * len(fileList)

	// Start workers
	wg := sync.WaitGroup{}
	jobs := make(chan string, maxConnections)
	// Return channel for any errors
	errs := make(chan error, nFiles)
	// Return channel for the resulting path + etag string
	etags := make(chan string, nFiles)
	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go getAndUpdate(&wg, jobs, etags, errs, verMap)
	}
	// Send jobs to workers
	for _, region := range regionMap {
		for _, sFormat := range fileList {
			jobs <- fmt.Sprintf(sFormat, region)
		}
	}
	close(jobs)
	wg.Wait()
	close(errs)
	close(etags)

	for err := range errs {
		l.Warnln(err)
	}

	lines := make([]string, 0, nFiles)
	for etag := range etags {
		lines = append(lines, etag)
	}
	// Sort the lines for consistent output
	sort.Strings(lines)
	// write .version file
	f, err := os.Create(path.Join(utils.BinDir, "data/.version"))
	if err != nil {
		l.Warnln(err)
		return
	}
	for _, line := range lines {
		_, err = f.WriteString(line)
		if err != nil {
			l.Warnln(err)
			return
		}
		_, err = f.Write([]byte("\n"))
		if err != nil {
			l.Warnln(err)
			return
		}
	}
	updateChecked = true
	l.Println("Game data updated.")
}

// getAndUpdate is the worker thread spawned by updateGameData, it concatenates
// the raw base URL with the path sent in the jobs channel and GETs it with the
// corresponding etag from the verMap if it is available. If the response code is
// 200, the worker minifies the json and saves it to the disk.
func getAndUpdate(wg *sync.WaitGroup, jobs chan string, etags chan string,
	errs chan error, verMap map[string]string) {
	client := http.DefaultClient
	defer wg.Done()
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	for reqPath := range jobs {
		// do job
		req, err := http.NewRequest(http.MethodGet, rawBaseURL+reqPath, nil)
		if err != nil {
			errs <- err
			continue
		}
		etag := verMap[reqPath]
		if etag != "" {
			req.Header.Add("If-None-Match", etag)
		}
		resp, err := client.Do(req)
		if err != nil {
			errs <- err
			continue
		}
		if etag != "" && resp.StatusCode == 304 {
			etags <- reqPath + string(versionFileDelimiter) + etag
			continue
		}
		if resp.StatusCode != 200 {
			errs <- fmt.Errorf("Unexpected status %d when fetching %s",
				resp.StatusCode, rawBaseURL+reqPath)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errs <- err
			continue
		}
		fileName := path.Join(utils.BinDir, "data", reqPath)
		err = os.MkdirAll(path.Dir(fileName), 0755)
		if err != nil {
			errs <- err
			continue
		}
		f, err := os.Create(fileName)
		if err != nil {
			errs <- err
			continue
		}
		err = m.Minify("application/json", f, bytes.NewBuffer(body))
		if err != nil {
			errs <- err
			continue
		}
		etag = resp.Header.Get("ETag")
		etags <- reqPath + string(versionFileDelimiter) + etag
	}
}

// loadVersionFile loads the version file from the data folder and returns a map
// of the files to their etags.
func loadVersionFile() map[string]string {
	ret := make(map[string]string)
	bindir := utils.BinDir
	verPath := path.Join(bindir, "/data/.version")
	_, err := os.Stat(verPath)
	if os.IsNotExist(err) {
		return ret
	}
	f, err := os.Open(verPath)
	if err != nil {
		return ret
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		toks := bytes.Split(scanner.Bytes(), versionFileDelimiter)
		if len(toks) < 2 {
			// Malformed line, abort and return
			return ret
		}
		path := toks[0]
		etag := toks[1]
		ret[string(path)] = string(etag)
	}
	return ret
}

func loadExcelJSON(region, table string) []byte {
	fileMutex.Lock()
	defer fileMutex.Unlock()
	b, err := os.ReadFile(path.Join(utils.BinDir, "data", regionMap[region], "excel", table+".json"))
	utils.Check(err)
	return b
}

func (d *GameData) loadStageTable(region string) {
	b := loadExcelJSON(region, "stage_table")
	stageTable, err := stagetable.Unmarshal(b)
	utils.Check(err)
	state.stageTableMap[region] = &stageTable
}

func (d *GameData) loadItemTable(region string) {
	b := loadExcelJSON(region, "item_table")
	itemTable, err := itemtable.Unmarshal(b)
	utils.Check(err)
	state.itemTableMap[region] = &itemTable
}
