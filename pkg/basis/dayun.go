package basis

type DaYun struct {
	Index    int
	StartAge int
	EndAge   int
	Stem     Stem
	Branch   Branch
	Palace   Palace
}

func (d DaYun) String() string {
	return d.Branch.String()
}

type LiuNian struct {
	Year   int
	Stem   Stem
	Branch Branch
	Palace Palace
}

func (l LiuNian) String() string {
	return l.Branch.String()
}
