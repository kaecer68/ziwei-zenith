package basis

type SecondaryStar int

const (
	SecondaryHongluan SecondaryStar = iota // 紅鸞
	SecondaryTianxi                        // 天喜
	SecondaryGuchen                        // 孤辰
	SecondaryGuaxu                         // 寡宿
	SecondaryLongchi                       // 龍池
	SecondaryFengge                        // 鳳閣
	SecondaryTianxing                      // 天刑
	SecondaryTianyao                       // 天姚
	SecondaryJieshen                       // 解神
	SecondaryTianwu                        // 天巫
	SecondarySantai                        // 三台
	SecondaryBazuo                         // 八座
)

func (s SecondaryStar) String() string {
	names := []string{
		"紅鸞", "天喜", "孤辰", "寡宿", "龍池", "鳳閣",
		"天刑", "天姚", "解神", "天巫", "三台", "八座",
	}
	return names[s]
}
