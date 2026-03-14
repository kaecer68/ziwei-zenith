package basis

type Star int

const (
	StarZiwei     Star = iota // 紫微星
	StarTianfu                // 天府星
	StarTianji                // 天機星
	StarTaiyang               // 太陽星
	StarWuqu                  // 武曲星
	StarTiantong              // 天同星
	StarLianzhen              // 廉貞星
	StarTanlang               // 貪狼星
	StarJumen                 // 巨門星
	StarTianxiang             // 天相星
	StarTianliang             // 天梁星
	StarQisha                 // 七殺星
	StarPojun                 // 破軍星
	StarTaiyin                // 太陰星
)

func (s Star) String() string {
	names := []string{
		"紫微", "天府", "天機", "太陽", "武曲",
		"天同", "廉貞", "貪狼", "巨門", "天相",
		"天梁", "七殺", "破軍", "太陰",
	}
	return names[s]
}

func (s Star) Element() Element {
	elements := []Element{
		ElementEarth, // 紫微 - 己土
		ElementEarth, // 天府 - 戊土
		ElementWood,  // 天機 - 乙木
		ElementFire,  // 太陽 - 丙火
		ElementMetal, // 武曲 - 辛金
		ElementWater, // 天同 - 壬水
		ElementFire,  // 廉貞 - 丁火
		ElementWood,  // 貪狼 - 甲木/癸水 (氣屬甲木，體屬癸水)
		ElementWater, // 巨門 - 癸水
		ElementWater, // 天相 - 壬水
		ElementEarth, // 天梁 - 戊土
		ElementMetal, // 七殺 - 庚金
		ElementWater, // 破軍 - 癸水
		ElementWater, // 太陰 - 癸水
	}
	return elements[s]
}

func (s Star) Category() string {
	categories := []string{
		"北斗", "南斗", "南斗", "中天", "北斗",
		"南斗", "北斗", "北斗", "北斗", "南斗",
		"南斗", "北斗", "北斗", "中天",
	}
	return categories[s]
}

type StarPosition struct {
	Star      Star
	Palace    Palace
	IsBright  bool
	Intensity int
}

type StarGroup struct {
	Stars []Star
}
