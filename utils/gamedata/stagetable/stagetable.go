package stagetable

import "encoding/json"

func Unmarshal(data []byte) (StageTable, error) {
	var r StageTable
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *StageTable) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type StageTable struct {
	Stages              map[string]Stage          `json:"stages"`
	Campaigns           Campaigns                 `json:"campaigns"`
	CampaignGroups      CampaignGroups            `json:"campaignGroups"`
	RuneStageGroups     RuneStageGroups           `json:"runeStageGroups"`
	MapThemes           map[string]MapTheme       `json:"mapThemes"`
	TileInfo            map[string]TileInfo       `json:"tileInfo"`
	ForceOpenTable      map[string]ForceOpenTable `json:"forceOpenTable"`
	TimelyStageDropInfo RuneStageGroups           `json:"timelyStageDropInfo"`
	TimelyTable         RuneStageGroups           `json:"timelyTable"`
	StageValidInfo      map[string]StageValidInfo `json:"stageValidInfo"`
}

type CampaignGroups struct {
	CampG1 CampG `json:"camp_g_1"`
	CampG2 CampG `json:"camp_g_2"`
}

type CampG struct {
	GroupID     string   `json:"groupId"`
	ActiveCamps []string `json:"activeCamps"`
	StartTs     int64    `json:"startTs"`
	EndTs       int64    `json:"endTs"`
}

type Campaigns struct {
	Camp01 Camp0 `json:"camp_01"`
	Camp02 Camp0 `json:"camp_02"`
}

type Camp0 struct {
	StageID              string          `json:"stageId"`
	GainLadders          []GainLadder    `json:"gainLadders"`
	BreakLadders         []BreakLadder   `json:"breakLadders"`
	DropLadders          []DropLadder    `json:"dropLadders"`
	DisplayRewards       []DisplayReward `json:"displayRewards"`
	DisplayDetailRewards []DisplayReward `json:"displayDetailRewards"`
}

type BreakLadder struct {
	KillCnt     int64    `json:"killCnt"`
	BreakFeeAdd int64    `json:"breakFeeAdd"`
	Rewards     []Reward `json:"rewards"`
}

type Reward struct {
	ID    string     `json:"id"`
	Count int64      `json:"count"`
	Type  RewardType `json:"type"`
}

type DisplayReward struct {
	OccPercent *int64                  `json:"occPercent,omitempty"`
	Type       DisplayDetailRewardType `json:"type"`
	ID         string                  `json:"id"`
	DropType   int64                   `json:"dropType"`
}

type DropLadder struct {
	KillCnt  int64    `json:"killCnt"`
	DropInfo DropInfo `json:"dropInfo"`
}

type DropInfo struct {
	// FirstPassRewards     interface{}     `json:"firstPassRewards"`
	// PassRewards          interface{}     `json:"passRewards"`
	DisplayDetailRewards []DisplayReward `json:"displayDetailRewards"`
}

type GainLadder struct {
	KillCnt      int64 `json:"killCnt"`
	ApFailReturn int64 `json:"apFailReturn"`
	Favor        int64 `json:"favor"`
	ExpGain      int64 `json:"expGain"`
	GoldGain     int64 `json:"goldGain"`
}

type ForceOpenTable struct {
	ID            string   `json:"id"`
	StartTime     int64    `json:"startTime"`
	EndTime       int64    `json:"endTime"`
	ForceOpenList []string `json:"forceOpenList"`
}

type MapTheme struct {
	ThemeID   string  `json:"themeId"`
	UnitColor string  `json:"unitColor"`
	ThemeType *string `json:"themeType"`
}

type RuneStageGroups struct {
}

type StageValidInfo struct {
	StartTs int64 `json:"startTs"`
	EndTs   int64 `json:"endTs"`
}

type Stage struct {
	StageType          StageType         `json:"stageType"`
	Difficulty         Difficulty        `json:"difficulty"`
	UnlockCondition    []UnlockCondition `json:"unlockCondition"`
	StageID            string            `json:"stageId"`
	LevelID            string            `json:"levelId"`
	ZoneID             string            `json:"zoneId"`
	Code               string            `json:"code"`
	Name               *string           `json:"name"`
	Description        string            `json:"description"`
	HardStagedID       *string           `json:"hardStagedId"`
	DangerLevel        string            `json:"dangerLevel"`
	DangerPoint        float64           `json:"dangerPoint"`
	CanPractice        bool              `json:"canPractice"`
	CanBattleReplay    bool              `json:"canBattleReplay"`
	ApCost             int64             `json:"apCost"`
	ApFailReturn       int64             `json:"apFailReturn"`
	EtCost             int64             `json:"etCost"`
	EtFailReturn       int64             `json:"etFailReturn"`
	PracticeTicketCost int64             `json:"practiceTicketCost"`
	ExpGain            int64             `json:"expGain"`
	GoldGain           int64             `json:"goldGain"`
	PassFavor          int64             `json:"passFavor"`
	CompleteFavor      int64             `json:"completeFavor"`
	SlProgress         int64             `json:"slProgress"`
	DisplayMainItem    *string           `json:"displayMainItem"`
	HilightMark        bool              `json:"hilightMark"`
	BossMark           bool              `json:"bossMark"`
	IsStoryOnly        bool              `json:"isStoryOnly"`
	StageDropInfo      StageDropInfo     `json:"stageDropInfo"`
	MainStageID        *string           `json:"mainStageId"`
}

type StageDropInfo struct {
	DisplayRewards       []DisplayReward `json:"displayRewards"`
	DisplayDetailRewards []DisplayReward `json:"displayDetailRewards"`
}

type UnlockCondition struct {
	StageID       string `json:"stageId"`
	CompleteState int64  `json:"completeState"`
}

type TileInfo struct {
	TileKey      string `json:"tileKey"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsFunctional bool   `json:"isFunctional"`
}

type RewardType string

const (
	HggShd        RewardType = "HGG_SHD"
	LggShd        RewardType = "LGG_SHD"
	PurpleCARDEXP RewardType = "CARD_EXP"
	PurpleGOLD    RewardType = "GOLD"
)

type DisplayDetailRewardType string

const (
	ActivityCoin  DisplayDetailRewardType = "ACTIVITY_COIN"
	Char          DisplayDetailRewardType = "CHAR"
	Diamond       DisplayDetailRewardType = "DIAMOND"
	DiamondShd    DisplayDetailRewardType = "DIAMOND_SHD"
	FluffyCARDEXP DisplayDetailRewardType = "CARD_EXP"
	FluffyGOLD    DisplayDetailRewardType = "GOLD"
	Furn          DisplayDetailRewardType = "FURN"
	Material      DisplayDetailRewardType = "MATERIAL"
	TktRecruit    DisplayDetailRewardType = "TKT_RECRUIT"
)

type Difficulty string

const (
	FourStar Difficulty = "FOUR_STAR"
	Normal   Difficulty = "NORMAL"
)

type LoadingPicID string

const (
	Loading1  LoadingPicID = "loading1"
	Loading2  LoadingPicID = "loading2"
	Loading3  LoadingPicID = "loading3"
	Loading4  LoadingPicID = "loading4"
	LoadingE1 LoadingPicID = "loadingE1"
	LoadingE2 LoadingPicID = "loadingE2"
	LoadingS  LoadingPicID = "loadingS"
)

type StageType string

const (
	Activity StageType = "ACTIVITY"
	Campaign StageType = "CAMPAIGN"
	Daily    StageType = "DAILY"
	Guide    StageType = "GUIDE"
	Main     StageType = "MAIN"
	Sub      StageType = "SUB"
)
