// Package basis 提供紫微斗數核心數據模型和算法
package basis

import (
	"fmt"
)

// 地支（十二支）
type Branch int

const (
	BranchZi   Branch = iota // 子
	BranchChou               // 丑
	BranchYin                // 寅
	BranchMao                // 卯
	BranchChen               // 辰
	BranchSi                 // 巳
	BranchWu                 // 午
	BranchWei                // 未
	BranchShen               // 申
	BranchYou                // 酉
	BranchXu                 // 戌
	BranchHai                // 亥
)

func (b Branch) String() string {
	names := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	return names[b]
}

// BranchByName 根據名稱獲取地支
func BranchByName(name string) (Branch, error) {
	names := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	for i, n := range names {
		if n == name {
			return Branch(i), nil
		}
	}
	return 0, fmt.Errorf("未知的地支: %s", name)
}

// 天干（十干）
type Stem int

const (
	StemJia  Stem = iota // 甲
	StemYi               // 乙
	StemBing             // 丙
	StemDing             // 丁
	StemWu               // 戊
	StemJi               // 己
	StemGeng             // 庚
	StemXin              // 辛
	StemRen              // 壬
	StemGui              // 癸
)

func (s Stem) String() string {
	names := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	return names[s]
}

// StemByName 根據名稱獲取天干
func StemByName(name string) (Stem, error) {
	names := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	for i, n := range names {
		if n == name {
			return Stem(i), nil
		}
	}
	return 0, fmt.Errorf("未知的天干: %s", name)
}

// 五行
type Element int

const (
	ElementWood  Element = iota // 木
	ElementFire                 // 火
	ElementEarth                // 土
	ElementMetal                // 金
	ElementWater                // 水
)

func (e Element) String() string {
	names := []string{"木", "火", "土", "金", "水"}
	return names[e]
}

// 陰陽
type Polarity int

const (
	PolarityYang Polarity = iota // 陽
	PolarityYin                  // 陰
)

func (p Polarity) String() string {
	if p == PolarityYang {
		return "陽"
	}
	return "陰"
}

// 干支結構（年、月、日、時柱）
type Pillar struct {
	Stem   Stem
	Branch Branch
}

// Sex represents gender for fortune calculation
type Sex int

const (
	SexMale   Sex = iota // 男
	SexFemale            // 女
)

func (s Sex) String() string {
	if s == SexMale {
		return "男"
	}
	return "女"
}

// 時辰（12個時辰）
type HourBranch int

const (
	HourZi   HourBranch = iota // 子時 (23:00-00:59)
	HourChou                   // 丑時 (01:00-02:59)
	HourYin                    // 寅時 (03:00-04:59)
	HourMao                    // 卯時 (05:00-06:59)
	HourChen                   // 辰時 (07:00-08:59)
	HourSi                     // 巳時 (09:00-10:59)
	HourWu                     // 午時 (11:00-12:59)
	HourWei                    // 未時 (13:00-14:59)
	HourShen                   // 申時 (15:00-16:59)
	HourYou                    // 酉時 (17:00-18:59)
	HourXu                     // 戌時 (19:00-20:59)
	HourHai                    // 亥時 (21:00-22:59)
)

func (h HourBranch) String() string {
	names := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	return names[h]
}

// HourBranchFromTime converts hour to HourBranch
func HourBranchFromTime(hour int) HourBranch {
	// Convert 24-hour to zi-hour (子時 = 23:00-00:59)
	// Formula: (hour + 1) / 2 mod 12
	return HourBranch(((hour + 1) / 2) % 12)
}

// 農曆日期
type LunarDate struct {
	Year  int // 農曆年
	Month int // 農曆月（正月初一為正月）
	Day   int // 農曆日
}

// BirthInfo 出生資訊
type BirthInfo struct {
	SolarYear  int // 公曆年
	SolarMonth int // 公曆月
	SolarDay   int // 公曆日
	Hour       int // 小時 (0-23)
	Minute     int // 分鐘 (0-59)
	Sex        Sex // 性別
	// 以下由 lunar-zenith 轉換產生
	LunarYear   int        // 農曆年
	LunarMonth  int        // 農曆月（1-12，閏月為負數）
	LunarDay    int        // 農曆日
	HourBranch  HourBranch // 時辰地支
	YearPillar  Pillar     // 年柱
	MonthPillar Pillar     // 月柱
	DayPillar   Pillar     // 日柱
	HourPillar  Pillar     // 時柱
}
