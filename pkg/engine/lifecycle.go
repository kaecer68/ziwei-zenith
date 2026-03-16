package engine

import "github.com/kaecer68/ziwei-zenith/pkg/basis"

var changShengStartTable = map[basis.Wuxing]basis.Branch{
	basis.WuxingShui2: basis.BranchShen,
	basis.WuxingMu3:   basis.BranchHai,
	basis.WuxingJin4:  basis.BranchSi,
	basis.WuxingTu5:   basis.BranchShen,
	basis.WuxingHuo6:  basis.BranchYin,
}

func CalcChangShengStages(wuxing basis.Wuxing, sex basis.Sex, yearStem basis.Stem) map[basis.Branch]basis.ChangShengStage {
	result := make(map[basis.Branch]basis.ChangShengStage, 12)
	start, ok := changShengStartTable[wuxing]
	if !ok {
		start = basis.BranchYin
	}

	isYang := int(yearStem)%2 == 0
	forward := (sex == basis.SexMale && isYang) || (sex == basis.SexFemale && !isYang)

	for i := 0; i < 12; i++ {
		branchIdx := (int(start) + i) % 12
		stageIdx := i
		if !forward {
			branchIdx = (int(start) - i + 12) % 12
		}
		result[basis.Branch(branchIdx)] = basis.ChangShengStage(stageIdx)
	}

	return result
}

func CalcBoShiStars(luCunBranch basis.Branch, sex basis.Sex, yearStem basis.Stem) map[basis.Branch]basis.BoShiStar {
	result := make(map[basis.Branch]basis.BoShiStar, 12)
	isYang := int(yearStem)%2 == 0
	forward := (sex == basis.SexMale && isYang) || (sex == basis.SexFemale && !isYang)

	for i := 0; i < 12; i++ {
		branchIdx := (int(luCunBranch) + i) % 12
		if !forward {
			branchIdx = (int(luCunBranch) - i + 12) % 12
		}
		result[basis.Branch(branchIdx)] = basis.BoShiStar(i)
	}

	return result
}

func locateLuCun(yearStem basis.Stem) basis.Branch {
	luCunIdx := map[basis.Stem]int{
		basis.StemJia: 2,
		basis.StemYi: 3,
		basis.StemBing: 5,
		basis.StemDing: 6,
		basis.StemWu: 5,
		basis.StemJi: 6,
		basis.StemGeng: 8,
		basis.StemXin: 9,
		basis.StemRen: 11,
		basis.StemGui: 0,
	}[yearStem]
	return basis.Branch(luCunIdx)
}
