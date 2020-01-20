package convert

import (
	"fmt"
	"math"
)

func humanize(s uint64) string {
	base := float64(1024)
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	logn := func(n, b float64) float64 {
		return math.Log(n) / math.Log(b)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}
