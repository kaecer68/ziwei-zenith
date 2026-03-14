package engine

import (
	"fmt"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type ZiweiEngine struct{}

type ZiweiChart struct {
	LifePalace       LifePalace
	Palaces          map[basis.Palace]basis.Branch
	Stars            map[basis.Palace][]basis.Star
	AssistantStars   map[basis.Palace][]interface{}
	SecondaryStars   map[basis.Palace][]interface{}
	TransformedStars map[basis.Palace][]interface{}
	DaYun            []basis.DaYun
	LiuNian          []basis.LiuNian
	LiuYue           []basis.LiuYue
	LiuRi            []basis.LiuRi
	Wuxing           basis.Wuxing
	NaYin            basis.NaYin
	YearPillar       basis.Pillar
	MonthPillar      basis.Pillar
	DayPillar        basis.Pillar
	HourPillar       basis.Pillar
}

func New() *ZiweiEngine {
	return &ZiweiEngine{}
}

func (e *ZiweiEngine) BuildChart(birth BirthInfo) (*ZiweiChart, error) {
	yearPillar := birth.YearPillar
	monthPillar := birth.MonthPillar
	dayPillar := birth.DayPillar

	lifePalace := CalcLifePalace(birth.LunarMonth, birth.LunarDay, basis.Branch(birth.HourBranch))

	wuxing := CalcWuxingJu(yearPillar.Stem, basis.Branch(lifePalace.MingGong))
	naYin := basis.CalcNaYin(dayPillar.Stem, dayPillar.Branch)

	palaces := BuildPalaces(lifePalace.MingGong)

	stars := PlaceMainStars(lifePalace.MingGong, wuxing, birth.LunarDay, monthPillar.Branch)

	assistantStars := PlaceAssistantStars(lifePalace.MingGong, dayPillar.Stem, yearPillar.Stem)

	secondaryStars := PlaceSecondaryStars(lifePalace.MingGong, yearPillar.Branch, dayPillar.Branch)

	transformedStars := PlaceTransformationStars(lifePalace.MingGong, yearPillar.Stem, stars, palaces)

	dayuns := CalcDaYun(lifePalace.MingGong, birth.Sex, yearPillar.Stem, birth.LunarYear)

	liunians := CalcLiuNian(lifePalace.MingGong, dayPillar.Stem, birth.LunarYear)

	liuyues := CalcLiuYue(lifePalace.MingGong, yearPillar.Branch, birth.LunarMonth)

	liuris := CalcLiuRi(lifePalace.MingGong, dayPillar.Stem, birth.LunarDay)

	return &ZiweiChart{
		LifePalace:       lifePalace,
		Palaces:          palaces,
		Stars:            stars,
		AssistantStars:   assistantStars,
		SecondaryStars:   secondaryStars,
		TransformedStars: transformedStars,
		DaYun:            dayuns,
		LiuNian:          liunians,
		LiuYue:           liuyues,
		LiuRi:            liuris,
		Wuxing:           wuxing,
		NaYin:            naYin,
		YearPillar:       yearPillar,
		MonthPillar:      monthPillar,
		DayPillar:        dayPillar,
		HourPillar:       birth.HourPillar,
	}, nil
}

func (c *ZiweiChart) String() string {
	str := "紫微斗數命盤\n"
	str += fmt.Sprintf("年柱: %s\n", c.YearPillar)
	str += fmt.Sprintf("月柱: %s\n", c.MonthPillar)
	str += fmt.Sprintf("日柱: %s\n", c.DayPillar)
	str += fmt.Sprintf("時柱: %s\n", c.HourPillar)
	str += fmt.Sprintf("五行局: %s\n", c.Wuxing)
	str += fmt.Sprintf("納音: %s\n", c.NaYin)
	str += fmt.Sprintf("命宮: %s\n", c.LifePalace.MingGong)
	str += fmt.Sprintf("身宮: %s\n", c.LifePalace.ShenGong)
	str += "\n宮位分布:\n"
	for i := 0; i < 12; i++ {
		palace := basis.Palace(i)
		branch := c.Palaces[palace]
		stars := c.Stars[palace]
		assistantStars := c.AssistantStars[palace]
		secondaryStars := c.SecondaryStars[palace]
		transformedStars := c.TransformedStars[palace]

		starStr := ""
		for _, s := range stars {
			starStr += s.String() + " "
		}
		for _, as := range assistantStars {
			switch v := as.(type) {
			case basis.AuspiciousStar:
				starStr += v.String() + " "
			case basis.MaleficStar:
				starStr += v.String() + " "
			case basis.LuCunStar:
				starStr += v.String() + " "
			}
		}
		for _, ss := range secondaryStars {
			switch v := ss.(type) {
			case basis.SecondaryStar:
				starStr += v.String() + " "
			}
		}
		for _, ts := range transformedStars {
			switch v := ts.(type) {
			case basis.TransformedStar:
				starStr += v.String() + " "
			}
		}
		if starStr == "" {
			starStr = "(空宮)"
		}
		str += fmt.Sprintf("  %s(%s): %s\n", palace, branch, starStr)
	}

	str += "\n大運:\n"
	for _, dy := range c.DaYun {
		str += fmt.Sprintf("  第%d運(%d-%d歲): %s%s\n", dy.Index, dy.StartAge, dy.EndAge, dy.Stem, dy.Branch)
	}

	str += "\n流年:\n"
	for _, ln := range c.LiuNian {
		str += fmt.Sprintf("  %d年: %s%s\n", ln.Year, ln.Stem, ln.Branch)
	}

	str += "\n流月:\n"
	for _, ly := range c.LiuYue {
		str += fmt.Sprintf("  %d月: %s%s\n", ly.Month, ly.Stem, ly.Branch)
	}

	str += "\n流日:\n"
	for i := 0; i < min(15, len(c.LiuRi)); i++ {
		lr := c.LiuRi[i]
		str += fmt.Sprintf("  %d日: %s%s\n", lr.Day, lr.Stem, lr.Branch)
	}
	if len(c.LiuRi) > 15 {
		str += fmt.Sprintf("  ... (共%d日)\n", len(c.LiuRi))
	}

	return str
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
