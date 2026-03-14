package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"time"

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
	latitude  float64
	longitude float64
)

func init() {
	flag.IntVar(&year, "year", 0, "Solar year (e.g., 1990)")
	flag.IntVar(&month, "month", 0, "Solar month (1-12)")
	flag.IntVar(&day, "day", 0, "Solar day (1-31)")
	flag.IntVar(&hour, "hour", 0, "Hour (0-23)")
	flag.IntVar(&minute, "minute", 0, "Minute (0-59)")
	flag.StringVar(&gender, "gender", "male", "Gender (male/female)")
	flag.Float64Var(&latitude, "lat", 25.033, "Latitude (default: 25.033 - Taipei)")
	flag.Float64Var(&longitude, "lon", 121.565, "Longitude (default: 121.565 - Taipei)")
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

	yearSexagenary := zodiac.NewYearSexagenary(year)

	monthSexagenary := zodiac.GetMonthSexagenary(yearSexagenary.StemIndex, int(month))

	jd := julianDay(solarTime)
	daySexagenary := zodiac.GetDaySexagenary(jd)

	hourBranchIdx := zodiac.GetHourBranch(hour)
	hourSexagenary := zodiac.GetHourSexagenary(daySexagenary.StemIndex, hourBranchIdx)

	lunarMonth := calcLunarMonth(year, month, day, jd)
	lunarDay := calcLunarDay(jd)

	yearPillar := basis.Pillar{
		Stem:   basis.Stem(yearSexagenary.StemIndex),
		Branch: basis.Branch(yearSexagenary.BranchIndex),
	}
	monthPillar := basis.Pillar{
		Stem:   basis.Stem(monthSexagenary.StemIndex),
		Branch: basis.Branch(monthSexagenary.BranchIndex),
	}
	dayPillar := basis.Pillar{
		Stem:   basis.Stem(daySexagenary.StemIndex),
		Branch: basis.Branch(daySexagenary.BranchIndex),
	}
	hourPillar := basis.Pillar{
		Stem:   basis.Stem(hourSexagenary.StemIndex),
		Branch: basis.Branch(hourSexagenary.BranchIndex),
	}

	birth := engine.BirthInfo{
		LunarYear:   lunarMonth,
		LunarMonth:  int(math.Abs(float64(lunarMonth))),
		LunarDay:    lunarDay,
		HourBranch:  basis.HourBranch(hourBranchIdx),
		YearPillar:  yearPillar,
		MonthPillar: monthPillar,
		DayPillar:   dayPillar,
		HourPillar:  hourPillar,
		Sex:         sex,
	}

	e := engine.New()
	chart, err := e.BuildChart(birth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	chart.YearPillar = yearPillar
	chart.MonthPillar = monthPillar
	chart.DayPillar = dayPillar
	chart.HourPillar = hourPillar

	if jsonFlag {
		outputJSON(chart)
	} else {
		fmt.Print(chart.String())
	}
}

func julianDay(t time.Time) float64 {
	y := t.Year()
	m := t.Month()
	d := float64(t.Day())
	h := float64(t.Hour()) + float64(t.Minute())/60.0

	if m <= 2 {
		y -= 1
		m += 12
	}

	a := math.Floor(float64(y) / 100)
	b := 2 - a + math.Floor(a/4)

	return math.Floor(365.25*float64(y+4716)) + math.Floor(30.6001*float64(m+1)) + d + h/24.0 + b - 1524.5
}

func calcLunarMonth(year, month, day int, jd float64) int {
	newMoonJD := findPreviousNewMoon(jd)
	refDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	refJD := julianDay(refDate)

	months := int((refJD - newMoonJD) / 29.53)
	return (months % 12) + 1
}

func findPreviousNewMoon(jd float64) float64 {
	low := jd - 30
	high := jd

	for i := 0; i < 20; i++ {
		mid := (low + high) / 2
		phase := moonPhase(mid)
		if phase > 180 {
			phase -= 360
		}
		if math.Abs(high-low) < 0.01 {
			return mid
		}
		if phase < 0 {
			low = mid
		} else {
			high = mid
		}
	}
	return (low + high) / 2
}

func moonPhase(jde float64) float64 {
	t := (jde - 2451545.0) / 36525.0

	lPrime := 218.3164477 + 481267.88123421*t
	mPrime := 134.9633964 + 477198.8675055*t
	m := 357.5291092 + 35999.0502909*t
	f := 93.2720950 + 483202.0175233*t
	d := 297.8501921 + 445267.1114034*t

	lPrime = math.Mod(lPrime, 360.0)
	if lPrime < 0 {
		lPrime += 360.0
	}

	deg2Rad := math.Pi / 180.0
	lambda := lPrime +
		6.288774*math.Sin(mPrime*deg2Rad) +
		1.274027*math.Sin((2*d-mPrime)*deg2Rad) +
		0.658314*math.Sin(2*d*deg2Rad) +
		0.213118*math.Sin(2*mPrime*deg2Rad) -
		0.185116*math.Sin(m*deg2Rad) -
		0.114332*math.Sin(2*f*deg2Rad)

	sLon := solarLongitude(jde)
	mLon := math.Mod(lambda+360.0, 360.0)
	diff := mLon - sLon
	return math.Mod(diff+360.0, 360.0)
}

func solarLongitude(jde float64) float64 {
	t := (jde - 2451545.0) / 36525.0
	l0 := 280.46646 + 36000.76983*t + 0.0003032*t*t
	return math.Mod(l0, 360.0)
}

func calcLunarDay(jd float64) int {
	phase := moonPhase(jd)
	day := int((phase / 12.0)) + 1
	if day > 30 {
		day = 30
	}
	return day
}

func outputJSON(chart *engine.ZiweiChart) {
	palaces := make(map[string]v1.PalaceData)
	for i := 0; i < 12; i++ {
		palace := basis.Palace(i)
		branch := chart.Palaces[palace]
		stars := chart.Stars[palace]
		starNames := make([]string, len(stars))
		for j, s := range stars {
			starNames[j] = s.String()
		}
		palaces[palace.String()] = v1.PalaceData{
			Branch: branch.String(),
			Stars:  starNames,
		}
	}

	genderStr := "male"
	if chart.LifePalace.MingGong == basis.PalaceMing {
		genderStr = "female"
	}

	response := v1.ZiweiResponse{
		Gender:     genderStr,
		Wuxing:     chart.Wuxing.String(),
		NaYin:      chart.NaYin.String(),
		MingGong:   chart.LifePalace.MingGong.String(),
		ShenGong:   chart.LifePalace.ShenGong.String(),
		YearPillar: chart.YearPillar.String(),
		DayPillar:  chart.DayPillar.String(),
		Palaces:    palaces,
	}

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonBytes))
}
