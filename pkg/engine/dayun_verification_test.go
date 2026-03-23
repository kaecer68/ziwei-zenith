package engine

import (
	"testing"

	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

// TestWuxingValues 驗證五行局數值與起運歲數對應關系
func TestWuxingValues(t *testing.T) {
	testCases := []struct {
		wuxing       basis.Wuxing
		expectedVal  int
		expectedName string
	}{
		{basis.WuxingShui2, 2, "水二局"},
		{basis.WuxingMu3, 3, "木三局"},
		{basis.WuxingJin4, 4, "金四局"},
		{basis.WuxingTu5, 5, "土五局"},
		{basis.WuxingHuo6, 6, "火六局"},
	}

	for _, tc := range testCases {
		if tc.wuxing.Value() != tc.expectedVal {
			t.Errorf("%s: 期望數值 %d, 實際 %d", tc.expectedName, tc.expectedVal, tc.wuxing.Value())
		}
		if tc.wuxing.String() != tc.expectedName {
			t.Errorf("期望名稱 %s, 實際 %s", tc.expectedName, tc.wuxing.String())
		}
	}
}

// TestDaYunDirection 驗證大限順逆方向 (陽男陰女順行, 陰男陽女逆行)
func TestDaYunDirection(t *testing.T) {
	mingBranch := basis.BranchChen // 命宮在辰(4)
	wuxing := basis.WuxingShui2    // 水二局

	testCases := []struct {
		name         string
		yearStem     basis.Stem
		sex          basis.Sex
		expectedDir  int // 1=順行, -1=逆行
		secondBranch basis.Branch
	}{
		{"陽男-甲年順行", basis.StemJia, basis.SexMale, 1, basis.BranchSi},     // 甲(0)=陽, 男, 順行→巳(5)
		{"陽女-甲年逆行", basis.StemJia, basis.SexFemale, -1, basis.BranchMao}, // 甲(0)=陽, 女, 逆行→卯(3)
		{"陰男-乙年逆行", basis.StemYi, basis.SexMale, -1, basis.BranchMao},    // 乙(1)=陰, 男, 逆行→卯(3)
		{"陰女-乙年順行", basis.StemYi, basis.SexFemale, 1, basis.BranchSi},    // 乙(1)=陰, 女, 順行→巳(5)
		{"陽男-丙年順行", basis.StemBing, basis.SexMale, 1, basis.BranchSi},    // 丙(2)=陽
		{"陰女-癸年順行", basis.StemGui, basis.SexFemale, 1, basis.BranchSi},   // 癸(9)=陰
	}

	for _, tc := range testCases {
		dayuns := CalcDaYun(mingBranch, tc.yearStem, tc.sex, wuxing)

		// 驗證第二個大限的宮位
		if dayuns[1].Branch != tc.secondBranch {
			t.Errorf("%s: 期望第二限在 %v, 實際在 %v",
				tc.name, tc.secondBranch, dayuns[1].Branch)
		}
	}
}

// TestDaYunAges 驗證各五行局的起運歲數和10年區間
func TestDaYunAges(t *testing.T) {
	mingBranch := basis.BranchChen
	yearStem := basis.StemJia // 陽年
	sex := basis.SexMale      // 陽男=順行

	testCases := []struct {
		wuxing         basis.Wuxing
		firstStart     int
		firstEnd       int
		secondStart    int
		expectedPalace []basis.Branch // 預期的十二宮順序(順行)
	}{
		{basis.WuxingShui2, 2, 11, 12, []basis.Branch{4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3}}, // 水二局: 2-11, 12-21...
		{basis.WuxingMu3, 3, 12, 13, []basis.Branch{4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3}},   // 木三局: 3-12, 13-22...
		{basis.WuxingJin4, 4, 13, 14, []basis.Branch{4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3}},  // 金四局: 4-13, 14-23...
		{basis.WuxingTu5, 5, 14, 15, []basis.Branch{4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3}},   // 土五局: 5-14, 15-24...
		{basis.WuxingHuo6, 6, 15, 16, []basis.Branch{4, 5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3}},  // 火六局: 6-15, 16-25...
	}

	for _, tc := range testCases {
		dayuns := CalcDaYun(mingBranch, yearStem, sex, tc.wuxing)

		// 驗證第一個大限的歲數
		if dayuns[0].StartAge != tc.firstStart || dayuns[0].EndAge != tc.firstEnd {
			t.Errorf("%s: 第一限期望 %d-%d歲, 實際 %d-%d歲",
				tc.wuxing.String(), tc.firstStart, tc.firstEnd,
				dayuns[0].StartAge, dayuns[0].EndAge)
		}

		// 驗證十二個大限的宮位順序
		for i := 0; i < 12; i++ {
			if dayuns[i].Branch != tc.expectedPalace[i] {
				t.Errorf("%s: 第%d限期望在 %v, 實際在 %v",
					tc.wuxing.String(), i+1, tc.expectedPalace[i], dayuns[i].Branch)
				break
			}
		}
	}
}

// TestCurrentDaYun 驗證當前大限計算 (根據虛歲匹配)
func TestCurrentDaYun(t *testing.T) {
	// 測試數據: 1990年出生, 土五局, 陽男順行
	// 虛歲 = 當前年份 - 出生年份 + 1
	birthYear := 1990
	wuxing := basis.WuxingTu5 // 土五局, 5歲起運
	mingBranch := basis.BranchChen
	yearStem := basis.StemGeng // 庚=陽
	sex := basis.SexMale

	dayuns := CalcDaYun(mingBranch, yearStem, sex, wuxing)

	testCases := []struct {
		currentYear int
		expectedIdx int    // 期望的大限索引 (0-based)
		expectedAge string // 期望的歲數範圍
	}{
		{1994, 0, "5-14"},  // 虛歲5, 第一限
		{2000, 0, "5-14"},  // 虛歲11, 第一限
		{2005, 1, "15-24"}, // 虛歲16, 第二限
		{2010, 1, "15-24"}, // 虛歲21, 第二限
		{2020, 2, "25-34"}, // 虛歲31, 第三限
		{2026, 3, "35-44"}, // 虛歲37, 第四限
	}

	for _, tc := range testCases {
		currentAge := tc.currentYear - birthYear + 1 // 虛歲計算

		// 查找對應的大限
		var found *basis.DaYun
		for i := range dayuns {
			if currentAge >= dayuns[i].StartAge && currentAge <= dayuns[i].EndAge {
				found = &dayuns[i]
				break
			}
		}

		if found == nil {
			t.Errorf("%d年(虛歲%d): 未找到對應大限", tc.currentYear, currentAge)
			continue
		}

		// 驗證大限索引
		expectedDaYun := dayuns[tc.expectedIdx]
		if found.Index != expectedDaYun.Index {
			t.Errorf("%d年(虛歲%d): 期望第%d限(%s), 實際第%d限(%d-%d)",
				tc.currentYear, currentAge, tc.expectedIdx+1, tc.expectedAge,
				found.Index, found.StartAge, found.EndAge)
		}
	}
}

// TestLiuNian 驗證流年計算 (太歲地支)
func TestLiuNian(t *testing.T) {
	testCases := []struct {
		year       int
		yearBranch basis.Branch
		expected   basis.Branch
	}{
		{2024, basis.BranchChen, basis.BranchChen}, // 甲辰年, 辰宮
		{2025, basis.BranchSi, basis.BranchSi},     // 乙巳年, 巳宮
		{2026, basis.BranchWu, basis.BranchWu},     // 丙午年, 午宮
		{2027, basis.BranchWei, basis.BranchWei},   // 丁未年, 未宮
	}

	for _, tc := range testCases {
		ln := CalcLiuNian(tc.yearBranch, tc.year)
		if ln.Branch != tc.expected {
			t.Errorf("%d年: 期望流年命宮在 %v, 實際在 %v",
				tc.year, tc.expected, ln.Branch)
		}
	}
}

// TestLiuYueDouJun 驗證流月斗君計算
func TestLiuYueDouJun(t *testing.T) {
	// 測試案例: 2024甲辰年(辰=4), 農曆4月出生, 丑時(1)
	// 公式: (ln - (bM-1) + bH + 12) % 12
	// 正月: (4 - 3 + 1 + 12) % 12 = 2 (寅宮)
	lnBranch := basis.BranchChen
	birthMonth := 4
	birthHour := basis.BranchChou // 丑=1

	// 驗證正月位置
	month1 := CalcLiuYue(lnBranch, birthMonth, birthHour, 1)
	if month1 != basis.BranchYin {
		t.Errorf("正月: 期望在寅(2), 實際在 %v(%d)", month1, int(month1))
	}

	// 驗證各月位置 (順行)
	expectedMonths := []basis.Branch{
		basis.BranchYin,  // 正月
		basis.BranchMao,  // 二月
		basis.BranchChen, // 三月
		basis.BranchSi,   // 四月
		basis.BranchWu,   // 五月
		basis.BranchWei,  // 六月
		basis.BranchShen, // 七月
		basis.BranchYou,  // 八月
		basis.BranchXu,   // 九月
		basis.BranchHai,  // 十月
		basis.BranchZi,   // 十一月
		basis.BranchChou, // 十二月
	}

	for i, expected := range expectedMonths {
		month := CalcLiuYue(lnBranch, birthMonth, birthHour, i+1)
		if month != expected {
			t.Errorf("農曆%d月: 期望在 %v, 實際在 %v", i+1, expected, month)
		}
	}
}

// TestLiuRi 驗證流日計算 (從流月順行)
func TestLiuRi(t *testing.T) {
	// 測試: 流月在辰宮(4)
	lyBranch := basis.BranchChen

	testCases := []struct {
		day      int
		expected basis.Branch
	}{
		{1, basis.BranchChen},  // 初一: 辰
		{2, basis.BranchSi},    // 初二: 巳
		{3, basis.BranchWu},    // 初三: 午
		{13, basis.BranchChen}, // 十三: 回到辰 (12天一周)
	}

	for _, tc := range testCases {
		result := CalcLiuRi(lyBranch, tc.day)
		if result != tc.expected {
			t.Errorf("流月%s的第%d天: 期望在 %v, 實際在 %v",
				lyBranch, tc.day, tc.expected, result)
		}
	}
}

// TestTemporalChain 驗證流年-流月-流日鏈式計算
func TestTemporalChain(t *testing.T) {
	// 測試場景: 2024年3月15日 (農曆)
	// 需要模擬從公曆轉農曆的計算

	// 簡化測試: 直接使用農曆數據
	lnBranch := basis.BranchChen // 2024甲辰年
	birthMonth := 6
	birthHour := basis.BranchWu // 午時
	targetLunarMonth := 3       // 農曆三月
	targetLunarDay := 15        // 農曆十五

	// 計算流月
	liuYue := CalcLiuYue(lnBranch, birthMonth, birthHour, targetLunarMonth)

	// 計算流日
	liuRi := CalcLiuRi(liuYue, targetLunarDay)

	t.Logf("測試鏈式計算: 流年%s -> 流月%s -> 流日%s",
		lnBranch, liuYue, liuRi)

	// 基本驗證: 流日應該是從流月順行 (day-1) 個宮位
	expectedRiIdx := (int(liuYue) + (targetLunarDay - 1)) % 12
	if int(liuRi) != expectedRiIdx {
		t.Errorf("流日計算不一致: 期望索引%d, 實際%d", expectedRiIdx, int(liuRi))
	}
}

// GetCurrentDaYun 輔助函數: 根據虛歲查找當前大限
func GetCurrentDaYun(dayuns []basis.DaYun, birthYear int, currentYear int) *basis.DaYun {
	// 虛歲計算: 當前年份 - 出生年份 + 1
	age := currentYear - birthYear + 1

	for i := range dayuns {
		if age >= dayuns[i].StartAge && age <= dayuns[i].EndAge {
			return &dayuns[i]
		}
	}
	return nil
}

// TestGetCurrentDaYun 驗證當前大限查找函數
func TestGetCurrentDaYun(t *testing.T) {
	// 測試數據: 1985年出生, 土五局 (5歲起運)
	birthYear := 1985
	mingBranch := basis.BranchZi
	yearStem := basis.StemYi // 乙=陰
	sex := basis.SexFemale   // 陰女=順行
	wuxing := basis.WuxingTu5

	dayuns := CalcDaYun(mingBranch, yearStem, sex, wuxing)

	testCases := []struct {
		currentYear int
		expectIdx   int // 0-based 索引
	}{
		{1985, 0}, // 虛歲1, 第一限(5歲起, 未到運)
		{1989, 0}, // 虛歲5, 第一限(5-14)
		{1994, 0}, // 虛歲10, 第一限(5-14)
		{1995, 0}, // 虛歲11, 還在第一限(5-14)
		{2005, 1}, // 虛歲21, 第二限(15-24)
		{2025, 3}, // 虛歲41, 第四限(35-44)
	}

	for _, tc := range testCases {
		current := GetCurrentDaYun(dayuns, birthYear, tc.currentYear)
		if current == nil {
			// 如果虛歲小於起運歲數, 可能返回nil
			age := tc.currentYear - birthYear + 1
			if age >= dayuns[0].StartAge {
				t.Errorf("%d年(虛歲%d): 應找到大限但未找到", tc.currentYear, age)
			}
			continue
		}

		if current.Index != tc.expectIdx+1 { // DaYun.Index 是 1-based
			t.Errorf("%d年: 期望第%d限, 實際第%d限(%d-%d歲)",
				tc.currentYear, tc.expectIdx+1, current.Index,
				current.StartAge, current.EndAge)
		}
	}
}

// BenchmarkDaYunCalculation 性能測試
func BenchmarkDaYunCalculation(b *testing.B) {
	mingBranch := basis.BranchChen
	yearStem := basis.StemJia
	sex := basis.SexMale
	wuxing := basis.WuxingShui2

	for i := 0; i < b.N; i++ {
		CalcDaYun(mingBranch, yearStem, sex, wuxing)
	}
}

// TestRealWorldExample 真實案例驗證
func TestRealWorldExample(t *testing.T) {
	// 案例: 男性, 1990年6月15日10時 (庚午年)
	// 預期: 陽男, 五行局需要計算

	// 此測試需要完整的命宮和五行局計算
	// 這裡只做框架, 實際測試需要整合 lifepalace.go 和 wuxing.go

	t.Skip("需要整合五行局計算, 暫時跳過")
}
