package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type StarPlacement struct {
	Star   basis.Star
	Palace basis.Palace
}

func CalcWuxingJu(yearStem basis.Stem, mingGongBranch basis.Branch) basis.Wuxing {
	table := [10][12]basis.Wuxing{
		{basis.WuxingShuiEr, basis.WuxingLiu, basis.WuxingMuSan, basis.WuxingTuWu, basis.WuxingJinSi, basis.WuxingLiu},
		{basis.WuxingLiu, basis.WuxingTuWu, basis.WuxingJinSi, basis.WuxingMuSan, basis.WuxingShuiEr, basis.WuxingTuWu},
		{basis.WuxingTuWu, basis.WuxingShuiEr, basis.WuxingShuiEr, basis.WuxingJinSi, basis.WuxingLiu, basis.WuxingMuSan},
		{basis.WuxingMuSan, basis.WuxingJinSi, basis.WuxingLiu, basis.WuxingShuiEr, basis.WuxingTuWu, basis.WuxingJinSi},
		{basis.WuxingJinSi, basis.WuxingShuiEr, basis.WuxingTuWu, basis.WuxingLiu, basis.WuxingMuSan, basis.WuxingShuiEr},
	}

	return table[yearStem][mingGongBranch]
}

func CalcZiweiStar(wuxing basis.Wuxing, lunarDay int) int {
	division := wuxing.Value()
	multiple := lunarDay / division
	remainder := lunarDay % division

	startingPoint := multiple
	if remainder%2 == 1 {
		startingPoint = multiple - remainder
	}
	if startingPoint < 0 {
		startingPoint = 0
	}

	return startingPoint % 12
}

func PlaceMainStars(mingGong basis.Palace, wuxing basis.Wuxing, lunarDay int, monthBranch basis.Branch) map[basis.Palace][]basis.Star {
	palaceStars := make(map[basis.Palace][]basis.Star)

	ziweiIdx := CalcZiweiStar(wuxing, lunarDay)
	ziweiPalace := basis.Palace((int(mingGong) + ziweiIdx) % 12)
	palaceStars[ziweiPalace] = append(palaceStars[ziweiPalace], basis.StarZiwei)

	tianfuIdx := (ziweiIdx + 12 - 1) % 12
	tianfuPalace := basis.Palace((int(mingGong) + tianfuIdx) % 12)
	palaceStars[tianfuPalace] = append(palaceStars[tianfuPalace], basis.StarTianfu)

	tianjiIdx := (ziweiIdx + 12 - 1) % 12
	tianjiPalace := basis.Palace((int(mingGong) + tianjiIdx) % 12)
	palaceStars[tianjiPalace] = append(palaceStars[tianjiPalace], basis.StarTianji)

	taiyangIdx := (tianjiIdx + 12 - 2) % 12
	taiyangPalace := basis.Palace((int(mingGong) + taiyangIdx) % 12)
	palaceStars[taiyangPalace] = append(palaceStars[taiyangPalace], basis.StarTaiyang)

	wuquIdx := (taiyangIdx + 12 - 1) % 12
	wuquPalace := basis.Palace((int(mingGong) + wuquIdx) % 12)
	palaceStars[wuquPalace] = append(palaceStars[wuquPalace], basis.StarWuqu)

	tiantongIdx := (wuquIdx + 12 - 1) % 12
	tiantongPalace := basis.Palace((int(mingGong) + tiantongIdx) % 12)
	palaceStars[tiantongPalace] = append(palaceStars[tiantongPalace], basis.StarTiantong)

	lianzhenIdx := (tiantongIdx + 12 - 1) % 12
	lianzhenPalace := basis.Palace((int(mingGong) + lianzhenIdx) % 12)
	palaceStars[lianzhenPalace] = append(palaceStars[lianzhenPalace], basis.StarLianzhen)

	tanlangIdx := (lianzhenIdx + 12 - 1) % 12
	tanlangPalace := basis.Palace((int(mingGong) + tanlangIdx) % 12)
	palaceStars[tanlangPalace] = append(palaceStars[tanlangPalace], basis.StarTanlang)

	jumenIdx := (tanlangIdx + 12 - 1) % 12
	jumenPalace := basis.Palace((int(mingGong) + jumenIdx) % 12)
	palaceStars[jumenPalace] = append(palaceStars[jumenPalace], basis.StarJumen)

	tianxiangIdx := (jumenIdx + 12 - 1) % 12
	tianxiangPalace := basis.Palace((int(mingGong) + tianxiangIdx) % 12)
	palaceStars[tianxiangPalace] = append(palaceStars[tianxiangPalace], basis.StarTianxiang)

	tianliangIdx := (tianxiangIdx + 12 - 1) % 12
	tianliangPalace := basis.Palace((int(mingGong) + tianliangIdx) % 12)
	palaceStars[tianliangPalace] = append(palaceStars[tianliangPalace], basis.StarTianliang)

	qishaIdx := (tianliangIdx + 12 - 1) % 12
	qishaPalace := basis.Palace((int(mingGong) + qishaIdx) % 12)
	palaceStars[qishaPalace] = append(palaceStars[qishaPalace], basis.StarQisha)

	pojunIdx := (qishaIdx + 12 - 1) % 12
	pojunPalace := basis.Palace((int(mingGong) + pojunIdx) % 12)
	palaceStars[pojunPalace] = append(palaceStars[pojunPalace], basis.StarPojun)

	taiyinIdx := (pojunIdx + 12 - 1) % 12
	taiyinPalace := basis.Palace((int(mingGong) + taiyinIdx) % 12)
	palaceStars[taiyinPalace] = append(palaceStars[taiyinPalace], basis.StarTaiyin)

	_ = monthBranch

	return palaceStars
}
