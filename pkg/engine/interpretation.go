package engine

import (
	"fmt"
	"strings"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type Interpretation struct {
	OriginPalaceAnalysis string             `json:"origin_palace_analysis"`
	KarmicNarrative      []KarmicStep       `json:"karmic_narrative"`
	SanFangDiagnosis     []SanFangRole      `json:"san_fang_diagnosis"`
	StarDetails          []DeepStarAnalysis `json:"star_details"`
	OriginFlyHua         FlyHuaAnalysis     `json:"origin_fly_hua"`
	TemporalResonance    []ResonancePoint   `json:"temporal_resonance"`
	Summary              string             `json:"summary"`
	CharacterTraits      string             `json:"character_traits"`
	ClassicPatterns      []string           `json:"classic_patterns"`
}

type KarmicStep struct {
	Type   string `json:"type"` // 祿, 權, 科, 忌
	Role   string `json:"role"` // 緣起, 緣力, 緣續, 緣滅
	Star   string `json:"star"`
	Palace string `json:"palace"`
	Desc   string `json:"desc"`
}

type SanFangRole struct {
	Role      string   `json:"role"` // 本宮, 對宮, 氣數位, 財帛位
	Palace    string   `json:"palace"`
	Diagnosis string   `json:"diagnosis"`
	Stars     []string `json:"stars"`
}

type DeepStarAnalysis struct {
	Name       string `json:"name"`
	Verse      string `json:"verse"`
	Positive   string `json:"positive"`
	Negative   string `json:"negative"`
	Remedy     string `json:"remedy"`
	Evolution  string `json:"evolution"`
	Brightness string `json:"brightness"`
}

type FlyHuaAnalysis struct {
	FromPalace string     `json:"from_palace"`
	Stem       string     `json:"stem"`
	Stages     []FlyStage `json:"stages"`
}

type FlyStage struct {
	Type            string          `json:"type"`
	Effect          string          `json:"effect"`
	Star            string          `json:"star"`
	Target          string          `json:"target"`
	Motive          string          `json:"motive"`
	Action          string          `json:"action"`
	Trap            string          `json:"trap"`
	Interpretations MultiSchoolView `json:"interpretations"`
}

type MultiSchoolView struct {
	SanHe   string `json:"sanhe"`
	SiHua   string `json:"sihua"`
	QinTian string `json:"qintian"`
}

type ResonancePoint struct {
	Layer  string `json:"layer"` // 流年, 流月, 流日
	Type   string `json:"type"`  // 祿, 忌
	Star   string `json:"star"`
	Natal  string `json:"natal"` // 本命化XX
	Palace string `json:"palace"`
	Mood   string `json:"mood"`
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

	// 7. Classic Patterns
	interp.ClassicPatterns = detectPatterns(chart, chart.LifePalace.MingGong)

	// 8. Character Analysis
	interp.CharacterTraits = buildCharacterTraits(chart)

	// 9. Summary
	interp.Summary = buildSummary(interp.KarmicNarrative, interp.CharacterTraits)

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
		transType := basis.TransformationType(i)
		for b, stars := range chart.TransformedStars {
			for _, s := range stars {
				if ts, ok := s.(basis.TransformedStar); ok {
					if ts.Transformation == transType {
						steps[i].Star = ts.StarName
						steps[i].Palace = chart.Palaces[b].String()

						// Add deeper meaning based on philosophy_core
						switch transType {
						case basis.TransLu:
							steps[i].Desc = fmt.Sprintf("【%s】代表緣分的前世資助，這份資源由「%s」承載，是今生最好的發力點。", ts.StarName, steps[i].Palace)
						case basis.TransQuan:
							steps[i].Desc = fmt.Sprintf("【%s】是您解決問題的「心機方法」，入於「%s」位，暗示您必須在此領域展現掌控力。", ts.StarName, steps[i].Palace)
						case basis.TransKe:
							steps[i].Desc = fmt.Sprintf("【%s】是安全感與轉機的來源。當人生卡關時，「%s」的智慧與人脈是最好的緩衝。", ts.StarName, steps[i].Palace)
						case basis.TransJi:
							steps[i].Desc = fmt.Sprintf("【%s】是今生的「使命黑洞」。您對「%s」的過度執念或恐懼，是靈魂必須償還並昇華的課題。", ts.StarName, steps[i].Palace)
						}
					}
				}
			}
		}
	}
	return steps
}

func buildResonance(chart *ZiweiChart) []ResonancePoint {
	var points []ResonancePoint

	checkJi := func(layerName string, layerMap map[basis.Branch][]interface{}) {
		for b, stars := range layerMap {
			var layerJiStar string
			for _, s := range stars {
				if ts, ok := s.(basis.TransformedStar); ok && ts.Transformation == basis.TransJi {
					layerJiStar = ts.StarName
				}
			}

			if layerJiStar != "" {
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

	srcName := chart.Palaces[fromBranch].String()
	analysis := FlyHuaAnalysis{
		FromPalace: srcName,
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

		stage := FlyStage{
			Type:   t,
			Star:   starName,
			Target: targetPalace,
			Motive: themes[t].motive,
			Action: themes[t].action,
			Trap:   themes[t].trap,
		}

		// Multi-school data (Logic from master_engine_v4.js)
		stage.Interpretations = MultiSchoolView{
			SanHe:   fmt.Sprintf("【三合派】%s之%s氣飛入%s，代表該宮位對目標宮位的實質拉動。", srcName, t, targetPalace),
			SiHua:   fmt.Sprintf("【飛星四化】「我」主動將%s的能量投射到%s，這是一種主觀的執著與互動。", t, targetPalace),
			QinTian: fmt.Sprintf("【欽天門】前世因今世果，%s星在%s宮展現先天本命之必然軌跡與欠債關係。", starName, targetPalace),
		}

		analysis.Stages = append(analysis.Stages, stage)
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

func detectPatterns(chart *ZiweiChart, target basis.Branch) []string {
	var patterns []string
	idx := int(target)
	sanFangIndices := []int{idx, (idx + 4) % 12, (idx + 8) % 12, (idx + 6) % 12}

	allStars := make(map[string]bool)
	for _, i := range sanFangIndices {
		for _, s := range getStarNames(chart, basis.Branch(i)) {
			allStars[s] = true
		}
	}

	lifeStars := getStarNames(chart, target)
	hasStar := func(name string) bool {
		for _, s := range lifeStars {
			if s == name {
				return true
			}
		}
		return false
	}

	// 1. 機月同梁
	if allStars["天機"] && allStars["太陰"] && allStars["天同"] && allStars["天梁"] {
		patterns = append(patterns, "【機月同梁格】機月同梁作吏人。主幕僚、辦公室行政、制度內發展之命。")
	}

	// 2. 紫府同宮
	if hasStar("紫微") && hasStar("天府") {
		patterns = append(patterns, "【紫府同宮格】紫府同宮，終身福厚。極高基調的生命格局，天生具備穩定的資源支撐。")
	}

	// 3. 殺破狼
	if allStars["七殺"] && allStars["破軍"] && allStars["貪狼"] {
		patterns = append(patterns, "【殺破狼格】變動之星系，一生多震盪與革新。")
	}

	// 4. 石中隱玉 (子午巨門)
	if hasStar("巨門") && (target == basis.BranchZi || target == basis.BranchWu) {
		patterns = append(patterns, "【石中隱玉格】子午巨門，石中隱玉。需磨礪（祿科引動）方能顯貴，不宜過早出位。")
	}

	// 5. 雄宿乾元 (未申廉貞)
	if hasStar("廉貞") && (target == basis.BranchWei || target == basis.BranchShen) {
		patterns = append(patterns, "【雄宿乾元格】廉貞在未申宮守命，展示高度的外交手腕與開創力。")
	}

	// 6. 壽星入廟 (梁居午位)
	if hasStar("天梁") && target == basis.BranchWu {
		patterns = append(patterns, "【壽星入廟格】梁居午位，官資清顯。人格高潔，長壽且具備崇高地位。")
	}

	// 7. 英星入廟 (子午破軍)
	if hasStar("破軍") && (target == basis.BranchZi || target == basis.BranchWu) {
		patterns = append(patterns, "【英星入廟格】子午破軍，加官進祿。具備強大的開創力與顛覆魄力。")
	}

	// 8. 日月併明
	if allStars["太陽"] && allStars["太陰"] {
		// Simplified check for brightening
		patterns = append(patterns, "【日月併明格】日月並明，佐九重於堯殿。代表得位返照，貴人顯赫。")
	}

	return patterns
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

func buildSummary(steps []KarmicStep, traits string) string {
	formula := "能量循環解析中..."
	if len(steps) >= 4 && steps[0].Palace != "" && steps[3].Palace != "" {
		formula = fmt.Sprintf("緣起於【%s】之%s，緣滅於【%s】之%s。",
			steps[0].Palace, steps[0].Star, steps[3].Palace, steps[3].Star)
	}
	return fmt.Sprintf("%s %s", traits, formula)
}

func buildCharacterTraits(chart *ZiweiChart) string {
	stars := getStarNames(chart, chart.LifePalace.MingGong)
	if len(stars) == 0 {
		return "您的命盤呈現出一種極具彈性的能量場，容易受環境影響而展現多樣面貌。"
	}

	var mainTraits []string
	for _, s := range stars {
		if e, ok := StarEssenceTable[s]; ok {
			mainTraits = append(mainTraits, e.Trait)
		}
	}

	if len(mainTraits) > 0 {
		return fmt.Sprintf("您的生命底色以「%s」為主旋律，展現出獨特的精神海拔。", strings.Join(mainTraits, "、"))
	}

	return "您具備極強的適應力，命宮格局暗示了一種深藏不露的生命張力。"
}
