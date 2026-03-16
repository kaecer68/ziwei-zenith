package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
	"github.com/kaecer68/ziwei-zenith/pkg/api/v1"
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
	"github.com/kaecer68/ziwei-zenith/pkg/engine"
)

var (
	year      int
	month     int
	day       int
	hour      int
	minute    int
	gender    string
	jsonFlag  bool
	isLunar   bool
	latitude  float64
	longitude float64
)

func init() {
	flag.IntVar(&year, "year", 0, "Year")
	flag.IntVar(&month, "month", 0, "Month")
	flag.IntVar(&day, "day", 0, "Day")
	flag.IntVar(&hour, "hour", 0, "Hour")
	flag.IntVar(&minute, "minute", 0, "Minute")
	flag.StringVar(&gender, "gender", "male", "Gender (male/female)")
	flag.BoolVar(&isLunar, "lunar", false, "Use lunar date input directly")
	flag.Float64Var(&latitude, "lat", 25.033, "Latitude")
	flag.Float64Var(&longitude, "lon", 121.565, "Longitude")
	flag.BoolVar(&jsonFlag, "json", false, "Output JSON format")
}

func main() {
	flag.Parse()

	if year == 0 || month == 0 || day == 0 {
		fmt.Fprintf(os.Stderr, "Error: year, month, and day are required\n")
		flag.Usage()
		os.Exit(1)
	}

	loc := time.FixedZone("", int((longitude/15)*3600))
	solarTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc)

	sex := basis.SexMale
	if gender == "female" {
		sex = basis.SexFemale
	}

	var lYear, lMonth, lDay int
	var yPillar, mPillar, dPillar basis.Pillar

	pt := celestial.NewPrecisionTime(solarTime)
	pillar := zodiac.GetAstrologicalPillar(pt)

	yPillar = basis.Pillar{Stem: basis.Stem(pillar.Year.StemIndex), Branch: basis.Branch(pillar.Year.BranchIndex)}
	mPillar = basis.Pillar{Stem: basis.Stem(pillar.Month.StemIndex), Branch: basis.Branch(pillar.Month.BranchIndex)}
	// Fix: use local time JD for day pillar (library converts to UTC internally)
	localDayPillar := zodiac.GetDaySexagenary(celestial.TimeToJD(solarTime))
	dPillar = basis.Pillar{Stem: basis.Stem(localDayPillar.StemIndex), Branch: basis.Branch(localDayPillar.BranchIndex)}

	if isLunar {
		lYear = year
		lMonth = month
		lDay = day
	} else {
		jd := celestial.TimeToJD(solarTime)
		engine_lunar := &zodiac.LunarEngine{}
		lunarDate := engine_lunar.GetLunarDate(jd)

		lYear = lunarDate.Year
		lMonth = lunarDate.Month
		lDay = lunarDate.Day
	}

	hourBranchIdx := zodiac.GetHourBranch(hour)
	hourSexagenary := zodiac.GetHourSexagenary(int(dPillar.Stem), hourBranchIdx)

	birth := basis.BirthInfo{
		SolarYear:   year,
		SolarMonth:  month,
		SolarDay:    day,
		Hour:        hour,
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if jsonFlag {
		outputJSON(chart, gender)
	} else {
		fmt.Print(chart.String())
	}
}

func outputJSON(chart *engine.ZiweiChart, gender string) {
	palaces := make(map[string]v1.PalaceData)
	for i := 0; i < 12; i++ {
		b := basis.Branch(i)
		pType := chart.Palaces[b]
		stars := chart.Stars[b]
		starNames := make([]string, 0)
		for _, s := range stars {
			starNames = append(starNames, s.String())
		}
		lnStars := make([]string, 0)
		for _, s := range chart.LiuNianStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				lnStars = append(lnStars, strer.String())
			}
		}

		lyStars := make([]string, 0)
		for _, s := range chart.LiuYueStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				lyStars = append(lyStars, strer.String())
			}
		}

		lrStars := make([]string, 0)
		for _, s := range chart.LiuRiStars[b] {
			if strer, ok := s.(interface{ String() string }); ok {
				lrStars = append(lrStars, strer.String())
			}
		}

		palaces[pType.String()] = v1.PalaceData{
			Branch:       b.String(),
			Stars:        starNames,
			LiuNianStars: lnStars,
			LiuYueStars:  lyStars,
			LiuRiStars:   lrStars,
		}
	}

	patterns := make([]v1.PatternData, 0)
	for _, p := range chart.Patterns {
		patterns = append(patterns, v1.PatternData{
			Name:        p.Name,
			Description: p.Description,
			Level:       p.Level,
		})
	}

	narrative := make([]v1.KarmicStep, 0)
	for _, s := range chart.Interpretation.KarmicNarrative {
		narrative = append(narrative, v1.KarmicStep{
			Type:   s.Type,
			Role:   s.Role,
			Star:   s.Star,
			Palace: s.Palace,
			Desc:   s.Desc,
		})
	}

	diagnosis := make([]v1.SanFangRole, 0)
	for _, r := range chart.Interpretation.SanFangDiagnosis {
		diagnosis = append(diagnosis, v1.SanFangRole{
			Role:      r.Role,
			Palace:    r.Palace,
			Diagnosis: r.Diagnosis,
		})
	}

	resonance := make([]v1.ResonancePoint, 0)
	for _, r := range chart.Interpretation.TemporalResonance {
		resonance = append(resonance, v1.ResonancePoint{
			Layer:  r.Layer,
			Type:   r.Type,
			Star:   r.Star,
			Natal:  r.Natal,
			Palace: r.Palace,
			Mood:   r.Mood,
		})
	}

	response := v1.ZiweiResponse{
		Gender:       gender,
		Wuxing:       chart.Wuxing.String(),
		NaYin:        chart.NaYin.String(),
		OriginPalace: chart.Palaces[chart.OriginPalace].String(),
		MingGong:     chart.LifePalace.MingGong.String(),
		ShenGong:     chart.LifePalace.ShenGong.String(),
		YearPillar:   chart.YearPillar.String(),
		DayPillar:    chart.DayPillar.String(),
		Palaces:      palaces,
		Patterns:     patterns,
		Interpretation: v1.InterpretationData{
			Summary:           chart.Interpretation.Summary,
			KarmicNarrative:   narrative,
			SanFangDiagnosis:  diagnosis,
			TemporalResonance: resonance,
		},
	}

	jsonBytes, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonBytes))
}
