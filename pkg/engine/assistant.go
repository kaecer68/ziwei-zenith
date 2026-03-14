package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func PlaceAssistantStars(mingGong basis.Palace, dayStem basis.Stem, yearStem basis.Stem) map[basis.Palace][]interface{} {
	assistantStars := make(map[basis.Palace][]interface{})

	zuofuPalace := (mingGong + 1) % 12
	assistantStars[basis.Palace(zuofuPalace)] = append(assistantStars[basis.Palace(zuofuPalace)], basis.AuspiciousZuofu)

	youbiPalace := (mingGong + 11) % 12
	assistantStars[basis.Palace(youbiPalace)] = append(assistantStars[basis.Palace(youbiPalace)], basis.AuspiciousYoubi)

	wenchangPalace := (mingGong + 8) % 12
	assistantStars[basis.Palace(wenchangPalace)] = append(assistantStars[basis.Palace(wenchangPalace)], basis.AuspiciousWenchang)

	wenquPalace := (mingGong + 4) % 12
	assistantStars[basis.Palace(wenquPalace)] = append(assistantStars[basis.Palace(wenquPalace)], basis.AuspiciousWenqu)

	tiankuiPalace := (mingGong + 10) % 12
	assistantStars[basis.Palace(tiankuiPalace)] = append(assistantStars[basis.Palace(tiankuiPalace)], basis.AuspiciousTiankui)

	tianyuePalace := (mingGong + 2) % 12
	assistantStars[basis.Palace(tianyuePalace)] = append(assistantStars[basis.Palace(tianyuePalace)], basis.AuspiciousTianyue)

	qingyangPalace := (mingGong + 6) % 12
	assistantStars[basis.Palace(qingyangPalace)] = append(assistantStars[basis.Palace(qingyangPalace)], basis.MaleficQingyang)

	tuoluoPalace := (mingGong + 8) % 12
	assistantStars[basis.Palace(tuoluoPalace)] = append(assistantStars[basis.Palace(tuoluoPalace)], basis.MaleficTuoluo)

	huoxingPalace := (mingGong + 10) % 12
	assistantStars[basis.Palace(huoxingPalace)] = append(assistantStars[basis.Palace(huoxingPalace)], basis.MaleficHuoxing)

	lingxingPalace := (mingGong + 4) % 12
	assistantStars[basis.Palace(lingxingPalace)] = append(assistantStars[basis.Palace(lingxingPalace)], basis.MaleficLingxing)

	dikongPalace := (mingGong + 2) % 12
	assistantStars[basis.Palace(dikongPalace)] = append(assistantStars[basis.Palace(dikongPalace)], basis.MaleficDikong)

	dijiePalace := (mingGong + 6) % 12
	assistantStars[basis.Palace(dijiePalace)] = append(assistantStars[basis.Palace(dijiePalace)], basis.MaleficDijie)

	lucunOffset := (int(yearStem) + 6) % 12
	lucunPalace := basis.Palace((int(mingGong) + lucunOffset) % 12)
	assistantStars[lucunPalace] = append(assistantStars[lucunPalace], basis.LuCun)

	tianmaOffset := (int(dayStem) + 8) % 12
	tianmaPalace := basis.Palace((int(mingGong) + tianmaOffset) % 12)
	assistantStars[tianmaPalace] = append(assistantStars[tianmaPalace], basis.Tianma)

	_ = dayStem

	return assistantStars
}

func PlaceSecondaryStars(mingGong basis.Palace, yearBranch basis.Branch, dayBranch basis.Branch) map[basis.Palace][]interface{} {
	secondaryStars := make(map[basis.Palace][]interface{})

	hongluanOffset := (int(yearBranch) + 10) % 12
	hongluanPalace := basis.Palace((int(mingGong) + hongluanOffset) % 12)
	secondaryStars[hongluanPalace] = append(secondaryStars[hongluanPalace], basis.SecondaryHongluan)

	tianxiOffset := (int(yearBranch) + 4) % 12
	tianxiPalace := basis.Palace((int(mingGong) + tianxiOffset) % 12)
	secondaryStars[tianxiPalace] = append(secondaryStars[tianxiPalace], basis.SecondaryTianxi)

	guchenOffset := (int(yearBranch) / 3) % 12
	guchenPalace := basis.Palace((int(mingGong) + guchenOffset) % 12)
	secondaryStars[guchenPalace] = append(secondaryStars[guchenPalace], basis.SecondaryGuchen)

	guaxuOffset := (int(yearBranch)/3 + 6) % 12
	guaxuPalace := basis.Palace((int(mingGong) + guaxuOffset) % 12)
	secondaryStars[guaxuPalace] = append(secondaryStars[guaxuPalace], basis.SecondaryGuaxu)

	longchiOffset := (int(dayBranch) + 8) % 12
	longchiPalace := basis.Palace((int(mingGong) + longchiOffset) % 12)
	secondaryStars[longchiPalace] = append(secondaryStars[longchiPalace], basis.SecondaryLongchi)

	fenggeOffset := (int(dayBranch) + 4) % 12
	fenggePalace := basis.Palace((int(mingGong) + fenggeOffset) % 12)
	secondaryStars[fenggePalace] = append(secondaryStars[fenggePalace], basis.SecondaryFengge)

	tianxingOffset := (int(yearBranch) + 6) % 12
	tianxingPalace := basis.Palace((int(mingGong) + tianxingOffset) % 12)
	secondaryStars[tianxingPalace] = append(secondaryStars[tianxingPalace], basis.SecondaryTianxing)

	tianyaoOffset := (int(dayBranch) + 10) % 12
	tianyaoPalace := basis.Palace((int(mingGong) + tianyaoOffset) % 12)
	secondaryStars[tianyaoPalace] = append(secondaryStars[tianyaoPalace], basis.SecondaryTianyao)

	jieshenOffset := (int(dayBranch) + 6) % 12
	jieshenPalace := basis.Palace((int(mingGong) + jieshenOffset) % 12)
	secondaryStars[jieshenPalace] = append(secondaryStars[jieshenPalace], basis.SecondaryJieshen)

	tianwuOffset := (int(dayBranch) + 2) % 12
	tianwuPalace := basis.Palace((int(mingGong) + tianwuOffset) % 12)
	secondaryStars[tianwuPalace] = append(secondaryStars[tianwuPalace], basis.SecondaryTianwu)

	santai1Palace := basis.Palace((int(mingGong) + 0) % 12)
	secondaryStars[santai1Palace] = append(secondaryStars[santai1Palace], basis.SecondarySantai)
	santai2Palace := basis.Palace((int(mingGong) + 4) % 12)
	secondaryStars[santai2Palace] = append(secondaryStars[santai2Palace], basis.SecondarySantai)
	santai3Palace := basis.Palace((int(mingGong) + 8) % 12)
	secondaryStars[santai3Palace] = append(secondaryStars[santai3Palace], basis.SecondarySantai)

	bazuo1Palace := basis.Palace((int(mingGong) + 2) % 12)
	secondaryStars[bazuo1Palace] = append(secondaryStars[bazuo1Palace], basis.SecondaryBazuo)
	bazuo2Palace := basis.Palace((int(mingGong) + 6) % 12)
	secondaryStars[bazuo2Palace] = append(secondaryStars[bazuo2Palace], basis.SecondaryBazuo)
	bazuo3Palace := basis.Palace((int(mingGong) + 10) % 12)
	secondaryStars[bazuo3Palace] = append(secondaryStars[bazuo3Palace], basis.SecondaryBazuo)

	_ = yearBranch
	_ = dayBranch

	return secondaryStars
}
