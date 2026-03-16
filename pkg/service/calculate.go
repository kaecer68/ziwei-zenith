package service

import (
	"fmt"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
	"github.com/kaecer68/ziwei-zenith/pkg/engine"
)

// CalculateInput 排盤計算的統一輸入結構
type CalculateInput struct {
	Year      int
	Month     int
	Day       int
	Hour      int
	Minute    int
	Gender    string
	IsLunar   bool
	IsLeap    bool
	IsDST     bool
	Longitude float64
}

// Calculate 執行紫微斗數排盤計算，回傳 ZiweiChart
func Calculate(input CalculateInput) (*engine.ZiweiChart, error) {
	calcHour := input.Hour
	if input.IsDST {
		calcHour--
	}

	lon := input.Longitude
	if lon == 0 {
		lon = 121.565 // 預設台北
	}

	loc := time.FixedZone("", int((lon/15)*3600))
	solarTime := time.Date(input.Year, time.Month(input.Month), input.Day, calcHour, input.Minute, 0, 0, loc)

	sex := basis.SexMale
	if input.Gender == "female" {
		sex = basis.SexFemale
	}

	pt := celestial.NewPrecisionTime(solarTime)
	pillar := zodiac.GetAstrologicalPillar(pt)

	yPillar := basis.Pillar{Stem: basis.Stem(pillar.Year.StemIndex), Branch: basis.Branch(pillar.Year.BranchIndex)}
	mPillar := basis.Pillar{Stem: basis.Stem(pillar.Month.StemIndex), Branch: basis.Branch(pillar.Month.BranchIndex)}
	dPillar := basis.Pillar{Stem: basis.Stem(pillar.Day.StemIndex), Branch: basis.Branch(pillar.Day.BranchIndex)}

	// 使用本地時間 JD 計算日柱（避免 UTC 偏移問題）
	localDayPillar := zodiac.GetDaySexagenary(celestial.TimeToJD(solarTime))
	dPillar = basis.Pillar{Stem: basis.Stem(localDayPillar.StemIndex), Branch: basis.Branch(localDayPillar.BranchIndex)}

	var lYear, lMonth, lDay int
	if input.IsLunar {
		lYear = input.Year
		lMonth = input.Month
		if input.IsLeap {
			lMonth = -input.Month
		}
		lDay = input.Day
	} else {
		jd := celestial.TimeToJD(solarTime)
		lunarEngine := &zodiac.LunarEngine{}
		lunarDate := lunarEngine.GetLunarDate(jd)
		lYear = lunarDate.Year
		lMonth = lunarDate.Month
		lDay = lunarDate.Day
	}

	hourBranchIdx := zodiac.GetHourBranch(input.Hour)
	hourSexagenary := zodiac.GetHourSexagenary(int(dPillar.Stem), hourBranchIdx)

	birth := basis.BirthInfo{
		SolarYear:   input.Year,
		SolarMonth:  input.Month,
		SolarDay:    input.Day,
		Hour:        input.Hour,
		Sex:         sex,
		LunarYear:   lYear,
		LunarMonth:  lMonth,
		LunarDay:    lDay,
		HourBranch:  basis.HourBranch(hourBranchIdx),
		YearPillar:  yPillar,
		MonthPillar: mPillar,
		DayPillar:   dPillar,
		HourPillar:  basis.Pillar{Stem: basis.Stem(hourSexagenary.StemIndex), Branch: basis.Branch(hourSexagenary.BranchIndex)},
	}

	e := engine.New()
	chart, err := e.BuildChart(birth)
	if err != nil {
		return nil, fmt.Errorf("calculation error: %w", err)
	}
	return chart, nil
}
