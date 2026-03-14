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

	return &ZiweiChart{
		LifePalace:       lifePalace,
		Palaces:          palaces,
		Stars:            stars,
		AssistantStars:   assistantStars,
		SecondaryStars:   secondaryStars,
		TransformedStars: transformedStars,
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
	return str
}
