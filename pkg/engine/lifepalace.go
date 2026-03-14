package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type LifePalace struct {
	MingGong basis.Palace
	ShenGong basis.Palace
}

func CalcLifePalace(lunarMonth int, lunarDay int, hourBranch basis.Branch) LifePalace {
	monthBranch := basis.Branch(0)
	switch lunarMonth {
	case 1:
		monthBranch = basis.BranchYin
	case 2:
		monthBranch = basis.BranchMao
	case 3:
		monthBranch = basis.BranchChen
	case 4:
		monthBranch = basis.BranchSi
	case 5:
		monthBranch = basis.BranchWu
	case 6:
		monthBranch = basis.BranchWei
	case 7:
		monthBranch = basis.BranchShen
	case 8:
		monthBranch = basis.BranchYou
	case 9:
		monthBranch = basis.BranchXu
	case 10:
		monthBranch = basis.BranchHai
	case 11:
		monthBranch = basis.BranchZi
	case 12:
		monthBranch = basis.BranchChou
	}

	monthOffset := int(monthBranch)

	mingOffset := (monthOffset - int(hourBranch) + 12) % 12
	shenOffset := (monthOffset + int(hourBranch)) % 12

	return LifePalace{
		MingGong: basis.Palace(mingOffset),
		ShenGong: basis.Palace(shenOffset),
	}
}

func BuildPalaces(mingGong basis.Palace) map[basis.Palace]basis.Branch {
	palaces := make(map[basis.Palace]basis.Branch)
	mingIdx := int(mingGong)
	mingBranch := basis.Branch(mingIdx)

	for i := 0; i < 12; i++ {
		palace := basis.Palace(i)
		branchIdx := (mingIdx + i) % 12
		palaces[palace] = basis.Branch(branchIdx)
	}

	_ = mingBranch

	return palaces
}
