package engine

// StarEssence holds the philosophical and diagnostic data for a star
type StarEssence struct {
	AncientVerse string
	Positive     string
	Negative     string
	Remedy       string
	Trait        string // Summary trait (e.g. "領袖力")
	Shadow       string // Primary negative (e.g. "孤高")
	Core         string // Execution strategy
}

// OriginStory holds the XVIII Flying Stars origin data
type OriginStory struct {
	OriginalName string
	Essence      string
}

var StarEssenceTable = map[string]StarEssence{
	"紫微": {
		AncientVerse: "紫微原屬土，官祿宮主星。有相為有用，無相為孤君。",
		Positive:     "領袖星，善於協調，具有剛強表顯的領導統御能力。",
		Negative:     "孤高自賞、耳根軟、好大喜功、眼高手低。",
		Remedy:       "身段要特別柔軟，必須要特別放下身段、收斂光芒。",
		Trait:        "領袖統御與承擔",
		Shadow:       "高冷與狹隘",
		Core:         "協調與授權",
	},
	"天機": {
		AncientVerse: "天機兄弟主，南斗正曜星。做事有操略，稟性最高明。",
		Positive:     "智謀過人、反應敏捷、善於策劃與點子創新。",
		Negative:     "神經質、想太多做太少、容易焦慮、只相信自己。",
		Remedy:       "學會信任他人，學會做減法，享受簡單。",
		Trait:        "智謀策劃與變動",
		Shadow:       "焦慮與游移",
		Core:         "邏輯與預判",
	},
	"太陽": {
		AncientVerse: "太陽原屬火，正主官祿星。權柄能操縱，榮華各處名。",
		Positive:     "無私奉獻、光芒四射、熱情磊落、照顧他人。",
		Negative:     "燃燒自己照亮別人而過勞、愛面子、強加好意。",
		Remedy:       "發光也要有底線，要知道如何給自己補充能量。",
		Trait:        "光芒散發與奉獻",
		Shadow:       "虛榮與乾涸",
		Core:         "公開與透明",
	},
	"武曲": {
		AncientVerse: "武曲北斗第六星，司權司祿星中尊。",
		Positive:     "正財星，具有剛毅果決、執行力強、堅忍不拔。",
		Negative:     "冷酷無情、過度現實、固執己見、缺乏圓融。",
		Remedy:       "化剛為柔，學會借力使力，凡事不能只看利益。",
		Trait:        "執行力與財富敏銳",
		Shadow:       "孤克與短視",
		Core:         "標準與效率",
	},
	"天同": {
		AncientVerse: "天同南斗第四星，化福為祥樂太平。",
		Positive:     "知足常樂、福氣綿長、協調力強、與人為善。",
		Negative:     "好逸惡勞、逃避壓力、缺乏企圖心、依賴性強。",
		Remedy:       "福氣會用盡的，避免溫水煮青蛙，建立抗壓性。",
		Trait:        "和諧與福氣轉化",
		Shadow:       "墮性與逃避",
		Core:         "情緒調節與社交",
	},
	"廉貞": {
		AncientVerse: "廉貞屬火，化囚為殺。廉合巳亥宮，遇吉福盈豐。",
		Positive:     "外交手腕、公關魅力、第六感強、節度自律。",
		Negative:     "心高氣傲、畫地自限、爛桃花、情緒極端反覆。",
		Remedy:       "從宗教的角度看情感，勿在私情糾纏，換位思考。",
		Trait:        "自律與公關魅力",
		Shadow:       "偏激與執拗",
		Core:         "規範與社交邊界",
	},
	"天府": {
		AncientVerse: "天府為祿庫，入命終是富。萬頃置田莊，家資無論數。",
		Positive:     "穩重保守、善於理財、包容性強、厚道。",
		Negative:     "缺乏承擔、保守過頭、喜歡發號施令卻不想動手。",
		Remedy:       "擴大包容胸襟，將資源投資向未來而非死守現狀。",
		Trait:        "穩重管理與庫存",
		Shadow:       "保守與傲慢",
		Core:         "儲存與資源分配",
	},
	"太陰": {
		AncientVerse: "太陰之星屬水，為田宅之主宰，司一生之財帛。",
		Positive:     "溫柔細膩、蓄財持家、文靜重感情、母性光輝。",
		Negative:     "多愁善感、潔癖、內心戲太多、悲觀傾向。",
		Remedy:       "走出來迎向陽光，允許自己不完美，主動溝通。",
		Trait:        "細膩情感與守成",
		Shadow:       "抑鬱與封閉",
		Core:         "環境適應與沉澱",
	},
	"貪狼": {
		AncientVerse: "貪狼北斗第一星，化氣為桃花。主禍福，亦主解厄。",
		Positive:     "多才多藝、交際手腕強、求知欲旺盛、充滿生命力。",
		Negative:     "慾望太深、博而不精、見異思遷、交友氾濫。",
		Remedy:       "慾望需要聚焦、適中，將貪念鎖在一個專業點上。",
		Trait:        "慾望驅動與才藝",
		Shadow:       "貪婪與不專",
		Core:         "變通與資源博弈",
	},
	"巨門": {
		AncientVerse: "巨門屬水，北斗第二星，化氣為暗。主是非，需磨礪而後明。",
		Positive:     "精準剖析、深度研究、口才極佳、找出漏洞。",
		Negative:     "口舌是非、尖酸刻薄、容易招來暗箭、挑剔過度。",
		Remedy:       "挑剔眼光用在『事』上，讚美言語用在『人』上。",
		Trait:        "研究剖析與口才",
		Shadow:       "是非與屏蔽",
		Core:         "分析與問題發現",
	},
	"天相": {
		AncientVerse: "天相原屬水，化印主官祿。身命二宮逢，定主多財福。",
		Positive:     "輔佐長才、熱心公益、注重形象、協調性佳。",
		Negative:     "失去主見、愛管閒事引火上身。",
		Remedy:       "幫人要量力而為，確立底線與核心價值。",
		Trait:        "輔佐、信譽與邏輯",
		Shadow:       "隨波逐流與多事",
		Core:         "審核與平衡",
	},
	"天梁": {
		AncientVerse: "天梁原屬土，南斗最吉星。化蔭名延壽，父母宮主星。",
		Positive:     "老成持重、逢凶化吉、庇蔭後進、名士風度。",
		Negative:     "喜歡說教、倚老賣老、孤芳自賞、招惹麻煩。",
		Remedy:       "以身作則勝過千言萬語，主動應掉化吉的業力。",
		Trait:        "庇蔭、化解與老成",
		Shadow:       "孤傲與古板",
		Core:         "傳承與調和",
	},
	"七殺": {
		AncientVerse: "七殺寅申子午宮，四夷拱手服英雄。主權柄，亦主孤高。",
		Positive:     "衝鋒陷陣、大將之風、理智堅定、不怕困難。",
		Negative:     "莽撞行事、一生動盪辛勞、缺乏感性。",
		Remedy:       "戰將需要君王與軍師，學會授權與信任團隊。",
		Trait:        "果決行動與孤高",
		Shadow:       "衝擊與肅殺",
		Core:         "決斷與執行",
	},
	"破軍": {
		AncientVerse: "破軍北斗第七星，化氣曰耗。子午破軍，加官進祿。",
		Positive:     "破舊立新、顛覆格局、勇於改革、創意十足。",
		Negative:     "消耗資源、三分鐘熱度、叛逆不羈。",
		Remedy:       "不能衝動，懂得將力量花在對的事情上。",
		Trait:        "變革、損耗與創新",
		Shadow:       "破壞與動盪",
		Core:         "突破與更新",
	},
}

var HistoricalOriginTable = map[string]OriginStory{
	"廉貞": {OriginalName: "天杖", Essence: "原始特質為『刑杖之氣』，帶動法律紛爭與禁錮能量。"},
	"破軍": {OriginalName: "旄頭", Essence: "原始特質為『消耗與破壞』，帶有衝鋒陷陣後的損耗特質。"},
	"七殺": {OriginalName: "天刃", Essence: "原始特質為『利刃之災』，將凶氣化為理智的殺伐決斷。"},
	"天相": {OriginalName: "天印", Essence: "原始特質為『印信權限』。天相是蓋章的人，代表合約與信譽。"},
}
