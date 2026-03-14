package basis

type TransformationType int

const (
	TransformationLu   TransformationType = iota // 化祿
	TransformationQuan                           // 化權
	TransformationKe                             // 化科
	TransformationJi                             // 化忌
)

func (t TransformationType) String() string {
	names := []string{"化祿", "化權", "化科", "化忌"}
	return names[t]
}

type TransformedStar struct {
	Star           Star
	Transformation TransformationType
}

func (t TransformedStar) String() string {
	return t.Star.String() + t.Transformation.String()
}

var TransformationTable = map[Stem][4]Branch{
	StemJia:  {BranchHai, BranchSi, BranchShen, BranchYin},
	StemYi:   {BranchWu, BranchChen, BranchMao, BranchShen},
	StemBing: {BranchYin, BranchShen, BranchHai, BranchSi},
	StemDing: {BranchMao, BranchYou, BranchZi, BranchWu},
	StemWu:   {BranchSi, BranchHai, BranchShen, BranchYin},
	StemJi:   {BranchWu, BranchChen, BranchMao, BranchShen},
	StemGeng: {BranchShen, BranchYin, BranchHai, BranchChen},
	StemXin:  {BranchYou, BranchMao, BranchYin, BranchZi},
	StemRen:  {BranchHai, BranchSi, BranchShen, BranchChen},
	StemGui:  {BranchZi, BranchWu, BranchMao, BranchChen},
}
