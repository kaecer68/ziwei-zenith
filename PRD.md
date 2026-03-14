# Ziwei Zenith - 紫微斗數排盤引擎

## 1. Project Overview

**Project Name**: Ziwei Zenith  
**Project Type**: Go Library + CLI Tool  
**Core Functionality**: 紫微斗數 (Ziwei Douju) fate calculation engine - a complete Chinese astrology chart calculation system  
**Target Users**: Chinese astrology practitioners, developers building fortune-telling applications

## 2. Technical Architecture

### Technology Stack
- **Language**: Go 1.25+
- **Core Dependency**: lunar-zenith (for accurate stem-branch calculations)
- **Output Formats**: Text (CLI), JSON (API)
- **Logic Baseline**: **K-Master Engine v4.0** (Aligns with verified JavaScript implementation)

### Project Structure
```
ziwei-zenith/
├── cmd/ziwei-cli/main.go       # CLI tool
├── pkg/
│   ├── api/v1/types.go         # JSON API types
│   ├── basis/                  # Core data types & definitions
│   │   ├── definitions.go      # Stem, Branch, Element, Pillar
│   │   ├── stars.go            # 14 main stars (十四主星)
│   │   ├── palaces.go         # 12 palaces (十二宮)
│   │   ├── wuxing.go          # Five elements (五行), NaYin (納音)
│   │   ├── auspicious.go      # Auspicious stars (六吉星)
│   │   ├── malefic.go         # Malefic stars (六煞星)
│   │   ├── secondary.go       # Secondary stars (丙級星)
│   │   ├── transformation.go # Transformation stars (四化飛星)
│   │   ├── dayun.go           # DaYun/LiuNian/LiuYue/LiuRi types
│   │   ├── pattern.go         # Pattern definitions
│   │   └── brightness.go      # Star brightness (星曜亮度)
│   └── engine/
│       ├── engine.go          # Main ZiweiEngine
│       ├── lifepalace.go      # Life palace calculation
│       ├── starplacement.go   # Main star placement algorithm
│       ├── assistant.go       # Assistant + transformation stars
│       ├── dayun.go          # DaYun/LiuNian/LiuYue/LiuRi algorithms
│       ├── pattern.go        # Pattern detection
│       └── brightness.go     # Brightness calculation
└── go.mod
```

## 3. Implemented Features

### 3.1 Core Calculation
- ✅ **十四主星 (14 Main Stars)**: 紫微、天府、天機、太陽、武曲、天同、廉貞、巨門、天相、天梁、七殺、破軍、太陰、貪狼
- ✅ **六吉星 (6 Auspicious Stars)**: 左輔、右弼、文昌、文曲、天魁、天鉞
- ✅ **六煞星 (6 Malefic Stars)**: 擎羊、陀羅、火星、鈴星、地空、地劫
- ✅ **祿存、天馬 (LuCun & Tianma)**
- ✅ **丙級星 (12 Secondary Stars)**: 紅鸞、天喜、孤辰、寡宿、龍池、鳳閣、天刑、天姚、解神、天巫、三台、八座
- ✅ **四化飛星 (Transformation Stars)**: 化祿、化權、化科、化忌

### 3.2 Luck Cycles
- ✅ **大運 (DaYun)**: 10-year luck periods
- ✅ **流年 (LiuNian)**: Yearly luck (Fixed at Birth Year Branch)
- ✅ **流月 (LiuYue)**: Monthly luck (Dou Jun method)
- ✅ **流日 (LiuRi)**: Daily luck (Branch-step method)
- ✅ **疊併感應 (Resonance)**: Detecting clashing 'Ji' or 'Lu' between cycles.

### 3.3 Chart Analysis
- ✅ **命宮、身宮 (Life Palace & Body Palace)**
- ✅ **五行局 (Wuxing Ju)**: 金木水火土各局
- ✅ **納音五行 (NaYin)**
- ✅ **格局判斷 (Pattern Detection)**: 日月反背、祿馬交馳等
- ✅ **星曜亮度 (Star Brightness)**: 廟、旺、落、陷

### 3.4 Input/Output
- ✅ **Solar Date Calculation**: Convert Gregorian date to lunar date using lunar-zenith
- ✅ **Time System**: 支持時辰起盤 (12時辰制)
- ✅ **Gender Support**: 男命/女命
- ✅ **Text Output**: Human-readable console output
- ✅ **JSON Output**: `-json` flag for API integration

## 4. Algorithm Details

### 4.1 Life Palace (命宮) & Body Palace (身宮)
- **命宮**: 從寅宮起正月，順數生月，再從該宮起子時，逆數生時。
- **身宮**: 從寅宮起正月，順數生月，再從該宮起子時，順數生時。

- **紫微星定位**: 根據五行局與農曆生日，沿用 Master Engine 的 Parity 算法 (Odd Forward/Backward 邏輯)。
- **來因宮 (Origin Palace)**: 採用欽天門算法，定位生年天干所在宮位。
- **閏月處理**: 遵循 15 日分隔線法則 (Master Engine 規範)。
- **主星分佈**: 
  - 紫微星群：紫微、天機、太陽、武曲、天同、廉貞（逆時鐘分佈）。
  - 天府星群：天府、太陰、貪狼、巨門、天相、天梁、七殺、破軍（順時鐘分佈）。
  - 對稱規則：tf = (4 - zw + 12) % 12。

- **安四化規則**:
  - 甲：廉破武陽 (廉貞化祿、破軍化權、武曲化科、太陽化忌)
  - 乙：機梁紫陰 (天機化祿、天梁化權、紫微化科、太陰化忌)
  - 丙：同機昌廉 (天同化祿、天機化權、文昌化科、廉貞化忌)
  - 丁：陰同機巨 (太陰化祿、天同化權、天機化科、巨門化忌)
  - 戊：貪陰右機 (貪狼化祿、太陰化權、右弼化科、天機化忌)
  - 己：武貪梁曲 (武曲化祿、貪狼化權、天梁化科、文曲化忌)
  - 庚：陽武府同 (太陽化祿、武曲化權、天府化科、天同化忌)
  - 辛：巨陽曲昌 (巨門化祿、太陽化權、文曲化科、文昌化忌)
  - 壬：梁紫左武 (天梁化祿、紫微化權、左輔化科、武曲化忌)
  - 癸：破巨陰貪 (破軍化祿、巨門化權、太陰化科、貪狼化忌)

- **大運起點**: 五行局數（如水二局則自2歲起大運）。
- **行運方向**: 陽男陰女順行，陰男陽女逆行。
- **流年**: 固定在地盤各宮位（如甲子年流年命宮在子）。
- **流月**: 採用「斗君起月」法。自流年命宮起正月，逆數生月，再順數生時（正月點），隨後依農曆月份順行。
- **時空感應**: 自動檢測「歲運疊併」。當流年/月/日三層四化與本命四化產生重疊（如雙重化忌）時，觸發警告邏輯。

### 4.5 Star Brightness (6-Level System)
| 等級 | Description |
|------|-------------|
| 廟 (Ming) | Highest influence (宮位與星曜五行極其契合) |
| 旺 (Wang) | Strong influence |
| 利 (Li) | Moderate influence |
| 平 (Ping) | Neutral influence |
| 陷 (Xian) | Low influence (能量受挫) |
| 不 (Bu/None)| Weakest/Disharmonious |

## 5. Usage

### CLI
```bash
# Basic usage
ziwei-cli -year 1990 -month 6 -day 15 -hour 10

# With gender
ziwei-cli -year 1990 -month 6 -day 15 -hour 10 -gender female

# JSON output
ziwei-cli -year 1990 -month 6 -day 15 -hour 10 -json

# Custom location
ziwei-cli -year 1990 -month 6 -day 15 -hour 10 -lat 25.033 -lon 121.565
```

### Library
```go
package main

import (
    "fmt"
    "github.com/kaecer68/ziwei-zenith/pkg/engine"
)

func main() {
    e := engine.New()
    chart, _ := e.BuildChart(engine.BirthInfo{
        SolarYear:  1990,
        SolarMonth: 6,
        SolarDay:   15,
        Hour:       10,
        Gender:     "male",
    })
    fmt.Println(chart)
}
```

## 6. Data Models

### Core Types
- `Star`: 14 main stars + assistant + secondary + transformed
- `Palace`: 12 palaces (命宮、兄弟宮、夫妻宮、子女宮、財帛宮、疾厄宮、遷移宮、僕役宮、官祿宮、田宅宮、父母宮、福德宮)
- `Branch`: 12 earthly branches (子丑寅卯辰巳午未申酉戌亥)
- `Stem`: 10 heavenly stems (甲乙丙丁戊己庚辛壬癸)
- `Wuxing`: 五行 (木火土金水)
- `NaYin`: 納音五行

### Output Structure
```json
{
  "yearPillar": "庚午",
  "monthPillar": "癸未", 
  "dayPillar": "辛亥",
  "hourPillar": "癸巳",
  "wuxingJu": "水二局",
  "naYin": "己亥",
  "lifePalace": "田宅宮",
  "bodyPalace": "僕役宮",
  "palaces": {...},
  "stars": {...},
  "patterns": [...],
  "starBrightness": [...],
  "dayun": [...],
  "liunian": [...],
  "liuyue": [...],
  "liuri": [...]
}
```

## 7. Future Enhancements

### Phase 2 (Planned)
- [ ] 小限 (Xiao Xian) calculation
- [ ] 鐵板神數 (Tieban Shenshu) integration
- [ ] 更多格局判斷
- [ ] 星曜組合分析
- [ ] Web API server

### Phase 3 (Roadmap)
- [ ] GUI Application
- [ ] Mobile SDK
- [ ] Birth chart visualization

## 8. Dependencies

- **lunar-zenith**: High-precision Chinese lunar calendar library
  - Provides accurate 天干地支 (Stem-Branch) calculations
  - Handles 閏月, 節氣, and exact solar-lunar conversion

## 9. References

- 《紫微斗數全書》
- 《紫微斗數命理學》
- 《現代紫微斗數論命寶典》

## 10. License

MIT License
