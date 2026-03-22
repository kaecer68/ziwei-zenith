package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/kaecer68/ziwei-zenith/pkg/api/grpc/v1"
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
	"github.com/kaecer68/ziwei-zenith/pkg/engine"
)

// ZiweiGRPCServer 實作 gRPC ZiweiService
type ZiweiGRPCServer struct {
	pb.UnimplementedZiweiServiceServer
	mu      sync.RWMutex
	records []*pb.BirthRecord
	tags    []*pb.Tag
	// 資料持久化回呼（由外部注入）
	OnRecordsChanged func()
	OnTagsChanged    func()
}

func NewZiweiGRPCServer() *ZiweiGRPCServer {
	return &ZiweiGRPCServer{}
}

// SetRecords 從外部同步紀錄資料（REST ↔ gRPC 共享）
func (s *ZiweiGRPCServer) SetRecords(records []*pb.BirthRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = records
}

// SetTags 從外部同步標籤資料
func (s *ZiweiGRPCServer) SetTags(tags []*pb.Tag) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tags = tags
}

// Calculate 計算命盤
func (s *ZiweiGRPCServer) Calculate(ctx context.Context, req *pb.CalculateRequest) (*pb.CalculateResponse, error) {
	input := CalculateInput{
		Year:    int(req.Year),
		Month:   int(req.Month),
		Day:     int(req.Day),
		Hour:    int(req.Hour),
		Gender:  req.Gender,
		IsLunar: req.IsLunar,
		IsLeap:  req.IsLeap,
		IsDST:   req.IsDst,
	}

	chart, err := Calculate(input)
	if err != nil {
		return nil, err
	}

	return chartToProto(chart, req.Gender), nil
}

// ListRecords 列出所有紀錄
func (s *ZiweiGRPCServer) ListRecords(ctx context.Context, req *pb.ListRecordsRequest) (*pb.ListRecordsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &pb.ListRecordsResponse{Records: s.records}, nil
}

// CreateRecord 建立紀錄
func (s *ZiweiGRPCServer) CreateRecord(ctx context.Context, req *pb.CreateRecordRequest) (*pb.CreateRecordResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := &pb.BirthRecord{
		Id:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Name:      req.Name,
		Year:      req.Year,
		Month:     req.Month,
		Day:       req.Day,
		Hour:      req.Hour,
		Gender:    req.Gender,
		IsLunar:   req.IsLunar,
		IsLeap:    req.IsLeap,
		IsDst:     req.IsDst,
		Tags:      req.Tags,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	s.records = append(s.records, record)

	if s.OnRecordsChanged != nil {
		s.OnRecordsChanged()
	}

	return &pb.CreateRecordResponse{Record: record}, nil
}

// DeleteRecord 刪除紀錄
func (s *ZiweiGRPCServer) DeleteRecord(ctx context.Context, req *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, r := range s.records {
		if r.Id == req.Id {
			s.records = append(s.records[:i], s.records[i+1:]...)
			if s.OnRecordsChanged != nil {
				s.OnRecordsChanged()
			}
			return &pb.DeleteRecordResponse{Success: true}, nil
		}
	}
	return &pb.DeleteRecordResponse{Success: false}, nil
}

// ListTags 列出標籤
func (s *ZiweiGRPCServer) ListTags(ctx context.Context, req *pb.ListTagsRequest) (*pb.ListTagsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &pb.ListTagsResponse{Tags: s.tags}, nil
}

// ─── 轉換函式 ───

func chartToProto(chart *engine.ZiweiChart, gender string) *pb.CalculateResponse {
	resp := &pb.CalculateResponse{
		Gender:       gender,
		Wuxing:       chart.Wuxing.String(),
		NaYin:        chart.NaYin.String(),
		OriginPalace: chart.Palaces[chart.OriginPalace].String(),
		MingGong:     chart.LifePalace.MingGong.String(),
		ShenGong:     chart.LifePalace.ShenGong.String(),
		YearPillar:   chart.YearPillar.String(),
		DayPillar:    chart.DayPillar.String(),
		Palaces:      make(map[string]*pb.PalaceData),
	}

	for i := 0; i < 12; i++ {
		b := basis.Branch(i)
		pType := chart.Palaces[b]

		starNames := make([]string, 0)
		starDetails := make([]*pb.PalaceStar, 0)
		for _, s := range chart.Stars[b] {
			starNames = append(starNames, s.String())
			starDetails = append(starDetails, &pb.PalaceStar{
				Name:       s.String(),
				Brightness: basis.BrightnessLevel(s, b).String(),
			})
		}

		pd := &pb.PalaceData{
			Branch:            b.String(),
			PalaceGan:         chart.PalaceGans[b].String(),
			Stars:             starNames,
			StarDetails:       starDetails,
			AssistantStars:    stringifyInterfaces(chart.AssistantStars[b]),
			SecondaryStars:    stringifyInterfaces(chart.SecondaryStars[b]),
			ChangSheng:        chart.ChangSheng[b].String(),
			BoShi:             chart.BoShi[b].String(),
			NatalTransforms:   transformsToProto(chart.TransformedStars[b]),
			LiuNianStars:      stringifyInterfaces(chart.LiuNianStars[b]),
			LiuNianTransforms: transformsToProto(chart.LiuNianStars[b]),
			LiuYueStars:       stringifyInterfaces(chart.LiuYueStars[b]),
			LiuYueTransforms:  transformsToProto(chart.LiuYueStars[b]),
			LiuRiStars:        stringifyInterfaces(chart.LiuRiStars[b]),
			LiuRiTransforms:   transformsToProto(chart.LiuRiStars[b]),
		}

		daYunAges := make([]string, 0)
		for _, dy := range chart.DaYun {
			if dy.Branch == b {
				daYunAges = append(daYunAges, fmt.Sprintf("%d-%d", dy.StartAge, dy.EndAge))
			}
		}
		pd.DaYunAges = daYunAges

		resp.Palaces[pType.String()] = pd
	}

	// 格局
	for _, p := range chart.Patterns {
		resp.Patterns = append(resp.Patterns, &pb.PatternData{
			Name:        p.Name,
			Description: p.Description,
			Level:       p.Level,
		})
	}

	// 大運
	for _, dy := range chart.DaYun {
		resp.DaYun = append(resp.DaYun, &pb.DaYunData{
			Index:    int32(dy.Index),
			StartAge: int32(dy.StartAge),
			EndAge:   int32(dy.EndAge),
			Stem:     chart.PalaceGans[dy.Branch].String(),
			Branch:   dy.Branch.String(),
			Palace:   chart.Palaces[dy.Branch].String(),
		})
	}

	resp.LiuNian = &pb.TemporalPalaceData{
		Label:      "流年",
		Branch:     chart.LiuNian.Branch.String(),
		Palace:     chart.Palaces[chart.LiuNian.Branch].String(),
		Stem:       chart.LiuNian.Stem.String(),
		TimeBranch: chart.LiuNian.Branch.String(),
	}
	resp.LiuYue = &pb.TemporalPalaceData{
		Label:      "流月",
		Branch:     chart.LiuYue.String(),
		Palace:     chart.Palaces[chart.LiuYue].String(),
		Stem:       chart.MonthPillar.Stem.String(),
		TimeBranch: chart.MonthPillar.Branch.String(),
	}
	resp.LiuRi = &pb.TemporalPalaceData{
		Label:      "流日",
		Branch:     chart.LiuRi.String(),
		Palace:     chart.Palaces[chart.LiuRi].String(),
		Stem:       chart.DayPillar.Stem.String(),
		TimeBranch: chart.DayPillar.Branch.String(),
	}

	// 解讀
	interp := &pb.InterpretationData{
		Summary:              chart.Interpretation.Summary,
		OriginPalaceAnalysis: chart.Interpretation.OriginPalaceAnalysis,
		ClassicPatterns:      chart.Interpretation.ClassicPatterns,
	}
	for _, k := range chart.Interpretation.KarmicNarrative {
		interp.KarmicNarrative = append(interp.KarmicNarrative, &pb.KarmicStep{
			Type: k.Type, Role: k.Role, Star: k.Star, Palace: k.Palace, Desc: k.Desc,
		})
	}
	for _, r := range chart.Interpretation.SanFangDiagnosis {
		interp.SanFangDiagnosis = append(interp.SanFangDiagnosis, &pb.SanFangRole{
			Role: r.Role, Palace: r.Palace, Diagnosis: r.Diagnosis,
		})
	}
	for _, s := range chart.Interpretation.StarDetails {
		interp.StarDetails = append(interp.StarDetails, &pb.DeepStarAnalysis{
			Name: s.Name, Verse: s.Verse, Positive: s.Positive,
			Negative: s.Negative, Remedy: s.Remedy, Evolution: s.Evolution,
			Brightness: s.Brightness,
		})
	}
	if chart.Interpretation.OriginFlyHua.FromPalace != "" {
		flyHua := &pb.FlyHuaAnalysis{
			FromPalace: chart.Interpretation.OriginFlyHua.FromPalace,
			Stem:       chart.Interpretation.OriginFlyHua.Stem,
		}
		for _, st := range chart.Interpretation.OriginFlyHua.Stages {
			flyHua.Stages = append(flyHua.Stages, &pb.FlyStage{
				Type: st.Type, Star: st.Star, Target: st.Target,
				Motive: st.Motive, Action: st.Action, Trap: st.Trap,
			})
		}
		interp.OriginFlyHua = flyHua
	}
	for _, r := range chart.Interpretation.TemporalResonance {
		interp.TemporalResonance = append(interp.TemporalResonance, &pb.ResonancePoint{
			Layer: r.Layer, Type: r.Type, Star: r.Star,
			Natal: r.Natal, Palace: r.Palace, Mood: r.Mood,
		})
	}
	resp.Interpretation = interp

	return resp
}

func stringifyInterfaces(items []interface{}) []string {
	result := make([]string, 0)
	for _, item := range items {
		if strer, ok := item.(interface{ String() string }); ok {
			result = append(result, strer.String())
		}
	}
	return result
}

func transformsToProto(items []interface{}) []*pb.TransformData {
	result := make([]*pb.TransformData, 0)
	for _, item := range items {
		type transformStringer interface {
			TransformationType() string
			String() string
		}
		if ts, ok := item.(transformStringer); ok {
			result = append(result, &pb.TransformData{
				Star:           ts.String(),
				Transformation: ts.TransformationType(),
				Display:        ts.String() + ts.TransformationType(),
			})
		}
	}
	return result
}
