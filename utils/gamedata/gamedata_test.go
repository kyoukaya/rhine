package gamedata

import (
	"fmt"
	"testing"
	"time"

	"github.com/kyoukaya/rhine/utils"
)

type logShim struct{ t *testing.T }

func (l logShim) Flush()                              {}
func (l logShim) Println(i ...interface{})            { l.t.Log(i...) }
func (l logShim) Printf(s string, i ...interface{})   { l.t.Logf(s, i...) }
func (l logShim) Verboseln(i ...interface{})          { l.t.Log(i...) }
func (l logShim) Verbosef(s string, i ...interface{}) { l.t.Logf(s, i...) }
func (l logShim) Warnln(i ...interface{})             { l.t.Error(i...) }
func (l logShim) Warnf(s string, i ...interface{})    { l.t.Errorf(s, i...) }

func TestUpdateGameData(t *testing.T) {
	startT := time.Now()
	fileMutex.Lock()
	updateGameData(logShim{t})
	fmt.Printf("%dms\n", time.Since(startT).Milliseconds())
}

func TestStageTable(t *testing.T) {
	targetStage := "a001_05"
	targetRegion := "GL"
	d, err := New(targetRegion, logShim{t})
	utils.Check(err)
	table := d.GetStageInfo()
	_, exists := table.Stages[targetStage]
	if !exists {
		t.Error("Failed to find " + targetStage + " on " + targetRegion)
	}
}

func TestItemTable(t *testing.T) {
	targetItem := "1stact"
	targetRegion := "GL"
	d, err := New(targetRegion, logShim{t})
	utils.Check(err)
	table := d.GetItemInfo()
	item, exists := table.Items[targetItem]
	if !exists {
		t.Error("Failed to find " + targetItem + " on " + targetRegion)
	}
	t.Log(item)
}
