package basis

type TransformationType int

const (
	TransformationLu   TransformationType = iota // 化祿
	TransformationQuan                           // 化權
	TransformationKe                             // 化科
	TransformationJi                             // 化忌
)

const (
	TransLu   = TransformationLu
	TransQuan = TransformationQuan
	TransKe   = TransformationKe
	TransJi   = TransformationJi
)

func (t TransformationType) String() string {
	names := []string{"化祿", "化權", "化科", "化忌"}
	return names[t]
}

type TransformedStar struct {
	StarName       string
	Transformation TransformationType
}

func (t TransformedStar) String() string {
	return t.StarName + t.Transformation.String()
}

// TransformationTable maps Stem to [Lu, Quan, Ke, Ji] star names
var TransformationTable = map[Stem][4]string{
	StemJia:  {"廉貞", "破軍", "武曲", "太陽"}, // 廉破武陽
	StemYi:   {"天機", "天梁", "紫微", "太陰"}, // 機梁紫陰
	StemBing: {"天同", "天機", "文昌", "廉貞"}, // 同機昌廉
	StemDing: {"太陰", "天同", "天機", "巨門"}, // 陰同機巨
	StemWu:   {"貪狼", "太陰", "右弼", "天機"}, // 貪陰右機
	StemJi:   {"武曲", "貪狼", "天梁", "文曲"}, // 武貪梁曲
	StemGeng: {"太陽", "武曲", "天府", "天同"}, // 陽武府同 (or 陽武陰同, using 府 for classical alignment in some schools, let's use 陽武府同 as requested/typical for some)
	StemXin:  {"巨門", "太陽", "文曲", "文昌"}, // 巨陽曲昌
	StemRen:  {"天梁", "紫微", "左輔", "武曲"}, // 梁紫左武
	StemGui:  {"破軍", "巨門", "太陰", "貪狼"}, // 破巨陰貪
}
