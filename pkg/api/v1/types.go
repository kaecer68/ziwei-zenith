package v1

type ZiweiResponse struct {
	Gender         string                `json:"gender"`
	Wuxing         string                `json:"wuxing"`
	NaYin          string                `json:"na_yin"`
	OriginPalace   string                `json:"origin_palace"`
	MingGong       string                `json:"ming_gong"`
	ShenGong       string                `json:"shen_gong"`
	YearPillar     string                `json:"year_pillar"`
	DayPillar      string                `json:"day_pillar"`
	CurrentDaYun   *DaYunData            `json:"current_da_yun,omitempty"`
	DaYun          []DaYunData           `json:"da_yun,omitempty"`
	LiuNian        *TemporalPalaceData   `json:"liu_nian,omitempty"`
	LiuYue         *TemporalPalaceData   `json:"liu_yue,omitempty"`
	LiuRi          *TemporalPalaceData   `json:"liu_ri,omitempty"`
	Palaces        map[string]PalaceData `json:"palaces"`
	Patterns       []PatternData         `json:"patterns"`
	Interpretation InterpretationData    `json:"interpretation"`
}

type PalaceData struct {
	Branch            string          `json:"branch"`
	PalaceGan         string          `json:"palace_gan,omitempty"`
	Stars             []string        `json:"stars"`
	StarDetails       []PalaceStar    `json:"star_details,omitempty"`
	AssistantStars    []string        `json:"assistant_stars,omitempty"`
	SecondaryStars    []string        `json:"secondary_stars,omitempty"`
	ChangSheng        string          `json:"chang_sheng,omitempty"`
	BoShi             string          `json:"bo_shi,omitempty"`
	NatalTransforms   []TransformData `json:"natal_transforms,omitempty"`
	LiuNianStars      []string        `json:"liu_nian_stars,omitempty"`
	LiuNianTransforms []TransformData `json:"liu_nian_transforms,omitempty"`
	LiuYueStars       []string        `json:"liu_yue_stars,omitempty"`
	LiuYueTransforms  []TransformData `json:"liu_yue_transforms,omitempty"`
	LiuRiStars        []string        `json:"liu_ri_stars,omitempty"`
	LiuRiTransforms   []TransformData `json:"liu_ri_transforms,omitempty"`
	DaYunAges         []string        `json:"da_yun_ages,omitempty"`
	FlyHua            FlyHuaAnalysis  `json:"fly_hua,omitempty"`
}

type PalaceStar struct {
	Name       string `json:"name"`
	Brightness string `json:"brightness,omitempty"`
}

type TransformData struct {
	Star           string `json:"star"`
	Transformation string `json:"transformation"`
	Display        string `json:"display"`
}

type TemporalPalaceData struct {
	Label      string `json:"label"`
	Branch     string `json:"branch"`
	Palace     string `json:"palace"`
	Stem       string `json:"stem"`
	TimeBranch string `json:"time_branch"`
}

type DaYunData struct {
	Index    int    `json:"index"`
	StartAge int    `json:"start_age"`
	EndAge   int    `json:"end_age"`
	Stem     string `json:"stem"`
	Branch   string `json:"branch"`
	Palace   string `json:"palace"`
}

type PatternData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

type InterpretationData struct {
	Summary              string             `json:"summary"`
	CharacterTraits      string             `json:"character_traits"`
	OriginPalaceAnalysis string             `json:"origin_palace_analysis"`
	KarmicNarrative      []KarmicStep       `json:"karmic_narrative"`
	SanFangDiagnosis     []SanFangRole      `json:"san_fang_diagnosis"`
	StarDetails          []DeepStarAnalysis `json:"star_details,omitempty"`
	OriginFlyHua         FlyHuaAnalysis     `json:"origin_fly_hua,omitempty"`
	TemporalResonance    []ResonancePoint   `json:"temporal_resonance,omitempty"`
	ClassicPatterns      []string           `json:"classic_patterns,omitempty"`
}

type DeepStarAnalysis struct {
	Name       string `json:"name"`
	Verse      string `json:"verse"`
	Positive   string `json:"positive"`
	Negative   string `json:"negative"`
	Remedy     string `json:"remedy"`
	Evolution  string `json:"evolution"`
	Brightness string `json:"brightness"`
}

type FlyHuaAnalysis struct {
	FromPalace string     `json:"from_palace"`
	Stem       string     `json:"stem"`
	Stages     []FlyStage `json:"stages"`
}

type FlyStage struct {
	Type            string          `json:"type"`
	Star            string          `json:"star"`
	Target          string          `json:"target"`
	Motive          string          `json:"motive"`
	Action          string          `json:"action"`
	Trap            string          `json:"trap"`
	Interpretations MultiSchoolView `json:"interpretations"`
}

type MultiSchoolView struct {
	SanHe   string `json:"sanhe"`
	SiHua   string `json:"sihua"`
	QinTian string `json:"qintian"`
}

type ResonancePoint struct {
	Layer  string `json:"layer"`
	Type   string `json:"type"`
	Star   string `json:"star"`
	Natal  string `json:"natal"`
	Palace string `json:"palace"`
	Mood   string `json:"mood"`
}

type KarmicStep struct {
	Type   string `json:"type"`
	Role   string `json:"role"`
	Star   string `json:"star"`
	Palace string `json:"palace"`
	Desc   string `json:"desc"`
}

type SanFangRole struct {
	Role      string `json:"role"`
	Palace    string `json:"palace"`
	Diagnosis string `json:"diagnosis"`
}

type BirthRecord struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Year      int      `json:"year"`
	Month     int      `json:"month"`
	Day       int      `json:"day"`
	Hour      int      `json:"hour"`
	Gender    string   `json:"gender"`
	IsLunar   bool     `json:"is_lunar"`
	IsLeap    bool     `json:"is_leap"`
	IsDST     bool     `json:"is_dst"`
	Longitude float64  `json:"longitude"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
}

type Tag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
