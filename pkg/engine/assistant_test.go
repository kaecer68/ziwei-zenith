package engine

import (
	"testing"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func TestPlaceAssistantStars(t *testing.T) {
	// Ren Stem (9), Month 4, Hour Chou (1)
	yearStem := basis.StemRen
	lunarMonth := 4
	hourBranch := basis.HourBranch(1)

	stars := PlaceAssistantStars(yearStem, lunarMonth, basis.Branch(hourBranch))

	checkStar := func(b basis.Branch, s interface{}) {
		found := false
		str := ""
		if strer, ok := s.(interface{ String() string }); ok {
			str = strer.String()
		}

		for _, star := range stars[b] {
			if strer, ok := star.(interface{ String() string }); ok {
				if strer.String() == str {
					found = true
					break
				}
			}
		}
		if !found {
			t.Errorf("Star %v not found in branch %s", s, b)
		}
	}

	checkStar(basis.Branch(7), basis.AuspiciousZuofu)
	checkStar(basis.Branch(7), basis.AuspiciousYoubi)
	checkStar(basis.Branch(9), basis.AuspiciousWenchang)
	checkStar(basis.Branch(5), basis.AuspiciousWenqu)
	checkStar(basis.Branch(3), basis.AuspiciousTiankui)
	checkStar(basis.Branch(5), basis.AuspiciousTianyue)
	checkStar(basis.Branch(11), basis.LuCun)
	checkStar(basis.Branch(0), basis.MaleficQingyang)
	checkStar(basis.Branch(10), basis.MaleficTuoluo)
}

func TestPlaceTransformationStars(t *testing.T) {
	// Setup a simple chart with some stars
	chart := &ZiweiChart{
		Stars: map[basis.Branch][]basis.Star{
			basis.BranchChen: {basis.StarZiwei},
			basis.BranchWei:  {basis.StarLianzhen},
			basis.BranchMao:  {basis.StarTianji},
		},
		AssistantStars: map[basis.Branch][]interface{}{
			basis.BranchWei: {basis.AuspiciousZuofu},
		},
	}

	// Ren Year Stem: 梁紫左武 (Tianliang-Lu, Ziwei-Quan, Zuofu-Ke, Wuqu-Ji)
	transformed := PlaceTransformationStars(basis.StemRen, chart)

	checkTrans := func(b basis.Branch, starName string, trans basis.TransformationType) {
		found := false
		for _, s := range transformed[b] {
			ts := s.(basis.TransformedStar)
			if ts.StarName == starName && ts.Transformation == trans {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Transformation %s-%s not found in branch %s", starName, trans, b)
		}
	}

	checkTrans(basis.BranchChen, "紫微", basis.TransformationQuan)
	checkTrans(basis.BranchWei, "左輔", basis.TransformationKe)
}
