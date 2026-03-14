package basis

import "fmt"

type Wuxing int

const (
	WuxingShui2 Wuxing = 2 // 水二局
	WuxingMu3   Wuxing = 3 // 木三局
	WuxingJin4  Wuxing = 4 // 金四局
	WuxingTu5   Wuxing = 5 // 土五局
	WuxingHuo6  Wuxing = 6 // 火六局
)

func (w Wuxing) Value() int {
	return int(w)
}

func (w Wuxing) String() string {
	names := map[Wuxing]string{
		WuxingShui2: "水二局",
		WuxingMu3:   "木三局",
		WuxingJin4:  "金四局",
		WuxingTu5:   "土五局",
		WuxingHuo6:  "火六局",
	}
	return names[w]
}

// WuxingJuTable maps Palace GanZhi (e.g., "甲子") to Ju number.
// Based on the verified logic in master_engine_v4.js
var WuxingJuTable = map[string]int{
	"甲子": 4, "乙丑": 4, "丙寅": 6, "丁卯": 6, "戊辰": 3, "己巳": 3, "庚午": 5, "辛未": 5, "壬申": 4, "癸酉": 4, "甲戌": 6, "乙亥": 6,
	"丙子": 2, "丁丑": 2, "戊寅": 5, "己卯": 5, "庚辰": 4, "辛巳": 4, "壬午": 3, "癸未": 3, "甲申": 2, "乙酉": 2, "丙戌": 5, "丁亥": 5,
	"戊子": 6, "己丑": 6, "庚寅": 3, "辛卯": 3, "壬辰": 2, "癸巳": 2, "甲午": 4, "乙未": 4, "丙申": 6, "丁酉": 6, "戊戌": 3, "己亥": 3,
	"庚子": 5, "辛丑": 5, "壬寅": 4, "癸卯": 4, "甲辰": 6, "乙巳": 6, "丙午": 2, "丁未": 2, "戊申": 5, "己酉": 5, "庚戌": 4, "辛亥": 4,
	"壬子": 3, "癸丑": 3, "甲寅": 2, "乙卯": 2, "丙辰": 5, "丁巳": 5, "戊午": 6, "己未": 6, "庚申": 3, "辛酉": 3, "壬戌": 2, "癸亥": 2,
}

func GetWuxingJu(gan Stem, zhi Branch) Wuxing {
	key := fmt.Sprintf("%s%s", gan, zhi)
	if val, ok := WuxingJuTable[key]; ok {
		return Wuxing(val)
	}
	return WuxingShui2 // Default
}
