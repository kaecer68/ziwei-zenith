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
	Palaces        map[string]PalaceData `json:"palaces"`
	Patterns       []PatternData         `json:"patterns"`
	Interpretation InterpretationData    `json:"interpretation"`
}

type PalaceData struct {
	Branch       string   `json:"branch"`
	Stars        []string `json:"stars"`
	LiuNianStars []string `json:"liu_nian_stars,omitempty"`
	LiuYueStars  []string `json:"liu_yue_stars,omitempty"`
	LiuRiStars   []string `json:"liu_ri_stars,omitempty"`
}

type PatternData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
}

type InterpretationData struct {
	Summary           string           `json:"summary"`
	KarmicNarrative   []KarmicStep     `json:"karmic_narrative"`
	SanFangDiagnosis  []SanFangRole    `json:"san_fang_diagnosis"`
	TemporalResonance []ResonancePoint `json:"temporal_resonance,omitempty"`
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
