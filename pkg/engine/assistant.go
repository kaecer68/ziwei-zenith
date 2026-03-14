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

func PlaceSecondaryStars(yearBranch basis.Branch, lunarMonth int, lunarDay int, hourBranch basis.Branch) map[basis.Branch][]interface{} {
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
