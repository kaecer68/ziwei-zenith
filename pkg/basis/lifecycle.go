package basis

type ChangShengStage int

const (
	ChangSheng ChangShengStage = iota
	MuYu
	GuanDai
	LinGuan
	DiWang
	Shuai
	Bing
	Si
	Mu
	Jue
	Tai
	Yang
)

func (c ChangShengStage) String() string {
	names := []string{"長生", "沐浴", "冠帶", "臨官", "帝旺", "衰", "病", "死", "墓", "絕", "胎", "養"}
	return names[c]
}

type BoShiStar int

const (
	BoShi BoShiStar = iota
	LiShi
	QingLong
	XiaoHao
	JiangJun
	ZouShu
	FeiLian
	XiShen
	BingFu
	DaHao
	FuBing
	GuanFu
)

func (b BoShiStar) String() string {
	names := []string{"博士", "力士", "青龍", "小耗", "將軍", "奏書", "飛廉", "喜神", "病符", "大耗", "伏兵", "官符"}
	return names[b]
}
