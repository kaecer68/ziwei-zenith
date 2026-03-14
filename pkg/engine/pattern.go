package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func DetectPatterns(chart *ZiweiChart) []basis.Pattern {
	var result []basis.Pattern

	branchOfPalace := make(map[basis.Palace]basis.Branch)
	for b, p := range chart.Palaces {
		branchOfPalace[p] = b
	}

	getStars := func(p basis.Palace) []string {
		b := branchOfPalace[p]
		var names []string
		for _, s := range chart.Stars[b] {
			names = append(names, s.String())
		}
		for _, s := range chart.AssistantStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				names = append(names, strer.String())
			}
		}
		return names
	}

	hasStar := func(p basis.Palace, name string) bool {
		stars := getStars(p)
		for _, s := range stars {
			if s == name {
				return true
			}
		}
		return false
	}

	getSanFangSiZheng := func(p basis.Palace) []basis.Palace {
		return []basis.Palace{
			p,
			basis.Palace((int(p) + 4) % 12),
			basis.Palace((int(p) + 8) % 12),
			basis.Palace((int(p) + 6) % 12),
		}
	}

	target := basis.PalaceMing
	sfzs := getSanFangSiZheng(target)

	hasInSFZ := func(name string) bool {
		for _, p := range sfzs {
			if hasStar(p, name) {
				return true
			}
		}
		return false
	}

	// Pattern Logic Aligned with master_engine_v4.js

	// 1. 紫府同宮
	if hasStar(target, "紫微") && hasStar(target, "天府") {
		result = append(result, basis.Patterns[0])
	}
	// 2. 紫破同宮
	if hasStar(target, "紫微") && hasStar(target, "破軍") {
		result = append(result, basis.Patterns[1])
	}
	// 3. 殺破狼格
	if hasInSFZ("七殺") && hasInSFZ("破軍") && hasInSFZ("貪狼") {
		result = append(result, basis.Patterns[3])
	}
	// 4. 機月同梁格
	if hasInSFZ("天機") && hasInSFZ("太陰") && hasInSFZ("天同") && hasInSFZ("天梁") {
		result = append(result, basis.Patterns[4])
	}
	// 5. 府相朝垣
	if hasInSFZ("天府") && hasInSFZ("天相") {
		result = append(result, basis.Patterns[6])
	}
	// 6. 日月拱照
	if hasInSFZ("太陽") && hasInSFZ("太陰") {
		result = append(result, basis.Patterns[7])
	}
	// 7. 火貪格 / 鈴貪格
	if hasInSFZ("貪狼") {
		if hasInSFZ("火星") {
			result = append(result, basis.Patterns[11])
		} else if hasInSFZ("鈴星") {
			result = append(result, basis.Patterns[12])
		}
	}
	// 8. 祿馬交馳
	if hasInSFZ("祿存") && hasInSFZ("天馬") {
		result = append(result, basis.Patterns[9])
	}
	// 9. 三奇加會 (Hua Lu, Hua Quan, Hua Ke in SFZ)
	luCount, quanCount, keCount := 0, 0, 0
	for _, p := range sfzs {
		b := branchOfPalace[p]
		for _, s := range chart.TransformedStars[b] {
			if ts, ok := s.(basis.TransformedStar); ok {
				switch ts.Transformation {
				case basis.TransLu:
					luCount++
				case basis.TransQuan:
					quanCount++
				case basis.TransKe:
					keCount++
				}
			}
		}
	}
	if luCount > 0 && quanCount > 0 && keCount > 0 {
		result = append(result, basis.Pattern{Name: "三奇加會", Description: "化祿、化權、化科在三方四正會照，主名望與富貴", Level: "甲"})
	}

	// 10. 雙祿交流 / 雙祿齊臨
	hasLuCun := hasInSFZ("祿存")
	hasHuaLu := false
	for _, p := range sfzs {
		b := branchOfPalace[p]
		for _, s := range chart.TransformedStars[b] {
			if ts, ok := s.(basis.TransformedStar); ok && ts.Transformation == basis.TransLu {
				hasHuaLu = true
				break
			}
		}
	}
	if hasLuCun && hasHuaLu {
		result = append(result, basis.Pattern{Name: "雙祿交流", Description: "祿存與化祿在三方四正會合，財源廣進", Level: "甲"})
	}

	// 11. 命無正曜
	if len(chart.Stars[branchOfPalace[target]]) == 0 {
		result = append(result, basis.Patterns[14]) // 空宮
	}

	// 12. 廉貞貪狼在巳亥 (Ancient Warning)
	if hasStar(target, "廉貞") && hasStar(target, "貪狼") && (branchOfPalace[target] == basis.BranchSi || branchOfPalace[target] == basis.BranchHai) {
		result = append(result, basis.Pattern{Name: "廉貞貪狼在巳亥", Description: "男浪蕩，女多淫；刑戮及身 (需自律)", Level: "辛"})
	}

	// 13. 天機巨門在卯酉 (Classic)
	if hasStar(target, "天機") && hasStar(target, "巨門") && (branchOfPalace[target] == basis.BranchMao || branchOfPalace[target] == basis.BranchYou) {
		result = append(result, basis.Pattern{Name: "天機巨門在卯酉", Description: "機巨同臨，位至三公；雖富貴，不免酒色", Level: "乙"})
	}

	return result
}
