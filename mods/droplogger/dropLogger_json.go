package droplogger

import "encoding/json"

func unmarshalBattleFinish(data []byte) (battleFinish, error) {
	var r battleFinish
	err := json.Unmarshal(data, &r)
	return r, err
}

type battleFinish struct {
	Result            int64         `json:"result"`
	ExpScale          float64       `json:"expScale"`
	GoldScale         float64       `json:"goldScale"`
	Rewards           []reward      `json:"rewards"`
	FirstRewards      []reward      `json:"firstRewards"`
	UnlockStages      []string      `json:"unlockStages"`
	UnusualRewards    []reward      `json:"unusualRewards"`
	AdditionalRewards []reward      `json:"additionalRewards"`
	FurnitureRewards  []reward      `json:"furnitureRewards"`
	Alert             []interface{} `json:"alert"`
}

type reward struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
	Type  string `json:"type"`
}
