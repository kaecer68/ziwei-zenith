package engine

import (
	"github.com/kaecer68/ziwei-zenith/pkg/basis"
)

func CalcStarBrightness(chart *ZiweiChart) []basis.StarBrightness {
	var result []basis.StarBrightness

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
