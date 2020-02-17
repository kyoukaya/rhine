package statestruct

type User struct {
	Building *Building `json:"building"`
	Social   *struct {
		AssistCharList  []AssistCharList `json:"assistCharList"`
		YesterdayReward YesterdayReward  `json:"yesterdayReward"`
	} `json:"social"`
	Ticket     interface{} `json:"ticket"`
	Gacha      *Gacha      `json:"gacha"`
	OpenServer *struct {
		CheckIn    OpenServerCheckIn `json:"checkIn"`
		ChainLogin ChainLogin        `json:"chainLogin"`
	} `json:"openServer"`
	DexNav  *DexNav  `json:"dexNav"`
	Dungeon *Dungeon `json:"dungeon"`
	Shop    *Shop    `json:"shop"`
	Skin    *struct {
		CharacterSkins map[string]int64 `json:"characterSkins"`
	} `json:"skin"`
	PushFlags *PushFlags   `json:"pushFlags"`
	Troop     *Troop       `json:"troop"`
	CheckIn   *UserCheckIn `json:"checkIn"`
	Activity  *struct {
		Default     interface{}            `json:"DEFAULT"`
		MissionOnly map[string]interface{} `json:"MISSION_ONLY"`
	} `json:"activity"`
	Mission          *Mission `json:"mission"`
	CollectionReward *struct {
		Team map[string]int64 `json:"team"`
	} `json:"collectionReward"`
	Recruit    *Recruit                             `json:"recruit"`
	Status     *UserStatus                          `json:"status"`
	Consumable map[string]map[string]ConsumableInfo `json:"consumable"`
	Inventory  map[string]int64                     `json:"inventory"`
	Event      *struct {
		Building int64 `json:"building"`
	} `json:"event"`
}

type Building struct {
	Status *struct {
		Labor Labor `json:"labor"`
	} `json:"status"`
	Chars     map[string]BuildingChar `json:"chars"`
	RoomSlots map[string]RoomSlot     `json:"roomSlots"`
	Rooms     *Rooms                  `json:"rooms"`
	Furniture map[string]struct {
		Count int64 `json:"count"`
		InUse int64 `json:"inUse"`
	} `json:"furniture"`
	DiyPresetSolutions interface{} `json:"diyPresetSolutions"`
	Assist             []int64     `json:"assist"`
}

type RoomSlot struct {
	Level                 int64   `json:"level"`
	State                 int64   `json:"state"`
	RoomID                string  `json:"roomId"`
	CharInstIDS           []int64 `json:"charInstIds"`
	CompleteConstructTime int64   `json:"completeConstructTime"`
}

type BuildingChar struct {
	CharID        string `json:"charId"`
	LastApAddTime int64  `json:"lastApAddTime"`
	Ap            int64  `json:"ap"`
	RoomSlotID    string `json:"roomSlotId"`
	Index         int64  `json:"index"`
	ChangeScale   int64  `json:"changeScale"`
	Bubble        struct {
		Normal Assist `json:"normal"`
		Assist Assist `json:"assist"`
	} `json:"bubble"`
	WorkTime int64 `json:"workTime"`
}

type Assist struct {
	Add int64 `json:"add"`
	Ts  int64 `json:"ts"`
}

type Rooms struct {
	Control     map[string]ControlSlot     `json:"CONTROL"`
	Elevator    map[string]interface{}     `json:"ELEVATOR"`
	Power       map[string]PowerInfo       `json:"POWER"`
	Manufacture map[string]ManufactureInfo `json:"MANUFACTURE"`
	Trading     map[string]TradingInfo     `json:"TRADING"`
	Dormitory   map[string]DormInfo        `json:"DORMITORY"`
	Corridor    map[string]interface{}     `json:"CORRIDOR"`
	Workshop    map[string]WorkshopInfo    `json:"WORKSHOP"`
	Meeting     map[string]MeetingInfo     `json:"MEETING"`
	Hire        map[string]HireInfo        `json:"HIRE"`
}

type ControlSlot struct {
	Buff struct {
		Global struct {
			ApCost int64 `json:"apCost"`
		} `json:"global"`
		Manufacture SpeedContainer `json:"manufacture"`
		Trading     SpeedContainer `json:"trading"`
		ApCost      interface{}    `json:"apCost"`
	} `json:"buff"`
	ApCost int64 `json:"apCost"`
}

type SpeedContainer struct {
	Speed float64 `json:"speed"`
}

type DormInfo struct {
	Buff        DormInfoBuff `json:"buff"`
	Comfort     int64        `json:"comfort"`
	DiySolution DiySolution  `json:"diySolution"`
}

type DormInfoBuff struct {
	ApCost struct {
		All    int64 `json:"all"`
		Single struct {
			Target *int64 `json:"target"`
			Value  int64  `json:"value"`
		} `json:"single"`
		Self interface{} `json:"self"`
	} `json:"apCost"`
}

type DiySolution struct {
	WallPaper *string  `json:"wallPaper"`
	Floor     *string  `json:"floor"`
	Carpet    []Carpet `json:"carpet"`
	Other     []Carpet `json:"other"`
}

type Carpet struct {
	ID         string     `json:"id"`
	Coordinate Coordinate `json:"coordinate"`
}

type Coordinate struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type HireInfo struct {
	Buff             SpeedContainer `json:"buff"`
	State            int64          `json:"state"`
	RefreshCount     int64          `json:"refreshCount"`
	LastUpdateTime   int64          `json:"lastUpdateTime"`
	ProcessPoint     float64        `json:"processPoint"`
	Speed            float64        `json:"speed"`
	CompleteWorkTime int64          `json:"completeWorkTime"`
}

type ManufactureInfo struct {
	Buff struct {
		ApCost struct {
			Self interface{} `json:"self"`
		} `json:"apCost"`
		Speed    float64 `json:"speed"`
		Capacity int64   `json:"capacity"`
	} `json:"buff"`
	State             int64   `json:"state"`
	FormulaID         string  `json:"formulaId"`
	RemainSolutionCnt int64   `json:"remainSolutionCnt"`
	OutputSolutionCnt int64   `json:"outputSolutionCnt"`
	LastUpdateTime    int64   `json:"lastUpdateTime"`
	SaveTime          int64   `json:"saveTime"`
	TailTime          float64 `json:"tailTime"`
	ApCost            int64   `json:"apCost"`
	CompleteWorkTime  int64   `json:"completeWorkTime"`
	Capacity          int64   `json:"capacity"`
	ProcessPoint      float64 `json:"processPoint"`
}

type MeetingInfo struct {
	Buff struct {
		Speed  float64 `json:"speed"`
		Weight Weight  `json:"weight"`
	} `json:"buff"`
	State        int64   `json:"state"`
	Speed        float64 `json:"speed"`
	ProcessPoint float64 `json:"processPoint"`
	OwnStock     []Stock `json:"ownStock"`
	ReceiveStock []Stock `json:"receiveStock"`
	Board        struct {
		Rhine      string `json:"RHINE"`
		Blacksteel string `json:"BLACKSTEEL"`
		Ursus      string `json:"URSUS"`
		Glasgow    string `json:"GLASGOW"`
		Kjerag     string `json:"KJERAG"`
		Rhodes     string `json:"RHODES"`
	} `json:"board"`
	SocialReward struct {
		Daily  int64 `json:"daily"`
		Search int64 `json:"search"`
	} `json:"socialReward"`
	DailyReward   interface{} `json:"dailyReward"`
	ExpiredReward int64       `json:"expiredReward"`
	Received      int64       `json:"received"`
	InfoShare     struct {
		Ts     int64 `json:"ts"`
		Reward int64 `json:"reward"`
	} `json:"infoShare"`
	LastUpdateTime   int64 `json:"lastUpdateTime"`
	CompleteWorkTime int64 `json:"completeWorkTime"`
}

type Weight struct {
	Rhine      float64 `json:"RHINE"`
	Penguin    float64 `json:"PENGUIN"`
	Blacksteel float64 `json:"BLACKSTEEL"`
	Ursus      float64 `json:"URSUS"`
	Glasgow    float64 `json:"GLASGOW"`
	Kjerag     float64 `json:"KJERAG"`
	Rhodes     float64 `json:"RHODES"`
}

type Stock struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Number  int64         `json:"number"`
	UID     string        `json:"uid"`
	Name    string        `json:"name"`
	NickNum string        `json:"nickNum"`
	Chars   []CharElement `json:"chars"`
	InUse   int64         `json:"inUse"`
}

type CharElement struct {
	CharID      string `json:"charId"`
	Level       int64  `json:"level"`
	Skin        string `json:"skin"`
	EvolvePhase int64  `json:"evolvePhase"`
}

type PowerInfo struct {
	Buff struct {
		LaborSpeed float64 `json:"laborSpeed"`
	} `json:"buff"`
}

type TradingInfo struct {
	Buff struct {
		Speed  float64 `json:"speed"`
		Limit  int64   `json:"limit"`
		ApCost struct {
			All    int64            `json:"all"`
			Single interface{}      `json:"single"`
			Self   map[string]int64 `json:"self"`
		} `json:"apCost"`
		Rate interface{} `json:"rate"`
	} `json:"buff"`
	State          int64          `json:"state"`
	LastUpdateTime int64          `json:"lastUpdateTime"`
	Strategy       string         `json:"strategy"`
	StockLimit     int64          `json:"stockLimit"`
	ApCost         int64          `json:"apCost"`
	Stock          []StockElement `json:"stock"`
	Next           struct {
		Order        int64   `json:"order"`
		ProcessPoint float64 `json:"processPoint"`
		MaxPoint     int64   `json:"maxPoint"`
		Speed        float64 `json:"speed"`
	} `json:"next"`
	CompleteWorkTime int64 `json:"completeWorkTime"`
}

type StockElement struct {
	InstID   int64  `json:"instId"`
	Delivery []Gain `json:"delivery"`
	Type     string `json:"type"`
	Gain     *Gain  `json:"gain"`
}

type Gain struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type WorkshopInfo struct {
	Buff struct {
		Rate *struct {
			All       int64   `json:"all"`
			WBuilding float64 `json:"W_BUILDING"`
		} `json:"rate"`
		Cost *struct {
			Type      string `json:"type"`
			Limit     int64  `json:"limit"`
			Reduction int64  `json:"reduction"`
		} `json:"cost"`
	} `json:"buff"`
}

type Labor struct {
	BuffSpeed      float64 `json:"buffSpeed"`
	ProcessPoint   float64 `json:"processPoint"`
	Value          int64   `json:"value"`
	LastUpdateTime int64   `json:"lastUpdateTime"`
	MaxValue       int64   `json:"maxValue"`
}

type UserCheckIn struct {
	CanCheckIn         int64   `json:"canCheckIn"`
	CheckInGroupID     string  `json:"checkInGroupId"`
	CheckInRewardIndex int64   `json:"checkInRewardIndex"`
	CheckInHistory     []int64 `json:"checkInHistory"`
}

type ConsumableInfo struct {
	Ts    int64 `json:"ts"`
	Count int64 `json:"count"`
}

type DexNav struct {
	Character map[string]Character        `json:"character"`
	Formula   *Formula                    `json:"formula"`
	Team      map[string]map[string]int64 `json:"team"`
	Enemy     *Enemy                      `json:"enemy"`
}

type Character struct {
	CharInstID int64 `json:"charInstId"`
	Count      int64 `json:"count"`
}

type Enemy struct {
	Enemies map[string]int64    `json:"enemies"`
	Stage   map[string][]string `json:"stage"`
}

type Formula struct {
	Shop        interface{}      `json:"shop"`
	Manufacture map[string]int64 `json:"manufacture"`
	Workshop    map[string]int64 `json:"workshop"`
}

type Dungeon struct {
	Stages    map[string]StageValue `json:"stages"`
	Campaigns *Campaigns            `json:"campaigns"`
}

type Campaigns struct {
	ActiveGroupID      string                  `json:"activeGroupId"`
	CampaignCurrentFee int64                   `json:"campaignCurrentFee"`
	CampaignTotalFee   int64                   `json:"campaignTotalFee"`
	Instances          map[string]InstanceInfo `json:"instances"`
}

type InstanceInfo struct {
	MaxKills     int64   `json:"maxKills"`
	RewardStatus []int64 `json:"rewardStatus"`
}

type StageValue struct {
	StageID         string `json:"stageId"`
	CompleteTimes   int64  `json:"completeTimes"`
	StartTimes      int64  `json:"startTimes"`
	PracticeTimes   int64  `json:"practiceTimes"`
	State           int64  `json:"state"`
	HasBattleReplay int64  `json:"hasBattleReplay"`
	NoCostCnt       int64  `json:"noCostCnt"`
}

type Gacha struct {
	Newbee map[string]int64     `json:"newbee"`
	Normal map[string]GachaInfo `json:"normal"`
}

type GachaInfo struct {
	Cnt    int64 `json:"cnt"`
	MaxCnt int64 `json:"maxCnt"`
	Rarity int64 `json:"rarity"`
	Avail  bool  `json:"avail"`
}

type Mission struct {
	Missions       *Missions        `json:"missions"`
	MissionRewards *MissionRewards  `json:"missionRewards"`
	MissionGroups  map[string]int64 `json:"missionGroups"`
}

type MissionRewards struct {
	DailyPoint  int64    `json:"dailyPoint"`
	WeeklyPoint int64    `json:"weeklyPoint"`
	Rewards     *Rewards `json:"rewards"`
}

type Rewards struct {
	Daily  map[string]int64 `json:"DAILY"`
	Weekly map[string]int64 `json:"WEEKLY"`
}

type Missions struct {
	Openserver map[string]MissionInfo `json:"OPENSERVER"`
	Daily      map[string]MissionInfo `json:"DAILY"`
	Weekly     map[string]MissionInfo `json:"WEEKLY"`
	Guide      map[string]MissionInfo `json:"GUIDE"`
	Main       map[string]MissionInfo `json:"MAIN"`
	Activity   map[string]MissionInfo `json:"ACTIVITY"`
	Sub        map[string]MissionInfo `json:"SUB"`
}

type MissionInfo struct {
	State    int64 `json:"state"`
	Progress []struct {
		Target *int64 `json:"target"`
		Value  int64  `json:"value"`
	} `json:"progress"`
}

type ChainLogin struct {
	IsAvailable bool    `json:"isAvailable"`
	NowIndex    int64   `json:"nowIndex"`
	History     []int64 `json:"history"`
}

type OpenServerCheckIn struct {
	IsAvailable bool    `json:"isAvailable"`
	History     []int64 `json:"history"`
}

type PushFlags struct {
	HasGifts         int64 `json:"hasGifts"`
	HasFriendRequest int64 `json:"hasFriendRequest"`
	HasClues         int64 `json:"hasClues"`
	HasFreeLevelGP   int64 `json:"hasFreeLevelGP"`
	Status           int64 `json:"status"`
}

type Recruit struct {
	Normal *RecruitNormal `json:"normal"`
}

type RecruitNormal struct {
	Slots map[string]RecruitSlot `json:"slots"`
}

type RecruitSlot struct {
	State         int64       `json:"state"`
	Tags          []int64     `json:"tags"`
	SelectTags    []SelectTag `json:"selectTags"`
	StartTs       int64       `json:"startTs"`
	DurationInSEC int64       `json:"durationInSec"`
	MaxFinishTs   int64       `json:"maxFinishTs"`
	RealFinishTs  int64       `json:"realFinishTs"`
}

type SelectTag struct {
	TagID int64 `json:"tagId"`
	Pick  int64 `json:"pick"`
}

type Shop struct {
	LS *struct {
		CurShopID  string `json:"curShopId"`
		CurGroupID string `json:"curGroupId"`
		Info       []Info `json:"info"`
	} `json:"LS"`
	HS *struct {
		CurShopID    string      `json:"curShopId"`
		Info         []Info      `json:"info"`
		ProgressInfo interface{} `json:"progressInfo"`
	} `json:"HS"`
	ES *struct {
		CurShopID string `json:"curShopId"`
		Info      []Info `json:"info"`
	} `json:"ES"`
	Cash InfoSlice `json:"CASH"`
	GP   *struct {
		OneTime InfoSlice `json:"oneTime"`
		Level   InfoSlice `json:"level"`
		Weekly  Monthly   `json:"weekly"`
		Monthly Monthly   `json:"monthly"`
	} `json:"GP"`
	Furni  InfoSlice `json:"FURNI"`
	Social Social    `json:"SOCIAL"`
}

type InfoSlice struct {
	Info []Info `json:"info"`
}

type Info struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
}

type Monthly struct {
	CurGroupID string `json:"curGroupId"`
	Info       []Info `json:"info"`
}

type Social struct {
	CurShopID    string           `json:"curShopId"`
	Info         []Info           `json:"info"`
	CharPurchase map[string]int64 `json:"charPurchase"`
}

type AssistCharList struct {
	CharInstID int64 `json:"charInstId"`
	SkillIndex int64 `json:"skillIndex"`
}

type YesterdayReward struct {
	CanReceive    int64 `json:"canReceive"`
	AssistAmount  int64 `json:"assistAmount"`
	ComfortAmount int64 `json:"comfortAmount"`
	First         int64 `json:"first"`
}

type UserStatus struct {
	NickName                     string           `json:"nickName"`
	NickNumber                   string           `json:"nickNumber"`
	Level                        int64            `json:"level"`
	Exp                          int64            `json:"exp"`
	SocialPoint                  int64            `json:"socialPoint"`
	GachaTicket                  int64            `json:"gachaTicket"`
	TenGachaTicket               int64            `json:"tenGachaTicket"`
	InstantFinishTicket          int64            `json:"instantFinishTicket"`
	HggShard                     int64            `json:"hggShard"`
	LggShard                     int64            `json:"lggShard"`
	RecruitLicense               int64            `json:"recruitLicense"`
	Progress                     int64            `json:"progress"`
	BuyApRemainTimes             int64            `json:"buyApRemainTimes"`
	ApLimitUpFlag                int64            `json:"apLimitUpFlag"`
	UID                          string           `json:"uid"`
	Flags                        map[string]int64 `json:"flags"`
	Ap                           int64            `json:"ap"`
	MaxAp                        int64            `json:"maxAp"`
	PayDiamond                   int64            `json:"payDiamond"`
	FreeDiamond                  int64            `json:"freeDiamond"`
	DiamondShard                 int64            `json:"diamondShard"`
	Gold                         int64            `json:"gold"`
	PracticeTicket               int64            `json:"practiceTicket"`
	LastRefreshTs                int64            `json:"lastRefreshTs"`
	LastApAddTime                int64            `json:"lastApAddTime"`
	MainStageProgress            string           `json:"mainStageProgress"`
	RegisterTs                   int64            `json:"registerTs"`
	LastOnlineTs                 int64            `json:"lastOnlineTs"`
	ServerName                   string           `json:"serverName"`
	AvatarID                     string           `json:"avatarId"`
	Resume                       string           `json:"resume"`
	FriendNumLimit               int64            `json:"friendNumLimit"`
	MonthlySubscriptionStartTime int64            `json:"monthlySubscriptionStartTime"`
	MonthlySubscriptionEndTime   int64            `json:"monthlySubscriptionEndTime"`
	Secretary                    string           `json:"secretary"`
	SecretarySkinID              string           `json:"secretarySkinId"`
}

type Troop struct {
	CurCharInstID int64                `json:"curCharInstId"`
	CurSquadCount int64                `json:"curSquadCount"`
	Squads        map[string]Squad     `json:"squads"`
	Chars         map[string]TroopChar `json:"chars"`
}

type TroopChar struct {
	InstID            int64   `json:"instId"`
	CharID            string  `json:"charId"`
	FavorPoint        int64   `json:"favorPoint"`
	PotentialRank     int64   `json:"potentialRank"`
	MainSkillLvl      int64   `json:"mainSkillLvl"`
	Skin              string  `json:"skin"`
	Level             int64   `json:"level"`
	Exp               int64   `json:"exp"`
	EvolvePhase       int64   `json:"evolvePhase"`
	DefaultSkillIndex int64   `json:"defaultSkillIndex"`
	GainTime          int64   `json:"gainTime"`
	Skills            []Skill `json:"skills"`
}

type Skill struct {
	SkillID             string `json:"skillId"`
	Unlock              int64  `json:"unlock"`
	State               int64  `json:"state"`
	SpecializeLevel     int64  `json:"specializeLevel"`
	CompleteUpgradeTime int64  `json:"completeUpgradeTime"`
}

type Squad struct {
	SquadID string           `json:"squadId"`
	Name    string           `json:"name"`
	Slots   []AssistCharList `json:"slots"`
}
