package basis

type Brightness int

const (
	BrightnessMing Brightness = iota // 廟 (M)
	BrightnessWang                   // 旺 (W)
	BrightnessLi                     // 利 (L)
	BrightnessPing                   // 平 (P)
	BrightnessXian                   // 陷 (X)
	BrightnessBu                     // 不 (B)
)

func (b Brightness) String() string {
	names := []string{"廟", "旺", "利", "平", "陷", "不"}
	return names[b]
}

// Global Aliases for Table
const (
	M = BrightnessMing
	W = BrightnessWang
	L = BrightnessLi
	P = BrightnessPing
	X = BrightnessXian
	B = BrightnessBu
)

var StarBrightnessTable = map[Star][]Brightness{
	StarZiwei:     {W, M, L, P, M, L, W, M, L, P, M, L},
	StarTianji:    {M, X, L, M, L, P, X, L, P, M, L, P},
	StarTaiyang:   {X, X, L, W, M, M, M, L, P, X, X, X},
	StarWuqu:      {W, M, L, L, M, P, P, M, L, W, M, L},
	StarTiantong:  {M, L, P, L, X, M, X, P, W, L, X, L},
	StarLianzhen:  {P, L, M, P, L, X, L, M, W, P, L, X},
	StarTianfu:    {M, W, M, L, W, M, W, M, M, L, W, M},
	StarTaiyin:    {M, M, L, L, X, X, X, X, P, W, M, M},
	StarTanlang:   {W, P, M, M, W, L, P, L, L, P, M, L},
	StarJumen:     {W, M, M, L, L, P, W, X, L, P, P, X},
	StarTianxiang: {M, M, L, P, P, L, M, W, L, P, P, L},
	StarTianliang: {M, M, L, P, L, M, W, X, L, P, P, X},
	StarQisha:     {M, W, M, P, L, L, M, W, M, P, L, L},
	StarPojun:     {W, M, L, X, W, P, L, M, L, X, W, P},
}

func BrightnessLevel(s Star, b Branch) Brightness {
	if table, ok := StarBrightnessTable[s]; ok {
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
