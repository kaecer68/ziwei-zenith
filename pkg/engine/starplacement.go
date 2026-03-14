package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

// CalcZiweiStarPos follows the exact logic of master_engine_v4.js: getZiweiIdx(d, ju)
func CalcZiweiStarPos(juValue int, lunarDay int) int {
	r := lunarDay % juValue
	q := lunarDay / juValue

	if r == 0 {
		return (2 + q - 1) % 12
	}

	x := juValue - r
	nQ := (lunarDay + x) / juValue
	bP := (2 + nQ - 1) % 12

	if x%2 != 0 {
		// If X is Odd: move backward X
		return (bP - x + 12) % 12
	} else {
		// If X is Even: move forward X
		return (bP + x) % 12
	}
}

func PlaceMainStars(ziweiIdx int) map[basis.Branch][]basis.Star {
	stars := make(map[basis.Branch][]basis.Star)

	// Ziwei Group (Counter-clockwise offsets from Ziwei)
	// Master Engine: [['紫微', 0], ['天機', -1], ['太陽', -3], ['武曲', -4], ['天同', -5], ['廉貞', -8]]
	ziweiGroup := []struct {
		star   basis.Star
		offset int
	}{
		{basis.StarZiwei, 0},
		{basis.StarTianji, -1},
		{basis.StarTaiyang, -3},
		{basis.StarWuqu, -4},
		{basis.StarTiantong, -5},
		{basis.StarLianzhen, -8},
	}

	for _, s := range ziweiGroup {
		pos := (ziweiIdx + s.offset + 12) % 12
		stars[basis.Branch(pos)] = append(stars[basis.Branch(pos)], s.star)
	}

	// Tianfu Symmetry: const tf = (4 - zw + 12) % 12;
	tianfuIdx := (4 - ziweiIdx + 12) % 12

	// Tianfu Group (Clockwise offsets from Tianfu)
	// Master Engine: [['天府', 0], ['太陰', 1], ['貪狼', 2], ['巨門', 3], ['天相', 4], ['天梁', 5], ['七殺', 6], ['破軍', 10]]
	tianfuGroup := []struct {
		star   basis.Star
		offset int
	}{
		{basis.StarTianfu, 0},
		{basis.StarTaiyin, 1},
		{basis.StarTanlang, 2},
		{basis.StarJumen, 3},
		{basis.StarTianxiang, 4},
		{basis.StarTianliang, 5},
		{basis.StarQisha, 6},
		{basis.StarPojun, 10},
	}

	for _, s := range tianfuGroup {
		pos := (tianfuIdx + s.offset) % 12
		stars[basis.Branch(pos)] = append(stars[basis.Branch(pos)], s.star)
	}

	return stars
}

func CalcWuxingJu(yearStem basis.Stem, mingBranch basis.Branch) basis.Wuxing {
	// Find MingGong Palace Gan based on Year Stem (tG)
	// Master Engine: const tG = ((yG % 5) * 2 + 2) % 10;
	// yG = 0(Jia)...
	yG := int(yearStem)
	tG := ((yG%5)*2 + 2) % 10

	// GongGans: GAN[(tG + (i - 2 + 12) % 12) % 10]
	// If mingBranch is at index i:
	mingIdx := int(mingBranch)
	mingStemIdx := (tG + (mingIdx-2+12)%12) % 10

	return basis.GetWuxingJu(basis.Stem(mingStemIdx), mingBranch)
}
