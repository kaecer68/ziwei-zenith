package engine

import (
	"testing"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func TestGenerateInterpretation(t *testing.T) {
	// Setup a complex chart to trigger interpretations
	chart := &ZiweiChart{
		LifePalace: LifePalace{MingGong: basis.BranchChen, ShenGong: basis.BranchShen},
		Palaces: map[basis.Branch]basis.Palace{
			basis.BranchChen: basis.PalaceMing,
			basis.BranchShen: basis.PalaceFuDe,
			basis.BranchZi:   basis.PalaceGuanLu,
			basis.BranchXu:   basis.PalaceQianYi,
		},
		Stars: map[basis.Branch][]basis.Star{
			basis.BranchChen: {basis.StarZiwei},
			basis.BranchZi:   {basis.StarTianji},
		},
		AssistantStars: make(map[basis.Branch][]interface{}),
		TransformedStars: map[basis.Branch][]interface{}{
			basis.BranchChen: {basis.TransformedStar{StarName: "紫微", Transformation: basis.TransQuan}},
			basis.BranchZi:   {basis.TransformedStar{StarName: "天機", Transformation: basis.TransJi}},
			basis.BranchXu:   {basis.TransformedStar{StarName: "太陽", Transformation: basis.TransLu}},
			basis.BranchShen: {basis.TransformedStar{StarName: "武曲", Transformation: basis.TransKe}},
		},
		LiuNianStars: map[basis.Branch][]interface{}{
			basis.BranchZi: {basis.TransformedStar{StarName: "天機", Transformation: basis.TransJi}},
		},
		OriginPalace: basis.BranchChen,
		PalaceGans: map[basis.Branch]basis.Stem{
			basis.BranchChen: basis.StemRen,
		},
	}

	interp := GenerateInterpretation(chart)

	// Check Summary
	if interp.Summary == "" || interp.Summary == "能量循環解析中..." {
		t.Errorf("Interpretation summary should be generated, got: %s", interp.Summary)
	}

	// Check Resonance (Ling Nian Ji + Natal Ji in Zi)
	foundResonance := false
	for _, r := range interp.TemporalResonance {
		if r.Palace == "官祿宮" && r.Type == "疊忌" {
			foundResonance = true
			break
		}
	}
	if !foundResonance {
		t.Errorf("Temporal resonance (疊忌) not detected for Zi branch")
	}

	// Check San Fang Analysis
	if len(interp.SanFangDiagnosis) == 0 {
		t.Errorf("San Fang diagnosis should not be empty")
	}

	// Check Fly-Hua
	if interp.OriginFlyHua.Stem == "" {
		t.Errorf("Origin Fly-Hua analysis should have stem data")
	}
}
