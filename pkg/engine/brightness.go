package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func CalcStarBrightness(chart *ZiweiChart) []basis.StarBrightness {
	var result []basis.StarBrightness

	// 計算主星亮度
	for b, stars := range chart.Stars {
		for _, s := range stars {
			brightness := basis.BrightnessLevel(s, b)
			result = append(result, basis.StarBrightness{
				Star:       s,
				Branch:     b,
				Palace:     chart.Palaces[b],
				Brightness: brightness,
			})
		}
	}
	return result
}

// AuspiciousStarBrightness 六吉星亮度資訊
type AuspiciousStarBrightness struct {
	Star       basis.AuspiciousStar
	Branch     basis.Branch
	Palace     basis.Palace
	Brightness basis.Brightness
}

// CalcAuspiciousStarBrightness 計算六吉星亮度
func CalcAuspiciousStarBrightness(chart *ZiweiChart) []AuspiciousStarBrightness {
	var result []AuspiciousStarBrightness

	for b, stars := range chart.AssistantStars {
		for _, s := range stars {
			switch star := s.(type) {
			case basis.AuspiciousStar:
				brightness := basis.AuspiciousBrightnessLevel(star, b)
				result = append(result, AuspiciousStarBrightness{
					Star:       star,
					Branch:     b,
					Palace:     chart.Palaces[b],
					Brightness: brightness,
				})
			case basis.LuCunStar:
				brightness := basis.LuCunBrightnessLevel(star, b)
				result = append(result, AuspiciousStarBrightness{
					Star:       basis.AuspiciousStar(star),
					Branch:     b,
					Palace:     chart.Palaces[b],
					Brightness: brightness,
				})
			}
		}
	}
	return result
}

// MaleficStarBrightness 六煞星亮度資訊
type MaleficStarBrightness struct {
	Star       basis.MaleficStar
	Branch     basis.Branch
	Palace     basis.Palace
	Brightness basis.Brightness
}

// CalcMaleficStarBrightness 計算六煞星亮度
func CalcMaleficStarBrightness(chart *ZiweiChart) []MaleficStarBrightness {
	var result []MaleficStarBrightness

	for b, stars := range chart.AssistantStars {
		for _, s := range stars {
			if star, ok := s.(basis.MaleficStar); ok {
				brightness := basis.MaleficBrightnessLevel(star, b)
				result = append(result, MaleficStarBrightness{
					Star:       star,
					Branch:     b,
					Palace:     chart.Palaces[b],
					Brightness: brightness,
				})
			}
		}
	}
	return result
}
