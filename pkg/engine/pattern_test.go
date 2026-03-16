package engine

import (
	"testing"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func TestDetectPatterns(t *testing.T) {
	// Setup a chart for "Sha-Po-Lang" (七殺、破軍、貪狼 in SFZ)
	chart := &ZiweiChart{
		Palaces: map[basis.Branch]basis.Palace{
			basis.BranchChen: basis.PalaceMing,
			basis.BranchShen: basis.PalaceCaiBang,
			basis.BranchZi:   basis.PalaceGuanLu,
			basis.BranchXu:   basis.PalaceQianYi,
		},
		Stars: map[basis.Branch][]basis.Star{
			basis.BranchChen: {basis.StarQisha},
			basis.BranchShen: {basis.StarPojun},
			basis.BranchZi:   {basis.StarTanlang},
		},
		AssistantStars:   make(map[basis.Branch][]interface{}),
		TransformedStars: make(map[basis.Branch][]interface{}),
	}

	patterns := DetectPatterns(chart)
	found := false
	for _, p := range patterns {
		if p.Name == "殺破狼格" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Sha-Po-Lang pattern not detected")
	}

	// Test "Zi-Fu" (紫微、天府 in Ming)
	chart.Stars[basis.BranchChen] = []basis.Star{basis.StarZiwei, basis.StarTianfu}
	patterns2 := DetectPatterns(chart)
	foundZiFu := false
	for _, p := range patterns2 {
		if p.Name == "紫府同宮" {
			foundZiFu = true
			break
		}
	}
	if !foundZiFu {
		t.Errorf("Zi-Fu pattern not detected")
	}
}
