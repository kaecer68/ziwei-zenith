package engine

import (
	"testing"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func TestCalcDaYun(t *testing.T) {
	// Yang Male (Ren Year Stem), Water 2 Ju, Ming at Chen (4)
	// Ren (9) is Odd? StemRen is index 8 (Jia=0, Yi=1... Ren=8).
	// 8 % 2 == 0 -> isYang = true.
	// Yang Male -> Direction = 1 (Clockwise)
	mingBranch := basis.Branch(4)
	yearStem := basis.StemRen
	sex := basis.SexMale
	wuxing := basis.WuxingShui2 // 2

	dayuns := CalcDaYun(mingBranch, yearStem, sex, wuxing)

	if len(dayuns) != 12 {
		t.Errorf("Expected 12 DaYun periods, got %d", len(dayuns))
	}

	// First DaYun: Age 2-11, Branch 4 (Chen)
	if dayuns[0].StartAge != 2 || dayuns[0].Branch != basis.Branch(4) {
		t.Errorf("First DaYun mismatch: age %d, branch %v", dayuns[0].StartAge, dayuns[0].Branch)
	}

	// Second DaYun: Age 12-21, Branch 5 (Si)
	if dayuns[1].StartAge != 12 || dayuns[1].Branch != basis.Branch(5) {
		t.Errorf("Second DaYun mismatch: age %d, branch %v", dayuns[1].StartAge, dayuns[1].Branch)
	}

	// Reverse Test: Yin Male (Gui Year Stem)
	// Gui (9) % 2 != 0 -> isYang = false.
	// Yin Male -> Direction = -1 (Counter-clockwise)
	dayunsRev := CalcDaYun(mingBranch, basis.StemGui, basis.SexMale, wuxing)
	// Second DaYun: Age 12-21, Branch 3 (Mao)
	if dayunsRev[1].Branch != basis.Branch(3) {
		t.Errorf("Reverse DaYun mismatch: branch %v, expected %v", dayunsRev[1].Branch, basis.Branch(3))
	}
}

func TestCalcLiuYue(t *testing.T) {
	// Liu Nian 2024 (Chen Branch = 4)
	// Birth Month 4, Birth Hour Chou (1)
	// Formula: (ln - (bM-1) + bH + 12) % 12
	// (4 - 3 + 1 + 12) % 12 = 2 (Yin) -> Month 1
	lnBranch := basis.Branch(4)
	birthMonth := 4
	birthHour := basis.Branch(1)

	month1 := CalcLiuYue(lnBranch, birthMonth, birthHour, 1)
	if month1 != basis.Branch(2) {
		t.Errorf("Liu Yue 1 mismatch: expected Yin (2), got %v", month1)
	}

	// Month 4: 2 + 3 = 5 (Si)
	month4 := CalcLiuYue(lnBranch, birthMonth, birthHour, 4)
	if month4 != basis.Branch(5) {
		t.Errorf("Liu Yue 4 mismatch: expected Si (5), got %v", month4)
	}
}
