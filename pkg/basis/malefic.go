package basis

type MaleficStar int

const (
	MaleficQingyang MaleficStar = iota // 擎羊星
	MaleficTuoluo                      // 陀羅星
	MaleficHuoxing                     // 火星
	MaleficLingxing                    // 鈴星
	MaleficDikong                      // 地空星
	MaleficDijie                       // 地劫星
)

func (s MaleficStar) String() string {
	names := []string{
		"擎羊", "陀羅", "火星", "鈴星", "地空", "地劫",
	}
	return names[s]
}
