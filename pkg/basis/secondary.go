package basis

type SecondaryStar int

const (
	SecondaryHongluan  SecondaryStar = iota // 紅鸞
	SecondaryTianxi                         // 天喜
	SecondaryGuchen                         // 孤辰
	SecondaryGuaxu                          // 寡宿
	SecondaryLongchi                        // 龍池
	SecondaryFengge                         // 鳳閣
	SecondaryTianxing                       // 天刑
	SecondaryTianyao                        // 天姚
	SecondaryJieshen                        // 解神
	SecondaryTianwu                         // 天巫
	SecondarySantai                         // 三台
	SecondaryBazuo                          // 八座
	SecondaryXianchi                        // 咸池
	SecondaryTianyue                        // 天月
	SecondaryYinsha                         // 陰煞
	SecondaryTianning                       // 台輔
	SecondaryFenggao                        // 封誥
	SecondaryTianku                         // 天哭
	SecondaryTianxu                         // 天虛
	SecondaryHuagai                         // 華蓋
	SecondaryPosui                          // 破碎
	SecondaryFeilian                        // 蜚廉
	SecondaryTianguan                       // 天官
	SecondaryTianfu                         // 天福
	SecondaryEnguang                        // 恩光
	SecondaryTiangui                        // 天貴
	SecondaryTiancai                        // 天才
	SecondaryTianshou                       // 天壽
	SecondaryTiande                         // 天德
	SecondaryYuede                          // 月德
	SecondaryTianshang                      // 天傷
	SecondaryTianshi                        // 天使
	SecondaryTiankong                       // 天空
	SecondaryJiesha                         // 劫煞
)

func (s SecondaryStar) String() string {
	names := []string{
		"紅鸞", "天喜", "孤辰", "寡宿", "龍池", "鳳閣",
		"天刑", "天姚", "解神", "天巫", "三台", "八座",
		"咸池", "天月", "陰煞", "台輔", "封誥",
		"天哭", "天虛", "華蓋", "破碎", "蜚廉",
		"天官", "天福", "恩光", "天貴",
		"天才", "天壽", "天德", "月德", "天傷", "天使",
		"天空", "劫煞",
	}
	return names[s]
}
