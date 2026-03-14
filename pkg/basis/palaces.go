package basis

type Palace int

const (
	PalaceMing     Palace = iota // 命宮
	PalaceXiongDi                // 兄弟宮
	PalaceFuQi                   // 夫妻宮
	PalaceZiNv                   // 子女宮
	PalaceCaiBang                // 財帛宮
	PalaceJiE                    // 疾厄宮
	PalaceQianYi                 // 遷移宮
	PalacePuYi                   // 僕役宮
	PalaceGuanLu                 // 官祿宮
	PalaceTianZhai               // 田宅宮
	PalaceFuDe                   // 福德宮
	PalaceFuMu                   // 父母宮
)

func (p Palace) String() string {
	names := []string{
		"命宮", "兄弟宮", "夫妻宮", "子女宮",
		"財帛宮", "疾厄宮", "遷移宮", "僕役宮",
		"官祿宮", "田宅宮", "福德宮", "父母宮",
	}
	return names[p]
}

func (p Palace) Index() int {
	return int(p)
}

var PalaceOrder = []Palace{
	PalaceMing, PalaceXiongDi, PalaceFuQi, PalaceZiNv,
	PalaceCaiBang, PalaceJiE, PalaceQianYi, PalacePuYi,
	PalaceGuanLu, PalaceTianZhai, PalaceFuDe, PalaceFuMu,
}

func PalaceFromIndex(idx int) Palace {
	return PalaceOrder[idx%12]
}

func (p Palace) Opposite() Palace {
	return Palace((int(p) + 6) % 12)
}

func (p Palace) Next(n int) Palace {
	return Palace((int(p) + n) % 12)
}

func (p Palace) Prev(n int) Palace {
	return Palace((int(p) - n + 12) % 12)
}
