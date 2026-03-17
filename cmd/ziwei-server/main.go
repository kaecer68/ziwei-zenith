package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
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
	records     []v1.BirthRecord
	tags        []v1.Tag
	mu          sync.RWMutex
	recordsFile = "records.json"
	tagsFile    = "tags.json"
)

func main() {
	loadData()

	// ─── gRPC Server (port 50051) ───
	go startGRPCServer()

	// ─── REST Server (port 8081) ───
	http.HandleFunc("/api/v1/health", healthHandler)
	http.HandleFunc("/api/v1/calculate", calculateHandler)
	http.HandleFunc("/api/v1/records", recordsHandler)
	http.HandleFunc("/api/v1/records/", recordItemHandler)
	http.HandleFunc("/api/v1/tags", tagsHandler)

	port := "8081"
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50053"
	}
	fmt.Printf("Ziwei Zenith REST API on :%s | gRPC on :%s\n", port, grpcPort)
	if err := http.ListenAndServe(":"+port, corsMiddleware(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}

func startGRPCServer() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50053"
	}
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

	response := mapChartToResponse(chart, req.Gender)
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

func mapChartToResponse(chart *engine.ZiweiChart, gender string) v1.ZiweiResponse {
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
		if dy.Index == 1 && currentDaYun == nil {
			copyItem := item
			currentDaYun = &copyItem
		}
	}

	liuNian := &v1.TemporalPalaceData{
		Label:  "流年",
		Branch: chart.LiuNian.Branch.String(),
		Palace: chart.Palaces[chart.LiuNian.Branch].String(),
		Stem:   chart.LiuNian.Stem.String(),
	}
	liuYue := &v1.TemporalPalaceData{
		Label:  "流月",
		Branch: chart.LiuYue.String(),
		Palace: chart.Palaces[chart.LiuYue].String(),
		Stem:   chart.MonthPillar.Stem.String(),
	}
	liuRi := &v1.TemporalPalaceData{
		Label:  "流日",
		Branch: chart.LiuRi.String(),
		Palace: chart.Palaces[chart.LiuRi].String(),
		Stem:   chart.DayPillar.Stem.String(),
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
