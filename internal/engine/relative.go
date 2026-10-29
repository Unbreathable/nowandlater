package engine

import (
	"fmt"
	"math"
	"time"
)

// Keep calendar units separate from elapsed seconds: a local day need not last
// 24 hours, and time.Duration cannot represent offsets beyond about 292 years.
type relativeDelta struct {
	years, months, days int
	seconds             int64
}

// Relative results are limited to years 1..9999, the RFC 3339 year range.
// This also bounds all arithmetic before converting counts to machine ints.
const maxRelativeSeconds = 10000 * 366 * 86400

func relativeSlots(n float64, period Period, sign int) (*ParsedDateSlots, error) {
	unit := periodToSeconds[period]
	if unit == 0 || math.IsNaN(n) || math.IsInf(n, 0) || n < 0 || n > float64(maxRelativeSeconds)/float64(unit) {
		return nil, fmt.Errorf("nowandlater: relative offset out of range")
	}
	whole, fraction := math.Modf(n)
	delta := &relativeDelta{}
	switch period {
	case PeriodYear:
		delta.years = int(whole) * sign
	case PeriodMonth:
		delta.months = int(whole) * sign
	case PeriodWeek:
		delta.days = int(whole) * 7 * sign
	case PeriodFortnight:
		delta.days = int(whole) * 14 * sign
	case PeriodDay:
		delta.days = int(whole) * sign
	default:
		fraction = n
	}
	delta.seconds = int64(math.Round(fraction*float64(unit))) * int64(sign)
	slots := &ParsedDateSlots{relative: delta, Period: period}
	setLegacySeconds(slots, math.Round(n*float64(unit))*float64(sign))
	return slots, nil
}

// DeltaSeconds remains an approximate projection for existing engine callers.
// When the platform int cannot hold it, the precise relative delta still works.
func setLegacySeconds(slots *ParsedDateSlots, seconds float64) {
	n := 0
	if seconds >= -float64(int(^uint(0)>>1)) && seconds <= float64(int(^uint(0)>>1)) {
		n = int(seconds)
	}
	slots.DeltaSeconds = new(n)
}

func (d relativeDelta) apply(now time.Time, sign int) (time.Time, error) {
	base := now.AddDate(d.years*sign, d.months*sign, d.days*sign)
	// Add whole seconds through Unix time rather than nanosecond Duration.
	base = time.Unix(base.Unix()+d.seconds*int64(sign), int64(base.Nanosecond())).In(now.Location())
	if base.Year() < 1 || base.Year() > 9999 {
		return time.Time{}, fmt.Errorf("nowandlater: relative result outside years 1..9999")
	}
	return base, nil
}
