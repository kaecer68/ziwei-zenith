package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/kaecer68/lunar-zenith/pkg/celestial"
	"github.com/kaecer68/lunar-zenith/pkg/zodiac"
	pb "github.com/kaecer68/ziwei-zenith/pkg/api/grpc/v1"
	v1 "github.com/kaecer68/ziwei-zenith/pkg/api/v1"
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
	"github.com/kaecer68/ziwei-zenith/pkg/engine"
	"github.com/kaecer68/ziwei-zenith/pkg/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CalculateRequest struct {
	Year      int     `json:"year"`
	Month     int     `json:"month"`
	Day       int     `json:"day"`
	Hour      int     `json:"hour"`
	Minute    int     `json:"minute"`
	Gender    string  `json:"gender"`
	IsLunar   bool    `json:"is_lunar"`
	IsLeap    bool    `json:"is_leap"` // Only for Lunar
	IsDST     bool    `json:"is_dst"`
	Longitude float64 `json:"longitude"`
}

var (
	records      []v1.BirthRecord
	tags         []v1.Tag
	mu           sync.RWMutex
	recordsFile  = "records.json"
	tagsFile     = "tags.json"
	contractFile = "contracts/openapi/ziwei-zenith.yaml"
)

func main() {
	loadData()

	// ─── gRPC Server (requires GRPC_PORT env or .env.ports) ───
	go startGRPCServer()

	// ─── REST Server (contract-driven port) ───
	http.HandleFunc("/api/v1/health", healthHandler)
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.HandleFunc("/api/v1/calculate/temporal", temporalCalculateHandler)
	http.HandleFunc("/api/v1/records", recordsHandler)
	http.HandleFunc("/api/v1/records/", recordItemHandler)
	http.HandleFunc("/api/v1/tags", tagsHandler)

	port, err := resolveRESTPort()
	if err != nil {
		log.Fatalf("resolve REST port failed: %v", err)
	}
	grpcPort := getGRPCPort()
	fmt.Printf("Ziwei Zenith REST API on :%s | gRPC on :%s\n", port, grpcPort)
	if err := http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

func resolveRESTPort() (string, error) {
	if port := strings.TrimSpace(os.Getenv("REST_PORT")); port != "" {
		return port, nil
	}

	port, err := restPortFromContract(contractFile)
	if err != nil {
		return "", fmt.Errorf("REST_PORT not set and contract lookup failed: %w", err)
	}

	return port, nil
}

func getGRPCPort() string {
	// Priority: GRPC_PORT env > ZIWEI_GRPC_PORT env > .env.ports > fail
	if port := strings.TrimSpace(os.Getenv("GRPC_PORT")); port != "" {
		return port
	}
	if port := strings.TrimSpace(os.Getenv("ZIWEI_GRPC_PORT")); port != "" {
		return port
	}
	// Try to read from .env.ports if it exists
	if data, err := os.ReadFile(".env.ports"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "GRPC_PORT=") || strings.HasPrefix(line, "ZIWEI_GRPC_PORT=") {
				if port := strings.TrimSpace(strings.SplitN(line, "=", 2)[1]); port != "" {
					return port
				}
			}
		}
	}
	// No fallback - require explicit configuration
	log.Fatal("GRPC_PORT or ZIWEI_GRPC_PORT must be set, or .env.ports must exist with GRPC_PORT defined")
	return "" // unreachable
}

func restPortFromContract(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inServers := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "servers:" {
			inServers = true
			continue
		}
		if !inServers {
			continue
		}
		if strings.HasPrefix(line, "paths:") {
			break
		}
		if !strings.HasPrefix(line, "- url:") {
			continue
		}

		rawURL := strings.TrimSpace(strings.TrimPrefix(line, "- url:"))
		rawURL = strings.Trim(rawURL, `"'`)
		u, parseErr := url.Parse(rawURL)
		if parseErr != nil {
			return "", fmt.Errorf("invalid server url %q: %w", rawURL, parseErr)
		}
		if u.Port() == "" {
			return "", fmt.Errorf("server url %q has no explicit port", rawURL)
		}
		return u.Port(), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("no server url found under servers in contract")
}

func startGRPCServer() {
	grpcPort := getGRPCPort()
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("gRPC listen failed: %v", err)
	}

	s := grpc.NewServer()
	grpcSvc := service.NewZiweiGRPCServer()

	// 同步現有紀錄與標籤到 gRPC 服務
	mu.RLock()
	pbRecords := make([]*pb.BirthRecord, 0, len(records))
	for _, r := range records {
		pbRecords = append(pbRecords, &pb.BirthRecord{
			Id: r.ID, Name: r.Name,
			Year: int32(r.Year), Month: int32(r.Month), Day: int32(r.Day), Hour: int32(r.Hour),
			Gender: r.Gender, IsLunar: r.IsLunar, IsLeap: r.IsLeap, IsDst: r.IsDST,
			Tags: r.Tags, CreatedAt: r.CreatedAt,
		})
	}
	pbTags := make([]*pb.Tag, 0, len(tags))
	for _, t := range tags {
		pbTags = append(pbTags, &pb.Tag{Id: t.ID, Name: t.Name, Color: t.Color})
	}
	mu.RUnlock()

	grpcSvc.SetRecords(pbRecords)
	grpcSvc.SetTags(pbTags)

	pb.RegisterZiweiServiceServer(s, grpcSvc)
	reflection.Register(s)

	log.Printf("gRPC server listening on :%s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC serve failed: %v", err)
	}
}

func loadData() {
	rData, err := os.ReadFile(recordsFile)
	if err == nil {
		json.Unmarshal(rData, &records)
	}

	tData, err := os.ReadFile(tagsFile)
	if err == nil {
		json.Unmarshal(tData, &tags)
	} else {
		// Default tags if none exist
		tags = []v1.Tag{
			{ID: "1", Name: "家人", Color: "#EAB308"},
			{ID: "2", Name: "親戚", Color: "#3B82F6"},
			{ID: "3", Name: "朋友", Color: "#22C55E"},
			{ID: "4", Name: "同事", Color: "#8B5CF6"},
			{ID: "5", Name: "客戶", Color: "#EC4899"},
		}
		saveTags()
	}
}

func saveData() {
	data, _ := json.MarshalIndent(records, "", "  ")
	os.WriteFile(recordsFile, data, 0644)
}

func saveTags() {
	data, _ := json.MarshalIndent(tags, "", "  ")
	os.WriteFile(tagsFile, data, 0644)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Handle DST: Subtract 1 hour if DST is active
	calcHour := req.Hour
	if req.IsDST {
		calcHour--
	}

	loc := time.FixedZone("", int((req.Longitude/15)*3600))
	solarTime := time.Date(req.Year, time.Month(req.Month), req.Day, calcHour, req.Minute, 0, 0, loc)

	sex := basis.SexMale
	if req.Gender == "female" {
		sex = basis.SexFemale
	}

	var lYear, lMonth, lDay int
	var yPillar, mPillar, dPillar basis.Pillar

	pt := celestial.NewPrecisionTime(solarTime)
	pillar := zodiac.GetAstrologicalPillar(pt)

	yPillar = basis.Pillar{Stem: basis.Stem(pillar.Year.StemIndex), Branch: basis.Branch(pillar.Year.BranchIndex)}
	mPillar = basis.Pillar{Stem: basis.Stem(pillar.Month.StemIndex), Branch: basis.Branch(pillar.Month.BranchIndex)}
	dPillar = basis.Pillar{Stem: basis.Stem(pillar.Day.StemIndex), Branch: basis.Branch(pillar.Day.BranchIndex)}

	if req.IsLunar {
		lYear = req.Year
		lMonth = req.Month
		if req.IsLeap {
			lMonth = -req.Month
		}
		lDay = req.Day
	} else {
		jd := celestial.TimeToJD(solarTime)
		engine_lunar := &zodiac.LunarEngine{}
		lunarDate := engine_lunar.GetLunarDate(jd)

		lYear = lunarDate.Year
		lMonth = lunarDate.Month
		lDay = lunarDate.Day
	}

	hourBranchIdx := zodiac.GetHourBranch(req.Hour)
	hourSexagenary := zodiac.GetHourSexagenary(int(dPillar.Stem), hourBranchIdx)

	birth := basis.BirthInfo{
		SolarYear:   req.Year,
		SolarMonth:  req.Month,
		SolarDay:    req.Day,
		Hour:        req.Hour,
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
		http.Error(w, fmt.Sprintf("Calculation error: %v", err), http.StatusInternalServerError)
		return
	}

	response := mapChartToResponse(chart, req.Gender, req.Year)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TemporalCalculateRequest 動態運限計算請求
type TemporalCalculateRequest struct {
	BirthYear   int    `json:"birth_year"`
	BirthMonth  int    `json:"birth_month"`
	BirthDay    int    `json:"birth_day"`
	BirthHour   int    `json:"birth_hour"`
	Gender      string `json:"gender"`
	IsLunar     bool   `json:"is_lunar"`
	IsLeap      bool   `json:"is_leap"`
	DaYunIndex  int    `json:"da_yun_index"`
	TargetYear  int    `json:"target_year,omitempty"`
	TargetMonth int    `json:"target_month,omitempty"`
	TargetDay   int    `json:"target_day,omitempty"`
	LunarMonth  int    `json:"lunar_month,omitempty"`
}

// TemporalCalculateResponse 動態運限計算響應
type TemporalCalculateResponse struct {
	DaYun      v1.DaYunData          `json:"da_yun"`
	LiuNian    v1.TemporalPalaceData `json:"liu_nian"`
	LiuYue     v1.TemporalPalaceData `json:"liu_yue,omitempty"`
	LiuRi      v1.TemporalPalaceData `json:"liu_ri,omitempty"`
	LunarYear  int                   `json:"lunar_year,omitempty"`
	LunarMonth int                   `json:"lunar_month,omitempty"`
	LunarDay   int                   `json:"lunar_day,omitempty"`
	LunarDays  []int                 `json:"lunar_days,omitempty"`
	// 各宮位的流運四化數據 - 用於三方四正顯示
	DaYunPalaceTransforms   map[string][]v1.TransformData `json:"da_yun_palace_transforms,omitempty"`
	LiuNianPalaceTransforms map[string][]v1.TransformData `json:"liu_nian_palace_transforms,omitempty"`
	LiuYuePalaceTransforms  map[string][]v1.TransformData `json:"liu_yue_palace_transforms,omitempty"`
	LiuRiPalaceTransforms   map[string][]v1.TransformData `json:"liu_ri_palace_transforms,omitempty"`
}

func temporalCalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TemporalCalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 驗證必要參數
	if req.BirthYear == 0 || req.BirthMonth == 0 || req.BirthDay == 0 || req.Gender == "" {
		http.Error(w, "Missing required birth info", http.StatusBadRequest)
		return
	}
	if req.DaYunIndex < 0 || req.DaYunIndex > 11 {
		http.Error(w, "Invalid da_yun_index (must be 0-11)", http.StatusBadRequest)
		return
	}

	// 使用當前時間作為默認值
	now := time.Now()
	targetYear := req.TargetYear
	if targetYear == 0 {
		targetYear = now.Year()
	}
	targetMonth := req.TargetMonth
	if targetMonth == 0 {
		targetMonth = int(now.Month())
	}
	targetDay := req.TargetDay
	if targetDay == 0 {
		targetDay = now.Day()
	}

	// 構建出生信息並計算命盤
	sex := basis.SexMale
	if req.Gender == "female" {
		sex = basis.SexFemale
	}

	calcHour := req.BirthHour

	lon := 121.565 // 預設台北
	loc := time.FixedZone("", int((lon/15)*3600))
	solarTime := time.Date(req.BirthYear, time.Month(req.BirthMonth), req.BirthDay, calcHour, 0, 0, 0, loc)

	pt := celestial.NewPrecisionTime(solarTime)
	pillar := zodiac.GetAstrologicalPillar(pt)

	yPillar := basis.Pillar{Stem: basis.Stem(pillar.Year.StemIndex), Branch: basis.Branch(pillar.Year.BranchIndex)}
	mPillar := basis.Pillar{Stem: basis.Stem(pillar.Month.StemIndex), Branch: basis.Branch(pillar.Month.BranchIndex)}
	dPillar := basis.Pillar{Stem: basis.Stem(pillar.Day.StemIndex), Branch: basis.Branch(pillar.Day.BranchIndex)}

	localDayPillar := zodiac.GetDaySexagenary(celestial.TimeToJD(solarTime))
	dPillar = basis.Pillar{Stem: basis.Stem(localDayPillar.StemIndex), Branch: basis.Branch(localDayPillar.BranchIndex)}

	var lYear, lMonth, lDay int
	if req.IsLunar {
		lYear = req.BirthYear
		lMonth = req.BirthMonth
		if req.IsLeap {
			lMonth = -req.BirthMonth
		}
		lDay = req.BirthDay
	} else {
		jd := celestial.TimeToJD(solarTime)
		lunarEngine := &zodiac.LunarEngine{}
		lunarDate := lunarEngine.GetLunarDate(jd)
		lYear = lunarDate.Year
		lMonth = lunarDate.Month
		lDay = lunarDate.Day
	}

	hourBranchIdx := zodiac.GetHourBranch(req.BirthHour)
	hourSexagenary := zodiac.GetHourSexagenary(int(dPillar.Stem), hourBranchIdx)

	birth := basis.BirthInfo{
		SolarYear:   req.BirthYear,
		SolarMonth:  req.BirthMonth,
		SolarDay:    req.BirthDay,
		Hour:        req.BirthHour,
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
		http.Error(w, fmt.Sprintf("Calculation error: %v", err), http.StatusInternalServerError)
		return
	}

	// 獲取指定的大限
	if req.DaYunIndex >= len(chart.DaYun) {
		http.Error(w, "DaYun index out of range", http.StatusBadRequest)
		return
	}
	selectedDaYun := chart.DaYun[req.DaYunIndex]

	// 計算流年（使用 targetDate，年界以立春切換由 astrological pillar 處理）
	targetTime := time.Date(targetYear, time.Month(targetMonth), targetDay, 12, 0, 0, 0, loc)
	targetPillar := zodiac.GetAstrologicalPillar(celestial.NewPrecisionTime(targetTime))
	lnBranch := basis.Branch(targetPillar.Year.BranchIndex)
	lnStem := basis.Stem(targetPillar.Year.StemIndex)

	// 優先使用前端傳來的農曆月份，否則從 targetDate 計算
	var lunarMonth int
	var lunarDay int
	lunarEngine := &zodiac.LunarEngine{}
	if req.LunarMonth > 0 {
		lunarMonth = req.LunarMonth
		// 如果有傳 LunarDay 就用，否則用 targetDay
		if req.TargetDay > 0 {
			lunarDay = req.TargetDay
		} else {
			// 嘗試從 targetDate 計算農曆日
			lunarDay = targetDay
		}
	} else {
		// 從 targetDate 計算正確的農曆年月
		targetLunarDate := lunarEngine.GetLunarDate(celestial.TimeToJD(targetTime))
		lunarMonth = targetLunarDate.Month
		if lunarMonth < 0 {
			lunarMonth = -lunarMonth
		}
		lunarDay = targetLunarDate.Day
	}

	// 計算該月的農曆天數
	lunarDaysInMonth := 30 // 預設
	if targetYear > 0 && lunarMonth > 0 && lunarMonth <= 12 {
		// 嘗試獲取該月的天數：使用陽曆日期轉換
		// 假設該月15日，轉換為儒略日
		testDate := time.Date(targetYear, time.Month(targetMonth), 15, 12, 0, 0, 0, loc)
		jd := celestial.TimeToJD(testDate)
		testLunarDate := lunarEngine.GetLunarDate(jd)
		// 嘗試獲取該月最後一天
		for d := 28; d <= 30; d++ {
			testDate2 := time.Date(targetYear, time.Month(targetMonth), d, 12, 0, 0, 0, loc)
			jd2 := celestial.TimeToJD(testDate2)
			lunarDate2 := lunarEngine.GetLunarDate(jd2)
			if lunarDate2.Month != testLunarDate.Month && lunarDate2.Month != -testLunarDate.Month {
				lunarDaysInMonth = d - 1
				break
			}
			lunarDaysInMonth = d
		}
	}

	// 計算流月
	// 使用 Dou Jun 方法：從流年地支開始，逆數到出生月，再順數到出生時
	birthMonth := req.BirthMonth
	if req.IsLunar {
		birthMonth = req.BirthMonth
	}
	hourBranch := basis.HourBranch(hourBranchIdx)
	month1Idx := (int(lnBranch) - (birthMonth - 1) + int(hourBranch) + 12) % 12
	lYueIdx := (month1Idx + (lunarMonth - 1)) % 12
	lyBranch := basis.Branch(lYueIdx)

	// 計算流日（使用計算出的農曆日，但不能超過該月實際天數）
	actualDay := lunarDay
	if actualDay > lunarDaysInMonth {
		actualDay = lunarDaysInMonth
	}
	lRiIdx := (int(lyBranch) + (actualDay - 1)) % 12
	lrBranch := basis.Branch(lRiIdx)

	response := TemporalCalculateResponse{
		DaYun: v1.DaYunData{
			Index:    selectedDaYun.Index,
			StartAge: selectedDaYun.StartAge,
			EndAge:   selectedDaYun.EndAge,
			Stem:     chart.PalaceGans[selectedDaYun.Branch].String(),
			Branch:   selectedDaYun.Branch.String(),
			Palace:   chart.Palaces[selectedDaYun.Branch].String(),
		},
		LiuNian: v1.TemporalPalaceData{
			Label:      fmt.Sprintf("流年（%d）", targetYear),
			Branch:     lnBranch.String(),
			Palace:     chart.Palaces[lnBranch].String(),
			Stem:       lnStem.String(),
			TimeBranch: lnBranch.String(),
		},
		LiuYue: v1.TemporalPalaceData{
			Label:      fmt.Sprintf("流月（農曆%d月）", lunarMonth),
			Branch:     lyBranch.String(),
			Palace:     chart.Palaces[lyBranch].String(),
			Stem:       basis.Stem(targetPillar.Month.StemIndex).String(),
			TimeBranch: basis.Branch(targetPillar.Month.BranchIndex).String(),
		},
		LiuRi: v1.TemporalPalaceData{
			Label:      fmt.Sprintf("流日（農曆%d日）", lunarDay),
			Branch:     lrBranch.String(),
			Palace:     chart.Palaces[lrBranch].String(),
			Stem:       basis.Stem(targetPillar.Day.StemIndex).String(),
			TimeBranch: basis.Branch(targetPillar.Day.BranchIndex).String(),
		},
		LunarYear:  targetYear,
		LunarMonth: lunarMonth,
		LunarDay:   lunarDay,
		LunarDays: func() []int {
			days := make([]int, lunarDaysInMonth)
			for i := range days {
				days[i] = i + 1
			}
			return days
		}(),
		DaYunPalaceTransforms:   calculateTemporalTransforms(chart, chart.PalaceGans[selectedDaYun.Branch]),
		LiuNianPalaceTransforms: calculateTemporalTransforms(chart, lnStem),
		LiuYuePalaceTransforms:  calculateTemporalTransforms(chart, basis.Stem(targetPillar.Month.StemIndex)),
		LiuRiPalaceTransforms:   calculateTemporalTransforms(chart, basis.Stem(targetPillar.Day.StemIndex)),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func recordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.RLock()
		defer mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)
	case http.MethodPost:
		var record v1.BirthRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		record.ID = fmt.Sprintf("%d", time.Now().UnixNano())
		record.CreatedAt = time.Now().Format(time.RFC3339)
		mu.Lock()
		records = append(records, record)
		saveData()
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(record)
	}
}

func recordItemHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/records/")
	mu.Lock()
	defer mu.Unlock()

	idx := -1
	for i, r := range records {
		if r.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		http.Error(w, "Record not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodPut:
		var record v1.BirthRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		record.ID = id
		records[idx] = record
		saveData()
		json.NewEncoder(w).Encode(record)
	case http.MethodDelete:
		records = append(records[:idx], records[idx+1:]...)
		saveData()
		w.WriteHeader(http.StatusNoContent)
	}
}

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		mu.RLock()
		defer mu.RUnlock()
		json.NewEncoder(w).Encode(tags)
	case http.MethodPut:
		var newTags []v1.Tag
		if err := json.NewDecoder(r.Body).Decode(&newTags); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		tags = newTags
		saveTags()
		mu.Unlock()
		json.NewEncoder(w).Encode(tags)
	}
}

func mapChartToResponse(chart *engine.ZiweiChart, gender string, birthYear int) v1.ZiweiResponse {
	palaces := make(map[string]v1.PalaceData)
	for i := 0; i < 12; i++ {
		b := basis.Branch(i)
		pType := chart.Palaces[b]
		starNames := make([]string, 0)
		starDetails := make([]v1.PalaceStar, 0)
		for _, s := range chart.Stars[b] {
			starNames = append(starNames, s.String())
			starDetails = append(starDetails, v1.PalaceStar{
				Name:       s.String(),
				Brightness: basis.BrightnessLevel(s, b).String(),
			})
		}

		assistantStars := stringifyInterfaces(chart.AssistantStars[b])
		secondaryStars := stringifyInterfaces(chart.SecondaryStars[b])
		natalTransforms := transformDataFromInterfaces(chart.TransformedStars[b])
		lnStars := stringifyInterfaces(chart.LiuNianStars[b])
		lyStars := stringifyInterfaces(chart.LiuYueStars[b])
		lrStars := stringifyInterfaces(chart.LiuRiStars[b])
		lnTransforms := transformDataFromInterfaces(chart.LiuNianStars[b])
		lyTransforms := transformDataFromInterfaces(chart.LiuYueStars[b])
		lrTransforms := transformDataFromInterfaces(chart.LiuRiStars[b])

		daYunAges := make([]string, 0)
		for _, dy := range chart.DaYun {
			if dy.Branch == b {
				daYunAges = append(daYunAges, fmt.Sprintf("%d-%d", dy.StartAge, dy.EndAge))
			}
		}

		palaces[pType.String()] = v1.PalaceData{
			Branch:            b.String(),
			PalaceGan:         chart.PalaceGans[b].String(),
			Stars:             starNames,
			StarDetails:       starDetails,
			AssistantStars:    assistantStars,
			SecondaryStars:    secondaryStars,
			ChangSheng:        chart.ChangSheng[b].String(),
			BoShi:             chart.BoShi[b].String(),
			NatalTransforms:   natalTransforms,
			LiuNianStars:      lnStars,
			LiuNianTransforms: lnTransforms,
			LiuYueStars:       lyStars,
			LiuYueTransforms:  lyTransforms,
			LiuRiStars:        lrStars,
			LiuRiTransforms:   lrTransforms,
			DaYunAges:         daYunAges,
			FlyHua:            buildPalaceFlyHua(chart, b),
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

	starDetails := make([]v1.DeepStarAnalysis, 0)
	for _, s := range chart.Interpretation.StarDetails {
		starDetails = append(starDetails, v1.DeepStarAnalysis{
			Name:       s.Name,
			Verse:      s.Verse,
			Positive:   s.Positive,
			Negative:   s.Negative,
			Remedy:     s.Remedy,
			Evolution:  s.Evolution,
			Brightness: s.Brightness,
		})
	}

	stages := make([]v1.FlyStage, 0)
	for _, s := range chart.Interpretation.OriginFlyHua.Stages {
		stages = append(stages, v1.FlyStage{
			Type:   s.Type,
			Star:   s.Star,
			Target: s.Target,
			Motive: s.Motive,
			Action: s.Action,
			Trap:   s.Trap,
			Interpretations: v1.MultiSchoolView{
				SanHe:   s.Interpretations.SanHe,
				SiHua:   s.Interpretations.SiHua,
				QinTian: s.Interpretations.QinTian,
			},
		})
	}

	daYun := make([]v1.DaYunData, 0, len(chart.DaYun))
	var currentDaYun *v1.DaYunData

	// 計算當前年齡（虛歲）對應的大限
	currentYear := time.Now().Year()
	virtualAge := currentYear - birthYear + 1
	currentDaYunIndex := (virtualAge - 1) / 10 // 0-9歲為第1大限，10-19歲為第2大限，以此類推

	for _, dy := range chart.DaYun {
		item := v1.DaYunData{
			Index:    dy.Index,
			StartAge: dy.StartAge,
			EndAge:   dy.EndAge,
			Stem:     chart.PalaceGans[dy.Branch].String(),
			Branch:   dy.Branch.String(),
			Palace:   chart.Palaces[dy.Branch].String(),
		}
		daYun = append(daYun, item)
		// 根據當前年齡找到對應的當前大限
		if dy.Index == currentDaYunIndex+1 && currentDaYun == nil {
			copyItem := item
			currentDaYun = &copyItem
		}
	}

	liuNian := &v1.TemporalPalaceData{
		Label:      "流年",
		Branch:     chart.LiuNian.Branch.String(),
		Palace:     chart.Palaces[chart.LiuNian.Branch].String(),
		Stem:       chart.LiuNian.Stem.String(),
		TimeBranch: chart.LiuNian.Branch.String(),
	}
	liuYue := &v1.TemporalPalaceData{
		Label:      "流月",
		Branch:     chart.LiuYue.String(),
		Palace:     chart.Palaces[chart.LiuYue].String(),
		Stem:       chart.MonthPillar.Stem.String(),
		TimeBranch: chart.MonthPillar.Branch.String(),
	}
	liuRi := &v1.TemporalPalaceData{
		Label:      "流日",
		Branch:     chart.LiuRi.String(),
		Palace:     chart.Palaces[chart.LiuRi].String(),
		Stem:       chart.DayPillar.Stem.String(),
		TimeBranch: chart.DayPillar.Branch.String(),
	}

	return v1.ZiweiResponse{
		Gender:       gender,
		Wuxing:       chart.Wuxing.String(),
		NaYin:        chart.NaYin.String(),
		OriginPalace: chart.Palaces[chart.OriginPalace].String(),
		MingGong:     chart.LifePalace.MingGong.String(),
		ShenGong:     chart.LifePalace.ShenGong.String(),
		YearPillar:   chart.YearPillar.String(),
		DayPillar:    chart.DayPillar.String(),
		CurrentDaYun: currentDaYun,
		DaYun:        daYun,
		LiuNian:      liuNian,
		LiuYue:       liuYue,
		LiuRi:        liuRi,
		Palaces:      palaces,
		Patterns:     patterns,
		Interpretation: v1.InterpretationData{
			Summary:              chart.Interpretation.Summary,
			CharacterTraits:      chart.Interpretation.CharacterTraits,
			OriginPalaceAnalysis: chart.Interpretation.OriginPalaceAnalysis,
			KarmicNarrative:      narrative,
			SanFangDiagnosis:     diagnosis,
			StarDetails:          starDetails,
			OriginFlyHua: v1.FlyHuaAnalysis{
				FromPalace: chart.Interpretation.OriginFlyHua.FromPalace,
				Stem:       chart.Interpretation.OriginFlyHua.Stem,
				Stages:     stages,
			},
			TemporalResonance: resonance,
			ClassicPatterns:   chart.Interpretation.ClassicPatterns,
		},
	}
}

func stringifyInterfaces(items []interface{}) []string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		if strer, ok := item.(interface{ String() string }); ok {
			values = append(values, strer.String())
		}
	}
	return values
}

func transformDataFromInterfaces(items []interface{}) []v1.TransformData {
	values := make([]v1.TransformData, 0)
	for _, item := range items {
		if ts, ok := item.(basis.TransformedStar); ok {
			values = append(values, v1.TransformData{
				Star:           ts.StarName,
				Transformation: ts.Transformation.String(),
				Display:        ts.String(),
			})
		}
	}
	return values
}

func buildPalaceFlyHua(chart *engine.ZiweiChart, branch basis.Branch) v1.FlyHuaAnalysis {
	stem := chart.PalaceGans[branch]
	hua, ok := basis.TransformationTable[stem]
	if !ok {
		return v1.FlyHuaAnalysis{}
	}

	result := v1.FlyHuaAnalysis{
		FromPalace: chart.Palaces[branch].String(),
		Stem:       stem.String(),
		Stages:     make([]v1.FlyStage, 0, 4),
	}

	types := []string{"祿", "權", "科", "忌"}
	for i, t := range types {
		target := findStarTargetPalace(chart, hua[i])
		result.Stages = append(result.Stages, v1.FlyStage{
			Type:   t,
			Star:   hua[i],
			Target: target,
		})
	}

	return result
}

func findStarTargetPalace(chart *engine.ZiweiChart, starName string) string {
	for b, stars := range chart.Stars {
		for _, s := range stars {
			if s.String() == starName {
				return chart.Palaces[b].String()
			}
		}
	}
	for b, stars := range chart.AssistantStars {
		for _, s := range stars {
			if strer, ok := s.(interface{ String() string }); ok && strer.String() == starName {
				return chart.Palaces[b].String()
			}
		}
	}
	for b, stars := range chart.SecondaryStars {
		for _, s := range stars {
			if strer, ok := s.(interface{ String() string }); ok && strer.String() == starName {
				return chart.Palaces[b].String()
			}
		}
	}
	return ""
}

func calculateTemporalTransforms(chart *engine.ZiweiChart, stem basis.Stem) map[string][]v1.TransformData {
	result := make(map[string][]v1.TransformData)

	table, ok := basis.TransformationTable[stem]
	if !ok {
		return result
	}

	transTypes := []string{"祿", "權", "科", "忌"}

	for i, starName := range table {
		transType := transTypes[i]

		for branchIdx := 0; branchIdx < 12; branchIdx++ {
			b := basis.Branch(branchIdx)
			palaceName := chart.Palaces[b].String()

			allStars := []string{}
			for _, s := range chart.Stars[b] {
				allStars = append(allStars, s.String())
			}
			for _, s := range chart.AssistantStars[b] {
				if strer, ok := s.(interface{ String() string }); ok {
					allStars = append(allStars, strer.String())
				}
			}
			for _, s := range chart.SecondaryStars[b] {
				if strer, ok := s.(interface{ String() string }); ok {
					allStars = append(allStars, strer.String())
				}
			}

			for _, s := range allStars {
				if s == starName {
					result[palaceName] = append(result[palaceName], v1.TransformData{
						Star:           starName,
						Transformation: "化" + transType,
						Display:        starName + "化" + transType,
					})
					break
				}
			}
		}
	}

	return result
}
