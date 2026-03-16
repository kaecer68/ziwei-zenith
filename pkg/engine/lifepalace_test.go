package engine

import (
	"testing"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func TestCalcLifePalace(t *testing.T) {
	tests := []struct {
		name       string
		lunarMonth int
		hourBranch basis.Branch
		wantMing   basis.Branch
		wantShen   basis.Branch
	}{
		{
			name:       "1990-06-15_10:00_Lunar_5_Si",
			lunarMonth: 5,
			hourBranch: basis.Branch(5),  // 巳
			wantMing:   basis.Branch(1),  // 丑
			wantShen:   basis.Branch(11), // 戌
		},
		{
			name:       "Jan_Zi_Hour",
			lunarMonth: 1,
			hourBranch: basis.Branch(0), // 子
			wantMing:   basis.Branch(2), // 寅
			wantShen:   basis.Branch(2), // 寅
		},
		{
			name:       "Jul_Wu_Hour",
			lunarMonth: 7,
			hourBranch: basis.Branch(6), // 午
			wantMing:   basis.Branch(2), // 寅
			wantShen:   basis.Branch(2), // 寅
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalcLifePalace(tt.lunarMonth, tt.hourBranch)
			if got.MingGong != tt.wantMing {
				t.Errorf("CalcLifePalace() MingGong = %v, want %v", got.MingGong, tt.wantMing)
			}
			if got.ShenGong != tt.wantShen {
				t.Errorf("CalcLifePalace() ShenGong = %v, want %v", got.ShenGong, tt.wantShen)
			}
		})
	}
}

func TestBuildPalaces(t *testing.T) {
	// If Ming is Chou (1), then:
	// 0: 命 -> Chou (1)
	// 1: 兄 -> Zi (0)
	// 2: 妻 -> Hai (11)
	// ...
	ming := basis.Branch(1)
	palaces := BuildPalaces(ming)

	if palaces[basis.PalaceMing] != basis.Branch(1) {
		t.Errorf("PalaceMing should be Branch(1), got %v", palaces[basis.PalaceMing])
	}
	if palaces[basis.PalaceXiongDi] != basis.Branch(0) {
		t.Errorf("PalaceXiongDi should be Branch(0), got %v", palaces[basis.PalaceXiongDi])
	}
	if palaces[basis.PalaceFuQi] != basis.Branch(11) {
		t.Errorf("PalaceFuQi should be Branch(11), got %v", palaces[basis.PalaceFuQi])
	}
}
