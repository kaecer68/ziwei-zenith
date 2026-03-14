package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func CalcDaYun(mingGong basis.Palace, sex basis.Sex, yearStem basis.Stem, birthYear int) []basis.DaYun {
	var dayunStemOffsets []int
	var dayunBranchOffsets []int

	isYang := (int(yearStem) % 2) == 0

	if sex == basis.SexMale {
		if isYang {
			dayunStemOffsets = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
			dayunBranchOffsets = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0}
		} else {
			dayunStemOffsets = []int{5, 6, 7, 8, 9, 0, 1, 2, 3, 4}
			dayunBranchOffsets = []int{5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3, 4}
		}
	} else {
		if isYang {
			dayunStemOffsets = []int{5, 6, 7, 8, 9, 0, 1, 2, 3, 4}
			dayunBranchOffsets = []int{5, 6, 7, 8, 9, 10, 11, 0, 1, 2, 3, 4}
		} else {
			dayunStemOffsets = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
			dayunBranchOffsets = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0}
		}
	}

	mingGongIdx := int(mingGong)

	dayuns := make([]basis.DaYun, 0)
	for i := 0; i < 10; i++ {
		palaceIdx := (mingGongIdx + dayunBranchOffsets[i]) % 12

		dayun := basis.DaYun{
			Index:    i + 1,
			StartAge: i * 10,
			EndAge:   (i + 1) * 10,
			Stem:     basis.Stem((int(yearStem) + dayunStemOffsets[i]) % 10),
			Branch:   basis.Branch(dayunBranchOffsets[i]),
			Palace:   basis.Palace(palaceIdx),
		}
		dayuns = append(dayuns, dayun)
	}

	_ = birthYear

	return dayuns
}

func CalcLiuNian(mingGong basis.Palace, dayPillarStem basis.Stem, currentYear int) []basis.LiuNian {
	baseStem := int(dayPillarStem)
	baseBranch := 0

	mingGongIdx := int(mingGong)

	liunians := make([]basis.LiuNian, 0)
	for i := 0; i < 12; i++ {
		year := currentYear - 5 + i

		stemIdx := (baseStem + i) % 10
		branchIdx := (baseBranch + i) % 12
		palaceIdx := (mingGongIdx + branchIdx) % 12

		liunian := basis.LiuNian{
			Year:   year,
			Stem:   basis.Stem(stemIdx),
			Branch: basis.Branch(branchIdx),
			Palace: basis.Palace(palaceIdx),
		}
		liunians = append(liunians, liunian)
	}

	return liunians
}

func CalcLiuYue(mingGong basis.Palace, yearBranch basis.Branch, month int) []basis.LiuYue {
	baseBranch := int(yearBranch)

	mingGongIdx := int(mingGong)

	liuyues := make([]basis.LiuYue, 0)
	for i := 0; i < 12; i++ {
		monthNum := ((month - 1 + i) % 12) + 1

		stemIdx := (baseBranch + i) % 10
		branchIdx := (baseBranch + i) % 12
		palaceIdx := (mingGongIdx + branchIdx) % 12

		liuyue := basis.LiuYue{
			Month:  monthNum,
			Stem:   basis.Stem(stemIdx),
			Branch: basis.Branch(branchIdx),
			Palace: basis.Palace(palaceIdx),
		}
		liuyues = append(liuyues, liuyue)
	}

	return liuyues
}

func CalcLiuRi(mingGong basis.Palace, dayPillarStem basis.Stem, day int) []basis.LiuRi {
	baseStem := int(dayPillarStem)

	mingGongIdx := int(mingGong)

	liuris := make([]basis.LiuRi, 0)
	for i := 0; i < 30; i++ {
		dayNum := ((day - 1 + i) % 30) + 1

		stemIdx := (baseStem + i) % 10
		branchIdx := (i) % 12
		palaceIdx := (mingGongIdx + branchIdx) % 12

		liuri := basis.LiuRi{
			Day:    dayNum,
			Stem:   basis.Stem(stemIdx),
			Branch: basis.Branch(branchIdx),
			Palace: basis.Palace(palaceIdx),
		}
		liuris = append(liuris, liuri)
	}

	return liuris
}
