package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type LifePalace struct {
	MingGong basis.Branch
	ShenGong basis.Branch
}

func CalcLifePalace(lunarMonth int, hourBranch basis.Branch) LifePalace {
	// Month start from 寅 (2)
	monthPos := (2 + lunarMonth - 1) % 12

	// MingGong moves counter-clockwise from monthPos by hourBranch
	mingIdx := (monthPos - int(hourBranch) + 12) % 12

	// ShenGong moves clockwise from monthPos by hourBranch
	shenIdx := (monthPos + int(hourBranch)) % 12

	return LifePalace{
		MingGong: basis.Branch(mingIdx),
		ShenGong: basis.Branch(shenIdx),
	}
}

func BuildPalaces(mingBranch basis.Branch) map[basis.Palace]basis.Branch {
	palaceMap := make(map[basis.Palace]basis.Branch)
	mingIdx := int(mingBranch)

	for i := 0; i < 12; i++ {
		palaceType := basis.Palace(i)
		// Palaces flow counter-clockwise: 0:命, 1:兄, 2:妻...
		// Branch index decreases as Palace index increases
		branchIdx := (mingIdx - i + 12) % 12
		palaceMap[palaceType] = basis.Branch(branchIdx)
	}

	return palaceMap
}
