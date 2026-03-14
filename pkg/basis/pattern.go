package basis

type Pattern struct {
	Name        string
	Description string
	Level       string
}

var Patterns = []Pattern{
	{Name: "紫府同宮", Description: "紫微天府同宮，尊貴無極", Level: "甲"},
	{Name: "紫破同宮", Description: "紫微破軍同宮，權力慾望極強", Level: "甲"},
	{Name: "紫相拱照", Description: "紫微天相拱照，有權有勢", Level: "甲"},
	{Name: "殺破狼格", Description: "七殺破軍貪狼同宮或對宮，三方四正見殺破狼", Level: "甲"},
	{Name: "機月同梁格", Description: "天機太陰天同天梁在三方四正", Level: "甲"},
	{Name: "紫武廉府", Description: "紫微武曲廉貞天府四星會合", Level: "甲"},
	{Name: "府相朝垣", Description: "天府天相會合於三方四正", Level: "乙"},
	{Name: "日月拱照", Description: "太陽太陰在三方四正拱照命宮", Level: "甲"},
	{Name: "日月反背", Description: "太陽太陰在遷移宮太陽在陷地", Level: "辛"},
	{Name: "祿馬交馳", Description: "祿存與天馬同宮或對宮", Level: "乙"},
	{Name: "天馬拱命", Description: "天馬在命宮三方四正", Level: "乙"},
	{Name: "火貪格", Description: "火星或鈴星與貪狼同宮", Level: "乙"},
	{Name: "鈴貪格", Description: "鈴星與貪狼同宮", Level: "乙"},
	{Name: "擎羊入命", Description: "擎羊星在命宮", Level: "辛"},
	{Name: "陀羅入命", Description: "陀羅星在命宮", Level: "辛"},
	{Name: "空宮", Description: "命宮無主星", Level: "丁"},
	{Name: "桃花犯主", Description: "紅鸞天姚與主星同宮", Level: "丙"},
	{Name: "水木清華", Description: "天機與太陰同宮或對宮", Level: "乙"},
	{Name: "土金相生", Description: "天府與武曲同宮或相生", Level: "丙"},
	{Name: "日月並明", Description: "太陽在午宮，太陰在丑宮或未宮", Level: "乙"},
}
