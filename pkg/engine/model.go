package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

type BirthInfo struct {
	LunarYear   int
	LunarMonth  int
	LunarDay    int
	HourBranch  basis.HourBranch
	YearPillar  basis.Pillar
	MonthPillar basis.Pillar
	DayPillar   basis.Pillar
	HourPillar  basis.Pillar
	Sex         basis.Sex
}
