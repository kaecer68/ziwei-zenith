package basis

type Brightness int

const (
	BrightnessMing Brightness = iota // 廟 (M)
	BrightnessWang                   // 旺 (W)
	BrightnessDe                     // 得 (D)
	BrightnessLi                     // 利 (L)
	BrightnessPing                   // 平 (P)
	BrightnessXian                   // 陷 (X)
	BrightnessBu                     // 不 (B)
)

func (b Brightness) String() string {
	names := []string{"廟", "旺", "得", "利", "平", "陷", "不"}
	return names[b]
}

// Global Aliases for Table
const (
	M = BrightnessMing
	W = BrightnessWang
	D = BrightnessDe
	L = BrightnessLi
	P = BrightnessPing
	X = BrightnessXian
	B = BrightnessBu
)

// 子  丑  寅  卯  辰  巳  午  未  申  酉  戌  亥
var StarBrightnessTable = map[Star][]Brightness{
	StarZiwei:     {W, W, D, W, M, D, W, M, D, P, M, D},
	StarTianji:    {M, X, D, M, D, P, X, D, P, M, L, P},
	StarTaiyang:   {X, X, D, W, M, M, M, D, P, X, X, X},
	StarWuqu:      {W, M, D, L, M, P, P, M, D, W, M, P},
	StarTiantong:  {M, L, P, D, X, M, X, P, W, D, P, L},
	StarLianzhen:  {P, L, M, P, D, X, D, L, W, P, L, X},
	StarTianfu:    {M, M, M, D, W, M, W, M, M, D, W, M},
	StarTaiyin:    {M, M, W, L, X, X, X, X, P, W, M, M},
	StarTanlang:   {W, P, M, L, W, D, P, L, D, P, M, L},
	StarJumen:     {W, M, M, L, X, P, W, X, L, P, P, X},
	StarTianxiang: {M, M, D, P, P, D, M, W, D, P, P, D},
	StarTianliang: {M, M, D, P, D, M, M, X, D, P, P, X},
	StarQisha:     {M, W, M, P, D, D, M, M, M, P, D, D},
	StarPojun:     {W, M, D, X, W, P, L, M, D, X, W, P},
}

func BrightnessLevel(s Star, b Branch) Brightness {
	if table, ok := StarBrightnessTable[s]; ok {
		return table[int(b)]
	}
	return BrightnessWang
}

// 六吉星亮度表：左輔、右弼、文昌、文曲、天魁、天鉞
// 子  丑  寅  卯  辰  巳  午  未  申  酉  戌  亥
var AuspiciousStarBrightnessTable = map[AuspiciousStar][]Brightness{
	AuspiciousZuofu:    {W, M, D, W, M, D, W, M, D, W, M, D}, // 左輔
	AuspiciousYoubi:    {D, W, M, D, W, M, D, W, M, D, W, M}, // 右弼
	AuspiciousWenchang: {X, X, W, M, D, D, W, M, D, D, W, X}, // 文昌
	AuspiciousWenqu:    {W, M, D, D, W, X, X, W, M, D, D, W}, // 文曲
	AuspiciousTiankui:  {M, M, M, M, M, M, M, M, M, M, M, M}, // 天魁(固定旺)
	AuspiciousTianyue:  {M, M, M, M, M, M, M, M, M, M, M, M}, // 天鉞(固定旺)
}

func AuspiciousBrightnessLevel(s AuspiciousStar, b Branch) Brightness {
	if table, ok := AuspiciousStarBrightnessTable[s]; ok {
		return table[int(b)]
	}
	return BrightnessWang
}

// 祿存、天馬亮度表
// 子  丑  寅  卯  辰  巳  午  未  申  酉  戌  亥
var LuCunStarBrightnessTable = map[LuCunStar][]Brightness{
	LuCun:  {M, M, W, D, D, W, M, M, W, D, D, W}, // 祿存
	Tianma: {W, D, M, M, D, W, M, M, D, W, M, D}, // 天馬
}

func LuCunBrightnessLevel(s LuCunStar, b Branch) Brightness {
	if table, ok := LuCunStarBrightnessTable[s]; ok {
		return table[int(b)]
	}
	return BrightnessWang
}

// 六煞星亮度表：擎羊、陀羅、火星、鈴星、地空、地劫
// 子  丑  寅  卯  辰  巳  午  未  申  酉  戌  亥
var MaleficStarBrightnessTable = map[MaleficStar][]Brightness{
	MaleficQingyang: {X, M, W, D, D, W, M, X, W, D, D, W}, // 擎羊
	MaleficTuoluo:   {W, D, D, W, M, X, W, D, D, W, M, X}, // 陀羅
	MaleficHuoxing:  {W, D, M, X, W, D, M, D, W, X, W, D}, // 火星
	MaleficLingxing: {D, W, X, W, D, M, D, W, X, W, D, M}, // 鈴星
	MaleficDikong:   {P, P, P, P, P, P, P, P, P, P, P, P}, // 地空(固定平)
	MaleficDijie:    {P, P, P, P, P, P, P, P, P, P, P, P}, // 地劫(固定平)
}

func MaleficBrightnessLevel(s MaleficStar, b Branch) Brightness {
	if table, ok := MaleficStarBrightnessTable[s]; ok {
		return table[int(b)]
	}
	return BrightnessWang
}

type StarBrightness struct {
	Star       Star
	Branch     Branch
	Palace     Palace
	Brightness Brightness
}

func (s StarBrightness) String() string {
	return s.Star.String() + "(" + s.Brightness.String() + ")"
}
