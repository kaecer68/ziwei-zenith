package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func PlaceAssistantStars(yearStem basis.Stem, lunarMonth int, hourBranch basis.Branch) map[basis.Branch][]interface{} {
	stars := make(map[basis.Branch][]interface{})

	// 1. Zuofu (Clockwise from Chen) - Master: 4 + (calcMonth - 1)
	zuofuIdx := (4 + lunarMonth - 1) % 12
	stars[basis.Branch(zuofuIdx)] = append(stars[basis.Branch(zuofuIdx)], basis.AuspiciousZuofu)

	// 2. Youbi (Counter-clockwise from Xu) - Master: 10 - (calcMonth - 1)
	youbiIdx := (10 - (lunarMonth - 1) + 12) % 12
	stars[basis.Branch(youbiIdx)] = append(stars[basis.Branch(youbiIdx)], basis.AuspiciousYoubi)

	// 3. Wenchang (Counter-clockwise from Xu) - Master: 10 - timeIdx
	wenchangIdx := (10 - int(hourBranch) + 12) % 12
	stars[basis.Branch(wenchangIdx)] = append(stars[basis.Branch(wenchangIdx)], basis.AuspiciousWenchang)

	// 4. Wenqu (Clockwise from Chen) - Master: 4 + timeIdx
	wenquIdx := (4 + int(hourBranch)) % 12
	stars[basis.Branch(wenquIdx)] = append(stars[basis.Branch(wenquIdx)], basis.AuspiciousWenqu)

	// 5. Tiankui & Tianyue (Based on Year Stem) - Ported Table
	kuiYue := map[basis.Stem][2]int{
		basis.StemJia: {1, 7}, basis.StemYi: {0, 8}, basis.StemBing: {11, 9},
		basis.StemDing: {11, 9}, basis.StemWu: {1, 7}, basis.StemJi: {0, 8},
		basis.StemGeng: {1, 7}, basis.StemXin: {6, 2}, basis.StemRen: {3, 5},
		basis.StemGui: {3, 5},
	}[yearStem]
	stars[basis.Branch(kuiYue[0])] = append(stars[basis.Branch(kuiYue[0])], basis.AuspiciousTiankui)
	stars[basis.Branch(kuiYue[1])] = append(stars[basis.Branch(kuiYue[1])], basis.AuspiciousTianyue)

	// 6. Lucun, Qingyang, Tuoluo (Based on Year Stem)
	lucunIdx := map[basis.Stem]int{
		basis.StemJia: 2, basis.StemYi: 3, basis.StemBing: 5, basis.StemDing: 6,
		basis.StemWu: 5, basis.StemJi: 6, basis.StemGeng: 8, basis.StemXin: 9,
		basis.StemRen: 11, basis.StemGui: 0,
	}[yearStem]
	stars[basis.Branch(lucunIdx)] = append(stars[basis.Branch(lucunIdx)], basis.LuCun)
	stars[basis.Branch((lucunIdx+1)%12)] = append(stars[basis.Branch((lucunIdx+1)%12)], basis.MaleficQingyang)
	stars[basis.Branch((lucunIdx-1+12)%12)] = append(stars[basis.Branch((lucunIdx-1+12)%12)], basis.MaleficTuoluo)

	return stars
}

func PlaceSecondaryStars(yearStem basis.Stem, yearBranch basis.Branch, lunarMonth int, lunarDay int, hourBranch basis.Branch, mingBranch basis.Branch, shenBranch basis.Branch) map[basis.Branch][]interface{} {
	stars := make(map[basis.Branch][]interface{})

	// 1. Fire star & Bell star (Year Branch + Hour)
	// Master FIRE_BELL_START: '寅': { h: '丑', l: '卯' }...
	fbStarts := map[basis.Branch][2]int{
		2: {1, 3}, 6: {1, 3}, 10: {1, 3}, // Yin, Wu, Xu -> h:Chou(1), l:Mao(3)
		8: {2, 10}, 0: {2, 10}, 4: {2, 10}, // Shen, Zi, Chen -> h:Yin(2), l:Xu(10)
		5: {3, 10}, 9: {3, 10}, 1: {3, 10}, // Si, You, Chou -> h:Mao(3), l:Xu(10)
		11: {9, 10}, 3: {9, 10}, 7: {9, 10}, // Hai, Mao, Wei -> h:You(9), l:Xu(10)
	}
	starts := fbStarts[yearBranch]
	stars[basis.Branch((starts[0]+int(hourBranch))%12)] = append(stars[basis.Branch((starts[0]+int(hourBranch))%12)], basis.MaleficHuoxing)
	stars[basis.Branch((starts[1]+int(hourBranch))%12)] = append(stars[basis.Branch((starts[1]+int(hourBranch))%12)], basis.MaleficLingxing)

	// 2. Void stats: Dikong & Dijie
	stars[basis.Branch((11-int(hourBranch)+12)%12)] = append(stars[basis.Branch((11-int(hourBranch)+12)%12)], basis.MaleficDikong)
	stars[basis.Branch((11+int(hourBranch))%12)] = append(stars[basis.Branch((11+int(hourBranch))%12)], basis.MaleficDijie)

	// 3. Hongluan & Tianxi (Year Branch)
	hongluanIdx := (3 - int(yearBranch) + 12) % 12
	stars[basis.Branch(hongluanIdx)] = append(stars[basis.Branch(hongluanIdx)], basis.SecondaryHongluan)
	stars[basis.Branch((hongluanIdx+6)%12)] = append(stars[basis.Branch((hongluanIdx+6)%12)], basis.SecondaryTianxi)

	// 4. Longchi & Fengge (Year Branch)
	longchiIdx := (4 + int(yearBranch)) % 12
	stars[basis.Branch(longchiIdx)] = append(stars[basis.Branch(longchiIdx)], basis.SecondaryLongchi)
	fenggeIdx := (10 - int(yearBranch) + 12) % 12
	stars[basis.Branch(fenggeIdx)] = append(stars[basis.Branch(fenggeIdx)], basis.SecondaryFengge)

	// 5. Tianma (Year Branch)
	maTable := map[basis.Branch]int{
		2: 8, 6: 8, 10: 8, // Yin, Wu, Xu -> Shen(8)
		8: 2, 0: 2, 4: 2, // Shen, Zi, Chen -> Yin(2)
		5: 11, 9: 11, 1: 11, // Si, You, Chou -> Hai(11)
		11: 5, 3: 5, 7: 5, // Hai, Mao, Wei -> Si(5)
	}
	stars[basis.Branch(maTable[yearBranch])] = append(stars[basis.Branch(maTable[yearBranch])], basis.Tianma)

	// 6. Jieshen (Month) - Master uses jieTable[displayMonth] (1-indexed)
	// Month indices: 1:8, 2:8, 3:10, 4:10, 5:0, 6:0, 7:2, 8:2, 9:4, 10:4, 11:6, 12:6
	jieTable := map[int]int{
		1: 8, 2: 8, 3: 10, 4: 10, 5: 0, 6: 0, 7: 2, 8: 2, 9: 4, 10: 4, 11: 6, 12: 6,
	}
	stars[basis.Branch(jieTable[lunarMonth])] = append(stars[basis.Branch(jieTable[lunarMonth])], basis.SecondaryJieshen)

	// 7. Tianwu (Month)
	// wuTable: 1:5, 2:8, 3:2, 4:11, 5:5, 6:8, 7:2, 8:11, 9:5, 10:8, 11:2, 12:11
	wuTable := map[int]int{
		1: 5, 2: 8, 3: 2, 4: 11, 5: 5, 6: 8, 7: 2, 8: 11, 9: 5, 10: 8, 11: 2, 12: 11,
	}
	stars[basis.Branch(wuTable[lunarMonth])] = append(stars[basis.Branch(wuTable[lunarMonth])], basis.SecondaryTianwu)

	// 8. Tianxing & Tianyao (Month)
	stars[basis.Branch((9+lunarMonth-1)%12)] = append(stars[basis.Branch((9+lunarMonth-1)%12)], basis.SecondaryTianxing)
	stars[basis.Branch((1+lunarMonth-1)%12)] = append(stars[basis.Branch((1+lunarMonth-1)%12)], basis.SecondaryTianyao)

	// 9. Guchen & Guaxu (Year Branch)
	guchenTable := map[basis.Branch]int{
		11: 2, 0: 2, 1: 2, // Hai, Zi, Chou -> Yin(2)
		2: 5, 3: 5, 4: 5, // Yin, Mao, Chen -> Si(5)
		5: 8, 6: 8, 7: 8, // Si, Wu, Wei -> Shen(8)
		8: 11, 9: 11, 10: 11, // Shen, You, Xu -> Hai(11)
	}
	guaxuTable := map[basis.Branch]int{
		11: 10, 0: 10, 1: 10, // Hai, Zi, Chou -> Xu(10)
		2: 1, 3: 1, 4: 1, // Yin, Mao, Chen -> Chou(1)
		5: 4, 6: 4, 7: 4, // Si, Wu, Wei -> Chen(4)
		8: 7, 9: 7, 10: 7, // Shen, You, Xu -> Wei(7)
	}
	stars[basis.Branch(guchenTable[yearBranch])] = append(stars[basis.Branch(guchenTable[yearBranch])], basis.SecondaryGuchen)
	stars[basis.Branch(guaxuTable[yearBranch])] = append(stars[basis.Branch(guaxuTable[yearBranch])], basis.SecondaryGuaxu)

	// 10. Santai & Bazuo (Based on Zuofu/Youbi + Day)
	zuofuIdx := (4 + lunarMonth - 1) % 12
	youbiIdx := (10 - (lunarMonth - 1) + 12) % 12
	santaiIdx := (zuofuIdx + (lunarDay - 1)) % 12
	bazuoIdx := (youbiIdx - (lunarDay - 1) + 120) % 12
	stars[basis.Branch(santaiIdx)] = append(stars[basis.Branch(santaiIdx)], basis.SecondarySantai)
	stars[basis.Branch(bazuoIdx)] = append(stars[basis.Branch(bazuoIdx)], basis.SecondaryBazuo)

	// 11. Xianchi (Year Branch)
	xianchiTable := map[basis.Branch]int{
		8: 9, 0: 9, 4: 9, // Shen-Zi-Chen -> You(9)
		2: 3, 6: 3, 10: 3, // Yin-Wu-Xu -> Mao(3)
		5: 6, 9: 6, 1: 6, // Si-You-Chou -> Wu(6)
		11: 0, 3: 0, 7: 0, // Hai-Mao-Wei -> Zi(0)
	}
	stars[basis.Branch(xianchiTable[yearBranch])] = append(stars[basis.Branch(xianchiTable[yearBranch])], basis.SecondaryXianchi)

	// 12. Tianyue (Month)
	tianyueTable := map[int]int{1: 10, 2: 5, 3: 4, 4: 5, 5: 7, 6: 8, 7: 10, 8: 1, 9: 11, 10: 11, 11: 3, 12: 6}
	stars[basis.Branch(tianyueTable[lunarMonth])] = append(stars[basis.Branch(tianyueTable[lunarMonth])], basis.SecondaryTianyue)

	// 13. Yinsha (Month)
	yinshaTable := map[int]int{1: 2, 2: 0, 3: 10, 4: 8, 5: 6, 6: 4, 7: 2, 8: 0, 9: 10, 10: 8, 11: 6, 12: 4}
	stars[basis.Branch(yinshaTable[lunarMonth])] = append(stars[basis.Branch(yinshaTable[lunarMonth])], basis.SecondaryYinsha)

	// 14. Taifu & Fenggao (Hour)
	stars[basis.Branch((6+int(hourBranch))%12)] = append(stars[basis.Branch((6+int(hourBranch))%12)], basis.SecondaryTianning)
	stars[basis.Branch((2+int(hourBranch))%12)] = append(stars[basis.Branch((2+int(hourBranch))%12)], basis.SecondaryFenggao)

	// 15. 天哭 (Year Branch): (6 - yearBranch + 12) % 12
	tiankuIdx := (6 - int(yearBranch) + 12) % 12
	stars[basis.Branch(tiankuIdx)] = append(stars[basis.Branch(tiankuIdx)], basis.SecondaryTianku)

	// 16. 天虛 (Year Branch): (6 + yearBranch) % 12
	tianxuIdx := (6 + int(yearBranch)) % 12
	stars[basis.Branch(tianxuIdx)] = append(stars[basis.Branch(tianxuIdx)], basis.SecondaryTianxu)

	// 17. 華蓋 (Year Branch): 三合局墓庫位
	huagaiTable := map[basis.Branch]int{
		8: 4, 0: 4, 4: 4, // 申子辰 -> 辰(4)
		2: 10, 6: 10, 10: 10, // 寅午戌 -> 戌(10)
		5: 1, 9: 1, 1: 1, // 巳酉丑 -> 丑(1)
		11: 7, 3: 7, 7: 7, // 亥卯未 -> 未(7)
	}
	stars[basis.Branch(huagaiTable[yearBranch])] = append(stars[basis.Branch(huagaiTable[yearBranch])], basis.SecondaryHuagai)

	// 18. 破碎 (Year Branch): {子午卯酉}->巳, {丑辰未戌}->丑, {寅巳申亥}->酉
	posuiTable := map[basis.Branch]int{
		0: 5, 6: 5, 3: 5, 9: 5, // 子午卯酉 -> 巳(5)
		1: 1, 4: 1, 7: 1, 10: 1, // 丑辰未戌 -> 丑(1)
		2: 9, 5: 9, 8: 9, 11: 9, // 寅巳申亥 -> 酉(9)
	}
	stars[basis.Branch(posuiTable[yearBranch])] = append(stars[basis.Branch(posuiTable[yearBranch])], basis.SecondaryPosui)

	// 19. 蜚廉 (Year Branch): (yearBranch + 8) % 12
	feilianIdx := (int(yearBranch) + 8) % 12
	stars[basis.Branch(feilianIdx)] = append(stars[basis.Branch(feilianIdx)], basis.SecondaryFeilian)

	// 20. 天官 (Year Stem)
	tianguanTable := map[basis.Stem]int{
		basis.StemJia: 7, basis.StemYi: 4, basis.StemBing: 5, basis.StemDing: 2,
		basis.StemWu: 3, basis.StemJi: 9, basis.StemGeng: 11, basis.StemXin: 9,
		basis.StemRen: 10, basis.StemGui: 6,
	}
	stars[basis.Branch(tianguanTable[yearStem])] = append(stars[basis.Branch(tianguanTable[yearStem])], basis.SecondaryTianguan)

	// 21. 天福 (Year Stem)
	tianfuTable := map[basis.Stem]int{
		basis.StemJia: 9, basis.StemYi: 8, basis.StemBing: 0, basis.StemDing: 11,
		basis.StemWu: 3, basis.StemJi: 2, basis.StemGeng: 6, basis.StemXin: 5,
		basis.StemRen: 6, basis.StemGui: 5,
	}
	stars[basis.Branch(tianfuTable[yearStem])] = append(stars[basis.Branch(tianfuTable[yearStem])], basis.SecondaryTianfu)

	// 22. 恩光 (文昌位 + 日數 - 2): 以文昌宮起初一，順數至生日，再退一位
	wenchangPos := (10 - int(hourBranch) + 12) % 12
	enguangIdx := (wenchangPos + lunarDay - 2 + 12) % 12
	stars[basis.Branch(enguangIdx)] = append(stars[basis.Branch(enguangIdx)], basis.SecondaryEnguang)

	// 23. 天貴 (文曲位 + 日數 - 2): 以文曲宮起初一，順數至生日，再退一位
	wenquPos := (4 + int(hourBranch)) % 12
	tianguiIdx := (wenquPos + lunarDay - 2 + 12) % 12
	stars[basis.Branch(tianguiIdx)] = append(stars[basis.Branch(tianguiIdx)], basis.SecondaryTiangui)

	// 24. 天才 (命宮支 + 年支)
	tiancaiIdx := (int(mingBranch) + int(yearBranch)) % 12
	stars[basis.Branch(tiancaiIdx)] = append(stars[basis.Branch(tiancaiIdx)], basis.SecondaryTiancai)

	// 25. 天壽 (身宮支 + 年支)
	tianshouIdx := (int(shenBranch) + int(yearBranch)) % 12
	stars[basis.Branch(tianshouIdx)] = append(stars[basis.Branch(tianshouIdx)], basis.SecondaryTianshou)

	// 26. 天德 (Year Stem)
	tiandeTable := map[basis.Stem]int{
		basis.StemJia: 9, basis.StemYi: 8, basis.StemBing: 0, basis.StemDing: 11,
		basis.StemWu: 3, basis.StemJi: 2, basis.StemGeng: 6, basis.StemXin: 5,
		basis.StemRen: 9, basis.StemGui: 8,
	}
	stars[basis.Branch(tiandeTable[yearStem])] = append(stars[basis.Branch(tiandeTable[yearStem])], basis.SecondaryTiande)

	// 27. 月德 (Year Stem)
	yuedeTable := map[basis.Stem]int{
		basis.StemJia: 5, basis.StemYi: 6, basis.StemBing: 8, basis.StemDing: 9,
		basis.StemWu: 11, basis.StemJi: 0, basis.StemGeng: 2, basis.StemXin: 3,
		basis.StemRen: 5, basis.StemGui: 6,
	}
	stars[basis.Branch(yuedeTable[yearStem])] = append(stars[basis.Branch(yuedeTable[yearStem])], basis.SecondaryYuede)

	// 28. 天傷 (僕役宮 = Palace index 7 = mingBranch - 7)
	tianshangIdx := (int(mingBranch) - 7 + 12) % 12
	stars[basis.Branch(tianshangIdx)] = append(stars[basis.Branch(tianshangIdx)], basis.SecondaryTianshang)

	// 29. 天使 (疾厄宮 = Palace index 5 = mingBranch - 5)
	tianshiIdx := (int(mingBranch) - 5 + 12) % 12
	stars[basis.Branch(tianshiIdx)] = append(stars[basis.Branch(tianshiIdx)], basis.SecondaryTianshi)

	// 30. 天空 (Year Branch): (yearBranch + 1) % 12
	tiankongIdx := (int(yearBranch) + 1) % 12
	stars[basis.Branch(tiankongIdx)] = append(stars[basis.Branch(tiankongIdx)], basis.SecondaryTiankong)

	// 31. 劫煞 (Year Branch): 三合局
	jieshaTable := map[basis.Branch]int{
		8: 5, 0: 5, 4: 5, // 申子辰 -> 巳(5)
		2: 11, 6: 11, 10: 11, // 寅午戌 -> 亥(11)
		5: 2, 9: 2, 1: 2, // 巳酉丑 -> 寅(2)
		11: 8, 3: 8, 7: 8, // 亥卯未 -> 申(8)
	}
	stars[basis.Branch(jieshaTable[yearBranch])] = append(stars[basis.Branch(jieshaTable[yearBranch])], basis.SecondaryJiesha)

	return stars
}

func PlaceTransformationStars(yearStem basis.Stem, chart *ZiweiChart) map[basis.Branch][]interface{} {
	table, ok := basis.TransformationTable[yearStem]
	if !ok {
		return nil
	}

	result := make(map[basis.Branch][]interface{})

	for i, starName := range table {
		transType := basis.TransformationType(i)
		for b, stars := range chart.Stars {
			for _, s := range stars {
				if s.String() == starName {
					result[b] = append(result[b], basis.TransformedStar{
						Transformation: transType,
						StarName:       starName,
					})
				}
			}
		}
		for b, stars := range chart.AssistantStars {
			for _, s := range stars {
				if strer, ok := s.(interface{ String() string }); ok {
					if strer.String() == starName {
						result[b] = append(result[b], basis.TransformedStar{
							Transformation: transType,
							StarName:       starName,
						})
					}
				}
			}
		}
	}
	return result
}
func PlaceLayeredTransformations(stem basis.Stem, chart *ZiweiChart) map[basis.Branch][]interface{} {
	table, ok := basis.TransformationTable[stem]
	if !ok {
		return nil
	}

	result := make(map[basis.Branch][]interface{})

	for i, starName := range table {
		transType := basis.TransformationType(i)
		// Find the star in the chart and mark its transformation
		for b, stars := range chart.Stars {
			for _, s := range stars {
				if s.String() == starName {
					result[b] = append(result[b], basis.TransformedStar{
						Transformation: transType,
						StarName:       starName,
					})
				}
			}
		}
		// Also check in assistant stars if they are main stars typically but formatted differently
		// For simplicity, we mostly track the 14 main stars + assistant stars like Zuo/You/Wen/Wen
		for b, stars := range chart.AssistantStars {
			for _, s := range stars {
				if strer, ok := s.(interface{ String() string }); ok {
					if strer.String() == starName {
						result[b] = append(result[b], basis.TransformedStar{
							Transformation: transType,
							StarName:       starName,
						})
					}
				}
			}
		}
	}
	return result
}
