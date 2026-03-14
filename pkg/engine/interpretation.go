package engine

import (
	"fmt"
	"strings"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type Interpretation struct {
	OriginPalaceAnalysis string
	KarmicNarrative      []KarmicStep
	SanFangDiagnosis     []SanFangRole
	StarDetails          []DeepStarAnalysis
	OriginFlyHua         FlyHuaAnalysis
	TemporalResonance    []ResonancePoint
	Summary              string
}

type KarmicStep struct {
	Type   string // 祿, 權, 科, 忌
	Role   string // 緣起, 緣力, 緣續, 緣滅
	Star   string
	Palace string
	Desc   string
}

type SanFangRole struct {
	Role      string // 本宮, 對宮, 氣數位, 財帛位
	Palace    string
	Diagnosis string
	Stars     []string
}

type DeepStarAnalysis struct {
	Name         string
	Verse        string
	Positive     string
	Negative     string
	Remedy       string
	Evolution    string
	Brightness   string
}

type FlyHuaAnalysis struct {
	FromPalace string
	Stem       string
	Stages     []FlyStage
}

type FlyStage struct {
	Type   string
	Effect string
	Star   string
	Target string
	Motive string
	Action string
	Trap   string
}

type ResonancePoint struct {
	Layer string // 流年, 流月, 流日
	Type  string // 祿, 忌
	Star  string
	Natal string // 本命化XX
	Palace string
	Mood   string
}

func GenerateInterpretation(chart *ZiweiChart) Interpretation {
	interp := Interpretation{}

	// 1. Origin Palace Analysis
	originBranch := chart.OriginPalace
	originName := chart.Palaces[originBranch].String()
	interp.OriginPalaceAnalysis = fmt.Sprintf("來因宮位在【%s】，這是您靈魂的原始設定與能量發射源。", originName)

	// 2. Karmic Narrative (祿隨忌走)
	interp.KarmicNarrative = buildKarmicStory(chart)

	// 3. San Fang Si Zheng Synthesis (Life Palace)
	interp.SanFangDiagnosis = buildSanFangDiagnosis(chart, chart.LifePalace.MingGong)

	// 4. Origin Palace Fly-Hua Analysis
	interp.OriginFlyHua = buildFlyHuaAnalysis(chart, originBranch)

	// 5. Temporal Resonance (疊併感應)
	interp.TemporalResonance = buildResonance(chart)

	// 6. Star Deep Analysis for Life Palace
	interp.StarDetails = buildStarDetails(chart, chart.LifePalace.MingGong)

	// 7. Summary
	interp.Summary = buildSummary(interp.KarmicNarrative)

	return interp
}

func buildKarmicStory(chart *ZiweiChart) []KarmicStep {
	steps := []KarmicStep{
		{Type: "祿", Role: "緣起", Desc: "天降的資源與機緣。故事的起點。"},
		{Type: "權", Role: "緣力", Desc: "達成目標的方法與力量。業力的推動。"},
		{Type: "科", Role: "緣續", Desc: "秩序、緩衝與名聲。智慧與轉機。"},
		{Type: "忌", Role: "緣滅", Desc: "最終的結果、責任與執著。收藏與還債。"},
	}

	for i := range steps {
		for b, stars := range chart.TransformedStars {
			for _, s := range stars {
				if ts, ok := s.(basis.TransformedStar); ok {
					if int(ts.Transformation) == i {
						steps[i].Star = ts.StarName
						steps[i].Palace = chart.Palaces[b].String()
					}
				}
			}
		}
	}
	return steps
}

func buildResonance(chart *ZiweiChart) []ResonancePoint {
	var points []ResonancePoint
	
	// Check for Natal Ji vs Temporal Ji intersections
	checkJi := func(layerName string, layerMap map[basis.Branch][]interface{}) {
		for b, stars := range layerMap {
			var layerJiStar string
			for _, s := range stars {
				if ts, ok := s.(basis.TransformedStar); ok && ts.Transformation == basis.TransJi {
					layerJiStar = ts.StarName
				}
			}
			
			if layerJiStar != "" {
				// Search for Natal Ji in same palace
				for _, ns := range chart.TransformedStars[b] {
					if nts, ok := ns.(basis.TransformedStar); ok && nts.Transformation == basis.TransJi {
						points = append(points, ResonancePoint{
							Layer:  layerName,
							Type:   "疊忌",
							Star:   layerJiStar,
							Natal:  "本命忌",
							Palace: chart.Palaces[b].String(),
							Mood:   "⚠️ 雙重執念匯聚！壓力能量倍增，此宮位領域需謹慎處理，防範負面爆發。",
						})
					}
				}
				// Search for Natal Lu in same palace (Lu-Ji meeting)
				for _, ns := range chart.TransformedStars[b] {
					if nts, ok := ns.(basis.TransformedStar); ok && nts.Transformation == basis.TransLu {
						points = append(points, ResonancePoint{
							Layer:  layerName,
							Type:   "祿忌交戰",
							Star:   layerJiStar,
							Natal:  "本命祿",
							Palace: chart.Palaces[b].String(),
							Mood:   "🌀 祿忌同宮，得失參半。好處中帶有隱患，成敗繫於一念之間。",
						})
					}
				}
			}
		}
	}

	checkJi("流年", chart.LiuNianStars)
	checkJi("流月", chart.LiuYueStars)
	checkJi("流日", chart.LiuRiStars)

	return points
}

func buildFlyHuaAnalysis(chart *ZiweiChart, fromBranch basis.Branch) FlyHuaAnalysis {
	stem := chart.PalaceGans[fromBranch]
	hua, ok := basis.TransformationTable[stem]
	if !ok {
		return FlyHuaAnalysis{}
	}

	analysis := FlyHuaAnalysis{
		FromPalace: chart.Palaces[fromBranch].String(),
		Stem:       stem.String(),
	}

	types := []string{"祿", "權", "科", "忌"}
	themes := map[string]struct {
		motive string
		action string
		trap   string
	}{
		"祿": {"緣起能量", "順勢接納，資源分享", "暴殄天賦，錯過機會"},
		"權": {"掌控慾望", "自帶影響力，執行力強", "剛愎自用，過度干涉"},
		"科": {"理智秩序", "化解危機，合約名聲", "好面子，猶豫不決"},
		"忌": {"因果債務", "承擔責任，正視弱點", "鑽牛角尖，陷入死循環"},
	}

	for i, t := range types {
		starName := hua[i]
		targetPalace := findStarLocation(chart, starName)
		
		analysis.Stages = append(analysis.Stages, FlyStage{
			Type:   t,
			Star:   starName,
			Target: targetPalace,
			Motive: themes[t].motive,
			Action: themes[t].action,
			Trap:   themes[t].trap,
		})
	}
	return analysis
}

func findStarLocation(chart *ZiweiChart, starName string) string {
	for b, stars := range chart.Stars {
		for _, s := range stars {
			if s.String() == starName {
				return chart.Palaces[b].String()
			}
		}
	}
	// Check assistant stars
	for b, stars := range chart.AssistantStars {
		for _, s := range stars {
			if strer, ok := s.(interface{ String() string }); ok {
				if strer.String() == starName {
					return chart.Palaces[b].String()
				}
			}
		}
	}
	return "未知"
}

func buildSanFangDiagnosis(chart *ZiweiChart, target basis.Branch) []SanFangRole {
	idx := int(target)
	roles := []struct {
		name  string
		pIdx  int
		label string
	}{
		{"命宮(體)", idx, "本質能量"},
		{"遷移位(用)", (idx + 6) % 12, "外部環境與表現"},
		{"官祿位(氣)", (idx + 4) % 12, "執行邏輯與事業位"},
		{"財帛位(數)", (idx + 8) % 12, "資源獲取與生命支撐"},
	}

	var result []SanFangRole
	for _, r := range roles {
		branch := basis.Branch(r.pIdx)
		roleStars := getStarNames(chart, branch)
		diag := synthesizeDiagnosis(chart.Palaces[branch].String(), roleStars, r.name)
		result = append(result, SanFangRole{
			Role:      r.name,
			Palace:    chart.Palaces[branch].String(),
			Stars:     roleStars,
			Diagnosis: diag,
		})
	}
	return result
}

func synthesizeDiagnosis(palaceName string, stars []string, roleLabel string) string {
	if len(stars) == 0 {
		return "此宮位目前氣場較為空靈，受周邊宮位引動影響較大。"
	}

	var mainStar string
	for _, s := range stars {
		if _, ok := StarEssenceTable[s]; ok {
			mainStar = s
			break
		}
	}

	if mainStar != "" {
		se := StarEssenceTable[mainStar]
		base := fmt.Sprintf("主導星曜為【%s】。展現出「%s」的特質。", mainStar, se.Trait)
		if strings.Contains(roleLabel, "用") {
			base += fmt.Sprintf(" 面對環境需警惕「%s」，建議「%s」。", se.Shadow, se.Remedy)
		} else if strings.Contains(roleLabel, "氣") {
			base += fmt.Sprintf(" 執行邏輯為：%s。這是成敗關鍵。", se.Core)
		}
		return base
	}
	return fmt.Sprintf("由輔星「%s」主導，呈現多元且微觀的變動氣場。", strings.Join(stars, "、"))
}

func buildStarDetails(chart *ZiweiChart, branch basis.Branch) []DeepStarAnalysis {
	var result []DeepStarAnalysis
	stars := getStarNames(chart, branch)
	for _, name := range stars {
		analysis := DeepStarAnalysis{Name: name}
		if se, ok := StarEssenceTable[name]; ok {
			analysis.Verse = se.AncientVerse
			analysis.Positive = se.Positive
			analysis.Negative = se.Negative
			analysis.Remedy = se.Remedy
		}
		if origin, ok := HistoricalOriginTable[name]; ok {
			analysis.Evolution = fmt.Sprintf("溯源：原名「%s」。%s", origin.OriginalName, origin.Essence)
		}
		result = append(result, analysis)
	}
	return result
}

func getStarNames(chart *ZiweiChart, b basis.Branch) []string {
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

func buildSummary(steps []KarmicStep) string {
	if steps[0].Palace != "" && steps[3].Palace != "" {
		return fmt.Sprintf("命運公式：緣起於【%s】之%s，緣滅於【%s】之%s。所得終為債，此為祿隨忌走之必然。",
			steps[0].Palace, steps[0].Star, steps[3].Palace, steps[3].Star)
	}
	return "能量循環尚不完整，需進一步推演。"
}
