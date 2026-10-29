package nowandlater_test

import (
	"testing"
	"time"

	nl "github.com/Unbreathable/nowandlater"
)

func FuzzEnglishParser(f *testing.F) {
	for _, s := range []string{"next Monday at 9:30 AM", "tomorrow at midnight", "2026-02-31", "in 300 years", "in 9999999999 years", "99:99", "12:00 PM", "March 5 at 3 PM", "2026-12-04T09:30:00Z", "9999999999", "1:30am 29/02/20", "a week ago"} {
		f.Add(s)
	}
	now := time.Date(2026, 3, 22, 10, 15, 30, 0, time.UTC)
	p := nl.Parser{Now: func() time.Time { return now }}
	f.Fuzz(func(t *testing.T, s string) {
		p.Parse(s)
		p.ParseInterval(s)
	})
}
