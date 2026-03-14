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
