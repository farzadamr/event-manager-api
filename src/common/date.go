package common

import (
	"fmt"
	"math"
	"time"

	persian "github.com/yaa110/go-persian-calendar"
)

func ToShamsiString(t time.Time) string {
	p := persian.New(t)
	return fmt.Sprintf(
		"%04d/%02d/%02d",
		p.Year(), p.Month(), p.Day())
}

func CalculateDurationString(start, end time.Time) string {
	if end.Before(start) {
		start, end = end, start
	}
	hours := end.Sub(start).Hours()

	h := math.Round(hours*100) / 10
	if h == float64(int64(h)) {
		return fmt.Sprintf("%d", int64(h))
	}
	return fmt.Sprintf("%.1f", h)
}
