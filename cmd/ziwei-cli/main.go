package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/kaecer68/ziwei-zenith/pkg/api/v1"
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
	"github.com/kaecer68/ziwei-zenith/pkg/engine"
)

var (
	year     int
	month    int
	day      int
	hour     int
	minute   int
	gender   string
	jsonFlag bool
)

func init() {
	flag.IntVar(&year, "year", 0, "Solar year (e.g., 1990)")
	flag.IntVar(&month, "month", 0, "Solar month (1-12)")
	flag.IntVar(&day, "day", 0, "Solar day (1-31)")
	flag.IntVar(&hour, "hour", 0, "Hour (0-23)")
	flag.IntVar(&minute, "minute", 0, "Minute (0-59)")
	flag.StringVar(&gender, "gender", "male", "Gender (male/female)")
	flag.BoolVar(&jsonFlag, "json", false, "Output JSON format")
}

func main() {
	flag.Parse()

	if year == 0 || month == 0 || day == 0 {
		fmt.Fprintf(os.Stderr, "Error: year, month, and day are required\n")
		flag.Usage()
		os.Exit(1)
	}

	sex := basis.SexMale
	if gender == "female" {
		sex = basis.SexFemale
	}

	yearPillar := basis.Pillar{Stem: basis.StemJia, Branch: basis.BranchZi}
	monthPillar := basis.Pillar{Stem: basis.StemYi, Branch: basis.BranchYin}
	dayPillar := basis.Pillar{Stem: basis.StemBing, Branch: basis.BranchWu}
	hourPillar := basis.Pillar{Stem: basis.StemDing, Branch: basis.BranchShen}

	birth := engine.BirthInfo{
		LunarYear:   year,
		LunarMonth:  month,
		LunarDay:    day,
		HourBranch:  basis.HourBranch(hour),
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

	if jsonFlag {
		outputJSON(chart)
	} else {
		fmt.Print(chart.String())
	}
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

	gender := "male"
	if chart.LifePalace.MingGong == basis.PalaceMing {
		gender = "female"
	}

	response := v1.ZiweiResponse{
		Gender:   gender,
		Wuxing:   chart.Wuxing.String(),
		NaYin:    chart.NaYin.String(),
		MingGong: chart.LifePalace.MingGong.String(),
		ShenGong: chart.LifePalace.ShenGong.String(),
		Palaces:  palaces,
	}

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(jsonBytes))
}
