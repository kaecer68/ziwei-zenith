package engine

import (
	"fmt"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type ZiweiEngine struct{}

type ZiweiChart struct {
	LifePalace       LifePalace
	Palaces          map[basis.Branch]basis.Palace
	Stars            map[basis.Branch][]basis.Star
	AssistantStars   map[basis.Branch][]interface{}
	SecondaryStars   map[basis.Branch][]interface{}
	ChangSheng       map[basis.Branch]basis.ChangShengStage
	BoShi            map[basis.Branch]basis.BoShiStar
	TransformedStars map[basis.Branch][]interface{}
	DaYun            []basis.DaYun
	LiuNian          basis.LiuNian
	LiuYue           basis.Branch
	LiuRi            basis.Branch
	Patterns         []basis.Pattern
	Wuxing           basis.Wuxing
	NaYin            basis.NaYin
	YearPillar       basis.Pillar
	MonthPillar      basis.Pillar
	DayPillar        basis.Pillar
	HourPillar       basis.Pillar
	LiuNianStars     map[basis.Branch][]interface{}
	LiuYueStars      map[basis.Branch][]interface{}
	LiuRiStars       map[basis.Branch][]interface{}
	StarBrightness   []basis.StarBrightness
	OriginPalace     basis.Branch // 來因宮
	PalaceGans       map[basis.Branch]basis.Stem
	Interpretation   Interpretation
}

func New() *ZiweiEngine {
	return &ZiweiEngine{}
}

func (e *ZiweiEngine) BuildChart(birth basis.BirthInfo) (*ZiweiChart, error) {
	yearPillar := birth.YearPillar

	absMonth := birth.GetAbsMonth()
	calcMonth := absMonth
	// Master Split Rule: Leap Month Split at 15th
	if birth.IsLeap() && birth.LunarDay >= 16 {
		calcMonth = (absMonth % 12) + 1
	}

	lp := CalcLifePalace(calcMonth, basis.Branch(birth.HourBranch))

	palaceToBranch := BuildPalaces(lp.MingGong)
	branchToPalace := make(map[basis.Branch]basis.Palace)
	for p, b := range palaceToBranch {
		branchToPalace[b] = p
	}

	// Calculate Palace Gans (Based on Year Stem)
	yG := int(yearPillar.Stem)
	tG := ((yG%5)*2 + 2) % 10
	palaceGans := make(map[basis.Branch]basis.Stem)
	for i := 0; i < 12; i++ {
		b := basis.Branch(i)
		stemIdx := (tG + (i-2+12)%12) % 10
		palaceGans[b] = basis.Stem(stemIdx)
	}

	var origin basis.Branch
	for b, g := range palaceGans {
		if g == yearPillar.Stem {
			if b == basis.Branch(0) || b == basis.Branch(1) {
				if b == yearPillar.Branch {
					origin = b
					break
				}
				continue
			}
			origin = b
		}
	}

	wuxing := CalcWuxingJu(yearPillar.Stem, lp.MingGong)
	naYin := basis.CalcNaYin(yearPillar.Stem, yearPillar.Branch)

	ziweiIdx := CalcZiweiStarPos(wuxing.Value(), birth.LunarDay)
	stars := PlaceMainStars(ziweiIdx)

	chart := &ZiweiChart{
		LifePalace:   lp,
		Palaces:      branchToPalace,
		Stars:        stars,
		Wuxing:       wuxing,
		NaYin:        naYin,
		YearPillar:   yearPillar,
		MonthPillar:  birth.MonthPillar,
		DayPillar:    birth.DayPillar,
		HourPillar:   birth.HourPillar,
		OriginPalace: origin,
		PalaceGans:   palaceGans,
	}

	chart.AssistantStars = PlaceAssistantStars(yearPillar.Stem, calcMonth, basis.Branch(birth.HourBranch))
	chart.SecondaryStars = PlaceSecondaryStars(yearPillar.Stem, yearPillar.Branch, calcMonth, birth.LunarDay, basis.Branch(birth.HourBranch), lp.MingGong, lp.ShenGong)
	chart.ChangSheng = CalcChangShengStages(wuxing, birth.Sex, yearPillar.Stem)
	chart.BoShi = CalcBoShiStars(locateLuCun(yearPillar.Stem), birth.Sex, yearPillar.Stem)
	chart.TransformedStars = PlaceTransformationStars(yearPillar.Stem, chart)
	chart.DaYun = CalcDaYun(lp.MingGong, yearPillar.Stem, birth.Sex, wuxing)

	// 獲取當前系統時間計算流年、流月、流日
	now := time.Now()
	currentPillar := getCurrentPillar(now)

	chart.LiuNian = CalcLiuNian(currentPillar.yearBranch, now.Year())
	chart.LiuYue = CalcLiuYue(chart.LiuNian.Branch, birth.LunarMonth, basis.Branch(birth.HourBranch), currentPillar.lunarMonth)
	chart.LiuRi = CalcLiuRi(chart.LiuYue, currentPillar.lunarDay)

	// Temporal Transformations (使用當前時間的干支)
	chart.LiuNianStars = PlaceLayeredTransformations(currentPillar.yearStem, chart)
	chart.LiuYueStars = PlaceLayeredTransformations(currentPillar.monthStem, chart)
	chart.LiuRiStars = PlaceLayeredTransformations(currentPillar.dayStem, chart)

	chart.Patterns = DetectPatterns(chart)
	chart.StarBrightness = CalcStarBrightness(chart)
	chart.Interpretation = GenerateInterpretation(chart)

	return chart, nil
}

func (c *ZiweiChart) String() string {
	str := "紫微斗數命盤\n"
	str += fmt.Sprintf("年柱: %s\n", c.YearPillar)
	str += fmt.Sprintf("月柱: %s\n", c.MonthPillar)
	str += fmt.Sprintf("日柱: %s\n", c.DayPillar)
	str += fmt.Sprintf("時柱: %s\n", c.HourPillar)
	str += fmt.Sprintf("五行局: %s\n", c.Wuxing)
	str += fmt.Sprintf("納音: %s\n", c.NaYin)
	str += fmt.Sprintf("來因宮: %s (%s)\n", c.OriginPalace, c.Palaces[c.OriginPalace])
	str += fmt.Sprintf("命宮在: %s\n", c.LifePalace.MingGong)
	str += fmt.Sprintf("身宮在: %s\n", c.LifePalace.ShenGong)
	str += fmt.Sprintf("流年命宮: %s | 流年四化 (%s)\n", c.LiuNian.Branch, c.YearPillar.Stem)
	str += fmt.Sprintf("流月命宮: %s | 流月四化 (%s)\n", c.LiuYue, c.MonthPillar.Stem)
	str += fmt.Sprintf("流日命宮: %s\n", c.LiuRi)

	str += "\n能量循環 (祿隨忌走):\n"
	for _, s := range c.Interpretation.KarmicNarrative {
		str += fmt.Sprintf("  【%s】%s: %s落入%s — %s\n", s.Type, s.Role, s.Star, s.Palace, s.Desc)
	}
	str += fmt.Sprintf("\n大師總論: %s\n", c.Interpretation.Summary)

	str += "\n三方四正專業診斷 (Synthesis Diagnosis):\n"
	for _, r := range c.Interpretation.SanFangDiagnosis {
		str += fmt.Sprintf("  %-10s [%s]: %s\n", r.Role, r.Palace, r.Diagnosis)
	}

	str += fmt.Sprintf("\n來因宮動態因果 (Origin Fly-Hua: %s宮發射):\n", c.Interpretation.OriginFlyHua.FromPalace)
	for _, s := range c.Interpretation.OriginFlyHua.Stages {
		str += fmt.Sprintf("  • %s (%s) -> 飛入【%s】:\n", s.Type, s.Motive, s.Target)
		str += fmt.Sprintf("    動作: %s | 陷阱: %s\n", s.Action, s.Trap)
	}

	if len(c.Interpretation.TemporalResonance) > 0 {
		str += "\n時空感應 (Temporal Resonance - 歲運疊併):\n"
		for _, r := range c.Interpretation.TemporalResonance {
			str += fmt.Sprintf("  [%s][%s] %s 與 %s 於【%s】宮會合:\n", r.Layer, r.Type, r.Star, r.Natal, r.Palace)
			str += fmt.Sprintf("    解析: %s\n", r.Mood)
		}
	}

	str += "\n星曜深度解析 (Ancient Wisdom):\n"
	for _, sd := range c.Interpretation.StarDetails {
		if sd.Verse != "" {
			str += fmt.Sprintf("  ---【 %s 】---\n", sd.Name)
			str += fmt.Sprintf("  賦文: %s\n", sd.Verse)
			str += fmt.Sprintf("  正面: %s\n", sd.Positive)
			str += fmt.Sprintf("  負面: %s\n", sd.Negative)
			str += fmt.Sprintf("  修行: %s\n", sd.Remedy)
			if sd.Evolution != "" {
				str += fmt.Sprintf("  %s\n", sd.Evolution)
			}
		}
	}

	if len(c.Patterns) > 0 {
		str += "\n檢出格局:\n"
		for _, p := range c.Patterns {
			str += fmt.Sprintf("  [%s] %s — %s\n", p.Level, p.Name, p.Description)
		}
	}

	str += "\n宮位分布 (地支物理順序):\n"
	for i := 0; i < 12; i++ {
		b := basis.Branch(i)
		pType := c.Palaces[b]
		pGan := c.PalaceGans[b]

		starStr := ""
		for _, s := range c.Stars[b] {
			starStr += s.String() + " "
		}
		for _, s := range c.AssistantStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				starStr += strer.String() + " "
			}
		}
		for _, s := range c.SecondaryStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				starStr += strer.String() + " "
			}
		}
		for _, s := range c.TransformedStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				starStr += strer.String() + " "
			}
		}
		for _, s := range c.LiuNianStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				starStr += "[流年" + strer.String() + "] "
			}
		}
		for _, s := range c.LiuYueStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				starStr += "[流月" + strer.String() + "] "
			}
		}
		for _, s := range c.LiuRiStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				starStr += "[流日" + strer.String() + "] "
			}
		}

		if starStr == "" {
			starStr = "(空宮)"
		}
		prefix := ""
		if b == c.OriginPalace {
			prefix += "[來因]"
		}
		if b == c.LiuNian.Branch {
			prefix += "[流年]"
		}
		if b == c.LiuYue {
			prefix += "[流月]"
		}
		if b == c.LiuRi {
			prefix += "[流日]"
		}

		extra := ""
		if stage, ok := c.ChangSheng[b]; ok {
			extra += stage.String()
		}
		if boShi, ok := c.BoShi[b]; ok {
			if extra != "" {
				extra += " "
			}
			extra += boShi.String()
		}
		if extra != "" {
			starStr += " {" + extra + "}"
		}

		str += fmt.Sprintf("  %s%s%s(%s): %s\n", prefix, pGan, b, pType, starStr)
	}

	return str
}

// currentPillar 存储当前时间的干支信息
type currentPillar struct {
	yearStem   basis.Stem
	yearBranch basis.Branch
	monthStem  basis.Stem
	dayStem    basis.Stem
	lunarMonth int
	lunarDay   int
}

// getCurrentPillar 获取当前系统时间的干支信息
func getCurrentPillar(now time.Time) currentPillar {
	pt := celestial.NewPrecisionTime(now)
	pillar := zodiac.GetAstrologicalPillar(pt)

	// 获取农历日期
	jd := celestial.TimeToJD(now)
	lunarEngine := &zodiac.LunarEngine{}
	lunarDate := lunarEngine.GetLunarDate(jd)

	return currentPillar{
		yearStem:   basis.Stem(pillar.Year.StemIndex),
		yearBranch: basis.Branch(pillar.Year.BranchIndex),
		monthStem:  basis.Stem(pillar.Month.StemIndex),
		dayStem:    basis.Stem(pillar.Day.StemIndex),
		lunarMonth: lunarDate.Month,
		lunarDay:   lunarDate.Day,
	}
}
