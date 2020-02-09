package itemtable

import "encoding/json"

func Unmarshal(data []byte) (ItemTable, error) {
	var r ItemTable
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ItemTable) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ItemTable struct {
	Items          map[string]Item          `json:"items"`
	ExpItems       map[string]ExpItem       `json:"expItems"`
	PotentialItems map[string]PotentialItem `json:"potentialItems"`
	ApSupplies     ApSupplies               `json:"apSupplies"`
}

type ApSupplies struct {
	ApSupplyLt60  ApSupplyLt `json:"ap_supply_lt_60"`
	ApSupplyLt100 ApSupplyLt `json:"ap_supply_lt_100"`
}

type ApSupplyLt struct {
	ID    string `json:"id"`
	Ap    int64  `json:"ap"`
	HasTs bool   `json:"hasTs"`
}

type ExpItem struct {
	ID      string `json:"id"`
	GainExp int64  `json:"gainExp"`
}

type Item struct {
	ItemID              string                `json:"itemId"`
	Name                string                `json:"name"`
	Description         string                `json:"description"`
	Rarity              int64                 `json:"rarity"`
	IconID              string                `json:"iconId"`
	StackIconID         *string               `json:"stackIconId"`
	SortID              int64                 `json:"sortId"`
	Usage               string                `json:"usage"`
	ObtainApproach      *string               `json:"obtainApproach"`
	ClassifyType        ClassifyType          `json:"classifyType"`
	ItemType            string                `json:"itemType"`
	StageDropList       []StageDropList       `json:"stageDropList"`
	BuildingProductList []BuildingProductList `json:"buildingProductList"`
}

type BuildingProductList struct {
	RoomType  RoomType `json:"roomType"`
	FormulaID string   `json:"formulaId"`
}

type StageDropList struct {
	StageID    string     `json:"stageId"`
	Occurrance Occurrance `json:"occPer"`
}

type PotentialItem struct {
	Pioneer string `json:"PIONEER"`
	Warrior string `json:"WARRIOR"`
	Sniper  string `json:"SNIPER"`
	Tank    string `json:"TANK"`
	Medic   string `json:"MEDIC"`
	Support string `json:"SUPPORT"`
	Caster  string `json:"CASTER"`
	Special string `json:"SPECIAL"`
}

type RoomType string

const (
	Manufacture RoomType = "MANUFACTURE"
	Workshop    RoomType = "WORKSHOP"
)

type ClassifyType string

const (
	Consume  ClassifyType = "CONSUME"
	Material ClassifyType = "MATERIAL"
	None     ClassifyType = "NONE"
	Normal   ClassifyType = "NORMAL"
)

type Occurrance string

const (
	Almost    Occurrance = "ALMOST"
	Always    Occurrance = "ALWAYS"
	Often     Occurrance = "OFTEN"
	Sometimes Occurrance = "SOMETIMES"
	Usual     Occurrance = "USUAL"
)
