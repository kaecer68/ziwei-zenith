package engine

import (
	"testing"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func TestCalcZiweiStarPos(t *testing.T) {
	tests := []struct {
		name     string
		juValue  int
		lunarDay int
		want     int
	}{
		{"Water2_Day1", 2, 1, 1},  // Chou (1)
		{"Fire6_Day1", 6, 1, 9},   // Shen (9)
		{"Fire6_Day27", 6, 27, 3}, // Mao (3)
		{"Water2_Day2", 2, 2, 2},  // Yin (2)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcZiweiStarPos(tt.juValue, tt.lunarDay)
			if got != tt.want {
				t.Errorf("CalcZiweiStarPos(%d, %d) = %s (%d); want %s (%d)",
					tt.juValue, tt.lunarDay, basis.Branch(got), got, basis.Branch(tt.want), tt.want)
			}
		})
	}
}

func TestPlaceMainStars(t *testing.T) {
	// Test Ziwei at Mao (3)
	// Ziwei Group: ZW(3), TJ(2:Yin), TY(0:Zi), WQ(11:Hai), TT(10:Xu), LZ(7:Wei)
	// Tianfu Symmetry: (4 - 3 + 12) % 12 = 1 (Chou)
	// Tianfu Group: TF(1:Chou), TYin(2:Yin), TL(3:Mao), JM(4:Chen), TX(5:Si), TL(6:Wu), QS(7:Wei), PJ(11:Hai)

	ziweiIdx := 3
	stars := PlaceMainStars(ziweiIdx)

	checkStar := func(b basis.Branch, s basis.Star) {
		found := false
		for _, star := range stars[b] {
			if star == s {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Star %s not found in branch %s", s, b)
		}
	}

	checkStar(basis.Branch(3), basis.StarZiwei)
	checkStar(basis.Branch(2), basis.StarTianji)
	checkStar(basis.Branch(0), basis.StarTaiyang)
	checkStar(basis.Branch(11), basis.StarWuqu)
	checkStar(basis.Branch(10), basis.StarTiantong)
	checkStar(basis.Branch(7), basis.StarLianzhen)

	checkStar(basis.Branch(1), basis.StarTianfu)
	checkStar(basis.Branch(2), basis.StarTaiyin)
	checkStar(basis.Branch(3), basis.StarTanlang)
	checkStar(basis.Branch(4), basis.StarJumen)
	checkStar(basis.Branch(5), basis.StarTianxiang)
	checkStar(basis.Branch(6), basis.StarTianliang)
	checkStar(basis.Branch(7), basis.StarQisha)
	checkStar(basis.Branch(11), basis.StarPojun)
}
