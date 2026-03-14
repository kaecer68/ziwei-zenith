package basis

type Wuxing int

const (
	WuxingShuiEr Wuxing = iota // 水二局
	WuxingLiu                  // 火六局
	WuxingMuSan                // 木三局
	WuxingTuWu                 // 土五局
	WuxingJinSi                // 金四局
)

func (w Wuxing) String() string {
	names := []string{"水二局", "火六局", "木三局", "土五局", "金四局"}
	return names[w]
}

func (w Wuxing) Value() int {
	values := []int{2, 6, 3, 5, 4}
	return values[w]
}

func (w Wuxing) Element() Element {
	elements := []Element{
		ElementWater,
		ElementFire,
		ElementWood,
		ElementEarth,
		ElementMetal,
	}
	return elements[w]
}

type NaYin int

const (
	NaYinJiaZi    NaYin = iota // 甲子
	NaYinYiChou                // 乙丑
	NaYinBingYin               // 丙寅
	NaYinDingMao               // 丁卯
	NaYinWuChen                // 戊辰
	NaYinJiSi                  // 己巳
	NaYinGengWu                // 庚午
	NaYinXinWei                // 辛未
	NaYinRenShen               // 壬申
	NaYinGuiYou                // 癸酉
	NaYinJiaXu                 // 甲戌
	NaYinYiHai                 // 乙亥
	NaYinBingZi                // 丙子
	NaYinDingChou              // 丁丑
	NaYinWuYin                 // 戊寅
	NaYinJiMao                 // 己卯
	NaYinGengChen              // 庚辰
	NaYinXinSi                 // 辛巳
	NaYinRenWu                 // 壬午
	NaYinGuiWei                // 癸未
	NaYinJiaShen               // 甲申
	NaYinYiYou                 // 乙酉
	NaYinBingXu                // 丙戌
	NaYinDingHai               // 丁亥
	NaYinWuZi                  // 戊子
	NaYinJiChou                // 己丑
	NaYinGengYin               // 庚寅
	NaYinXinMao                // 辛卯
	NaYinRenChen               // 壬辰
	NaYinGuiSi                 // 癸巳
	NaYinJiaWu                 // 甲午
	NaYinYiWei                 // 乙未
	NaYinBingShen              // 丙申
	NaYinDingYou               // 丁酉
	NaYinWuXu                  // 戊戌
	NaYinJiHai                 // 己亥
	NaYinGengZi                // 庚子
	NaYinXinChou               // 辛丑
	NaYinRenYin                // 壬寅
	NaYinGuiMao                // 癸卯
	NaYinJiaChen               // 甲辰
	NaYinYiSi                  // 乙巳
	NaYinBingWu                // 丙午
	NaYinDingWei               // 丁未
	NaYinWuShen                // 戊申
	NaYinJiYou                 // 己酉
	NaYinGengXu                // 庚戌
	NaYinXinHai                // 辛亥
	NaYinRenZi                 // 壬子
	NaYinGuiChou               // 癸丑
)

func (n NaYin) Element() Element {
	elements := []Element{
		ElementWood, ElementWood,
		ElementFire, ElementFire,
		ElementEarth, ElementEarth,
		ElementMetal, ElementMetal,
		ElementWater, ElementWater,
		ElementWood, ElementWood,
		ElementFire, ElementFire,
		ElementEarth, ElementEarth,
		ElementMetal, ElementMetal,
		ElementWater, ElementWater,
		ElementWood, ElementWood,
		ElementFire, ElementFire,
		ElementEarth, ElementEarth,
		ElementMetal, ElementMetal,
		ElementWater, ElementWater,
		ElementWood, ElementWood,
		ElementFire, ElementFire,
		ElementEarth, ElementEarth,
		ElementMetal, ElementMetal,
		ElementWater, ElementWater,
		ElementWood, ElementWood,
		ElementFire, ElementFire,
		ElementEarth, ElementEarth,
		ElementMetal, ElementMetal,
		ElementWater, ElementWater,
		ElementWood, ElementWood,
		ElementFire, ElementFire,
		ElementEarth, ElementEarth,
		ElementMetal, ElementMetal,
		ElementWater, ElementWater,
	}
	return elements[n]
}

func CalcNaYin(stem Stem, branch Branch) NaYin {
	return NaYin(int(stem)*12 + int(branch))
}

func (n NaYin) String() string {
	names := []string{
		"甲子", "乙丑", "丙寅", "丁卯", "戊辰", "己巳", "庚午", "辛未", "壬申", "癸酉",
		"甲戌", "乙亥", "丙子", "丁丑", "戊寅", "己卯", "庚辰", "辛巳", "壬午", "癸未",
		"甲申", "乙酉", "丙戌", "丁亥", "戊子", "己丑", "庚寅", "辛卯", "壬辰", "癸巳",
		"甲午", "乙未", "丙申", "丁酉", "戊戌", "己亥", "庚子", "辛丑", "壬寅", "癸卯",
		"甲辰", "乙巳", "丙午", "丁未", "戊申", "己酉", "庚戌", "辛亥", "壬子", "癸丑",
	}
	return names[n]
}
