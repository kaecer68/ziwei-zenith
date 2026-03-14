package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type Pattern struct {
	Name        string
	Description string
	Level       string
}

func DetectPatterns(stars map[basis.Palace][]basis.Star, palaces map[basis.Palace]basis.Branch, assistantStars map[basis.Palace][]interface{}) []Pattern {
	var patterns []Pattern

	mingGongPalace := basis.PalaceMing
	for p := range stars {
		if palaces[p] == basis.BranchZi {
			mingGongPalace = p
			break
		}
	}

	hasStar := func(p basis.Palace, s basis.Star) bool {
		for _, star := range stars[p] {
			if star == s {
				return true
			}
		}
		return false
	}

	hasStarInAnyPalace := func(palaces []basis.Palace, s basis.Star) bool {
		for _, p := range palaces {
			if hasStar(p, s) {
				return true
			}
		}
		return false
	}

	getOppositePalace := func(p basis.Palace) basis.Palace {
		return basis.Palace((int(p) + 6) % 12)
	}

	getSanFangSiZheng := func(p basis.Palace) []basis.Palace {
		return []basis.Palace{
			basis.Palace((int(p) + 4) % 12),
			basis.Palace((int(p) + 8) % 12),
			basis.Palace((int(p) + 10) % 12),
		}
	}

	mingStars := stars[mingGongPalace]

	_ = mingStars

	if hasStar(mingGongPalace, basis.StarZiwei) && hasStar(mingGongPalace, basis.StarTianfu) {
		patterns = append(patterns, Pattern{Name: "紫府同宮", Description: "紫微天府同宮，尊貴無極", Level: "甲"})
	}

	if hasStar(mingGongPalace, basis.StarZiwei) && hasStar(mingGongPalace, basis.StarPojun) {
		patterns = append(patterns, Pattern{Name: "紫破同宮", Description: "紫微破軍同宮，權力慾望極強", Level: "甲"})
	}

	if hasStar(mingGongPalace, basis.StarZiwei) && hasStar(mingGongPalace, basis.StarTianxiang) {
		patterns = append(patterns, Pattern{Name: "紫相拱照", Description: "紫微天相拱照，有權有勢", Level: "甲"})
	}

	oppositePalace := getOppositePalace(mingGongPalace)
	sanfang := getSanFangSiZheng(mingGongPalace)

	hasQisha := hasStarInAnyPalace(sanfang, basis.StarQisha)
	hasPojun := hasStarInAnyPalace(sanfang, basis.StarPojun)
	hasTanlang := hasStarInAnyPalace(sanfang, basis.StarTanlang)
	if hasQisha && hasPojun && hasTanlang {
		patterns = append(patterns, Pattern{Name: "殺破狼格", Description: "七殺破軍貪狼三方四正會照", Level: "甲"})
	}

	hasTianji := hasStarInAnyPalace(sanfang, basis.StarTianji)
	hasTaiyinAtSf := hasStarInAnyPalace(sanfang, basis.StarTaiyin)
	hasTiantong := hasStarInAnyPalace(sanfang, basis.StarTiantong)
	hasTianliang := hasStarInAnyPalace(sanfang, basis.StarTianliang)
	if hasTianji && hasTaiyinAtSf && hasTiantong && hasTianliang {
		patterns = append(patterns, Pattern{Name: "機月同梁格", Description: "天機太陰天同天梁會合", Level: "甲"})
	}

	hasWuqu := hasStarInAnyPalace(sanfang, basis.StarWuqu)
	hasLianzhen := hasStarInAnyPalace(sanfang, basis.StarLianzhen)
	if hasStar(mingGongPalace, basis.StarZiwei) && hasWuqu && hasLianzhen && hasStar(mingGongPalace, basis.StarTianfu) {
		patterns = append(patterns, Pattern{Name: "紫武廉府", Description: "紫微武曲廉貞天府會合", Level: "甲"})
	}

	if hasStar(mingGongPalace, basis.StarTianfu) && hasStar(mingGongPalace, basis.StarTianxiang) {
		patterns = append(patterns, Pattern{Name: "府相朝垣", Description: "天府天相會合於命宮", Level: "乙"})
	}

	migrationPalace := basis.PalaceQianYi
	for p, b := range palaces {
		if b == basis.BranchYin {
			migrationPalace = p
			break
		}
	}

	hasTaiyang := hasStar(mingGongPalace, basis.StarTaiyang)
	hasTaiyinAtMing := hasStar(mingGongPalace, basis.StarTaiyin)
	if hasTaiyang && hasTaiyinAtMing {
		patterns = append(patterns, Pattern{Name: "日月拱照", Description: "太陽太陰拱照命宮", Level: "甲"})
	}

	if hasStar(migrationPalace, basis.StarTaiyang) {
		patterns = append(patterns, Pattern{Name: "日月反背", Description: "太陽在遷移宮且落陷", Level: "辛"})
	}

	hasLucun := false
	hasTianma := false
	for p, ss := range assistantStars {
		for _, as := range ss {
			if as == basis.LuCun {
				hasLucun = true
			}
			if as == basis.Tianma {
				hasTianma = true
			}
		}
		_ = p
	}
	if hasLucun && hasTianma {
		patterns = append(patterns, Pattern{Name: "祿馬交馳", Description: "祿存與天馬同宮或對宮", Level: "乙"})
	}

	hasHongluan := false
	hasTianyao := false
	for p, ss := range assistantStars {
		for _, as := range ss {
			if as == basis.SecondaryHongluan {
				hasHongluan = true
			}
			if as == basis.SecondaryTianyao {
				hasTianyao = true
			}
		}
		_ = p
	}
	if hasHongluan || hasTianyao {
		if len(mingStars) > 0 {
			patterns = append(patterns, Pattern{Name: "桃花犯主", Description: "桃花星與主星同宮", Level: "丙"})
		}
	}

	if len(mingStars) == 0 {
		patterns = append(patterns, Pattern{Name: "空宮", Description: "命宮無主星", Level: "丁"})
	}

	_ = oppositePalace
	_ = migrationPalace

	return patterns
}
