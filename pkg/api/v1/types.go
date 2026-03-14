package v1

type ZiweiResponse struct {
	Gender   string                `json:"gender"`
	Wuxing   string                `json:"wuxing"`
	NaYin    string                `json:"na_yin"`
	MingGong string                `json:"ming_gong"`
	ShenGong string                `json:"shen_gong"`
	Palaces  map[string]PalaceData `json:"palaces"`
}

type PalaceData struct {
	Branch string   `json:"branch"`
	Stars  []string `json:"stars"`
}
