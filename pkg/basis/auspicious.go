package basis

type AuspiciousStar int

const (
	AuspiciousZuofu    AuspiciousStar = iota // 左輔星
	AuspiciousYoubi                          // 右弼星
	AuspiciousWenchang                       // 文昌星
	AuspiciousWenqu                          // 文曲星
	AuspiciousTiankui                        // 天魁星
	AuspiciousTianyue                        // 天鉞星
)

func (s AuspiciousStar) String() string {
	names := []string{
		"左輔", "右弼", "文昌", "文曲", "天魁", "天鉞",
	}
	return names[s]
}

type LuCunStar int

const (
	LuCun  LuCunStar = iota // 祿存星
	Tianma                  // 天馬星
)

func (s LuCunStar) String() string {
	names := []string{"祿存", "天馬"}
	return names[s]
}
