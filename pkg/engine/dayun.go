package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func CalcDaYun(mingBranch basis.Branch, yearStem basis.Stem, sex basis.Sex, wuxing basis.Wuxing) []basis.DaYun {
	isYang := (int(yearStem) % 2) == 0
	
	// direction = 1 (陽男陰女): 順行 (Clockwise through Branches)
	// direction = -1 (陰男陽女): 逆行 (Counter-clockwise through Branches)
	direction := 1
	if (sex == basis.SexMale && !isYang) || (sex == basis.SexFemale && isYang) {
		direction = -1
	}

	startAge := wuxing.Value()
	dayuns := make([]basis.DaYun, 0)
	
	for i := 0; i < 12; i++ {
		// Calculate the physical branch for this DaYun
		// DaYun starts from Ming Palace (Index 0 in logic, but physical branch is mingBranch)
		// If direction=1, move to next branch (mingBranch + 1, + 2...)
		branchIdx := (int(mingBranch) + (direction * i) + 12) % 12
		
		ageStart := startAge + i*10
		ageEnd := ageStart + 9
		
		dayuns = append(dayuns, basis.DaYun{
			Index:    i + 1,
			StartAge: ageStart,
			EndAge:   ageEnd,
			Branch:   basis.Branch(branchIdx),
		})
	}
	return dayuns
}

// CalcLiuNian calculates the yearly luck for a target year
func CalcLiuNian(yearBranch basis.Branch, currentYear int) basis.LiuNian {
	// Simple rule: Liu Nian Palace is at the year's branch.
	// 2024 (Jia-Chen) -> Chen(4)
	// 2025 (Yi-Si) -> Si(5)
	// Standard formula: (Year - 4) % 12 gives branch index (0=Zi is at offset 4? No)
	// Let's use the actual yearBranch provided.
	return basis.LiuNian{
		Year:   currentYear,
		Branch: yearBranch,
	}
}

// CalcLiuYue calculates the monthly luck based on "Dou Jun" method
func CalcLiuYue(lnBranch basis.Branch, birthMonth int, birthHour basis.Branch, targetLunarMonth int) basis.Branch {
	// Dou Jun Rule: Starting from Liu Nian Branch, count counter-clockwise to Birth Month, 
	// then clockwise to Birth Hour to find the starting point (January).
	// month1Idx = (lnZhiIdx - bMonthIdx + bHourIdx + 12) % 12
	month1Idx := (int(lnBranch) - (birthMonth - 1) + int(birthHour) + 12) % 12
	
	// Each lunar month moves clockwise from month1Idx
	lYueIdx := (month1Idx + (targetLunarMonth - 1)) % 12
	return basis.Branch(lYueIdx)
}

// CalcLiuRi calculates daily luck
func CalcLiuRi(lyBranch basis.Branch, targetLunarDay int) basis.Branch {
	// Start from Liu Yue Branch, move clockwise by (targetLunarDay - 1)
	lRiIdx := (int(lyBranch) + (targetLunarDay - 1)) % 12
	return basis.Branch(lRiIdx)
}
