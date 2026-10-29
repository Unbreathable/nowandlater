package nowandlater_test

import (
	"testing"
	"time"

	nl "github.com/Unbreathable/nowandlater"
)

func TestRelativeCalendarArithmetic(t *testing.T) {
	now := time.Date(2026, 3, 22, 10, 15, 30, 123, time.UTC)
	p := nl.Parser{Now: func() time.Time { return now }}
	for _, tc := range []struct {
		input string
		want  time.Time
	}{
		{"in 300 years", now.AddDate(300, 0, 0)},
		{"300 years ago", now.AddDate(-300, 0, 0)},
		{"in 300.5 years", now.AddDate(300, 0, 0).Add(182*24*time.Hour + 12*time.Hour)},
		{"in 1 year and 6 months", now.AddDate(1, 6, 0)},
		{"1 year and 6 months ago", now.AddDate(-1, -6, 0)},
		{"in 1 month", now.AddDate(0, 1, 0)},
		{"in 9999999999 seconds", time.Unix(now.Unix()+9999999999, int64(now.Nanosecond()))},
		{"3 days before tomorrow", now.AddDate(0, 0, -2)},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, err := p.Parse(tc.input)
			if err != nil || !got.Equal(tc.want) {
				t.Fatalf("got %v, %v; want %v", got, err, tc.want)
			}
		})
	}
	p.Ambiguity = nl.AmbiguityHistorical
	got, err := p.Parse("300 years")
	if err != nil || !got.Equal(now.AddDate(-300, 0, 0)) {
		t.Fatalf("historical duration: %v, %v", got, err)
	}
	p.Ambiguity = nl.AmbiguityStrict
	if _, err := p.Parse("300 years"); err == nil {
		t.Fatal("strict mode accepted ambiguous duration")
	}
}

func TestRelativeRange(t *testing.T) {
	now := time.Date(2026, 3, 22, 10, 0, 0, 0, time.UTC)
	p := nl.Parser{Now: func() time.Time { return now }}
	for _, input := range []string{"in 9999999999 years", "in 9999999999.9999999999 years", "in 9999999999 months", "in 9999999999 days", "in 8000 years", "3000 years ago", "in 5000 years and 5000 years"} {
		t.Run(input, func(t *testing.T) {
			if got, err := p.Parse(input); err == nil {
				t.Errorf("accepted out-of-range result %v", got)
			}
			if _, _, err := p.ParseInterval(input); err == nil {
				t.Error("interval accepted out-of-range result")
			}
		})
	}
	for _, year := range []int{1, 9999} {
		boundary := time.Date(year, 6, 15, 12, 0, 0, 0, time.UTC)
		p.Now = func() time.Time { return boundary }
		got, err := p.Parse("now")
		if err != nil || !got.Equal(boundary) {
			t.Fatalf("year %d should work: %v, %v", year, got, err)
		}
	}
}

func TestExplicitMidnight(t *testing.T) {
	now := time.Date(2026, 3, 22, 10, 15, 30, 123, time.UTC)
	p := nl.Parser{Now: func() time.Time { return now }}
	for _, tc := range []struct {
		input       string
		day, second int
	}{
		{"midnight", 22, 0}, {"today at midnight", 22, 0},
		{"tomorrow at midnight", 23, 0}, {"tomorrow at 12 AM", 23, 0},
		{"tomorrow at 00:00:30", 23, 30}, {"yesterday at midnight", 21, 0},
		{"next Monday at midnight", 23, 0},
	} {
		t.Run(tc.input, func(t *testing.T) {
			want := time.Date(2026, 3, tc.day, 0, 0, tc.second, 0, time.UTC)
			got, err := p.Parse(tc.input)
			if err != nil || !got.Equal(want) {
				t.Fatalf("got %v, %v; want %v", got, err, want)
			}
			start, _, err := p.ParseInterval(tc.input)
			if err != nil || !start.Equal(want) {
				t.Fatalf("interval start %v, %v; want %v", start, err, want)
			}
		})
	}
	got, err := p.Parse("in 2 hours")
	if err != nil || !got.Equal(now.Add(2*time.Hour)) {
		t.Fatalf("duration lost its clock time: %v, %v", got, err)
	}
}

func TestRelativeLocalDST(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatal(err)
	}
	for _, now := range []time.Time{
		time.Date(2026, 3, 7, 23, 30, 0, 0, loc),
		time.Date(2026, 11, 1, 0, 30, 0, 0, loc),
	} {
		t.Run(now.Format(time.RFC3339), func(t *testing.T) {
			p := nl.Parser{Now: func() time.Time { return now }}
			for _, tc := range []struct {
				input string
				want  time.Time
			}{
				{"tomorrow", now.AddDate(0, 0, 1)},
				{"in 1 day", now.AddDate(0, 0, 1)},
				{"in 1 week", now.AddDate(0, 0, 7)},
				{"yesterday", now.AddDate(0, 0, -1)},
				{"in 24 hours", now.Add(24 * time.Hour)},
				{"3 days before tomorrow", now.AddDate(0, 0, -2)},
			} {
				got, err := p.Parse(tc.input)
				if err != nil || !got.Equal(tc.want) || got.Location() != loc {
					t.Errorf("%s: got %v, %v; want %v", tc.input, got, err, tc.want)
				}
			}
			start, _, err := p.ParseInterval("tomorrow")
			want := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, loc)
			if err != nil || !start.Equal(want) {
				t.Fatalf("interval start %v, %v; want %v", start, err, want)
			}
		})
	}
	for _, tc := range []struct {
		now   time.Time
		hours int
	}{
		{time.Date(2026, 3, 7, 12, 0, 0, 0, loc), 23},
		{time.Date(2026, 10, 31, 12, 0, 0, 0, loc), 25},
	} {
		p := nl.Parser{Now: func() time.Time { return tc.now }}
		start, end, err := p.ParseInterval("tomorrow")
		if err != nil || end.Sub(start) != time.Duration(tc.hours)*time.Hour {
			t.Fatalf("DST day length: [%v,%v), %v", start, end, err)
		}
	}
}

func TestDefaultAndExplicitLocations(t *testing.T) {
	// The zero-value parser uses the process's configured local timezone.
	before := time.Now()
	got, err := (nl.Parser{}).Parse("now")
	after := time.Now()
	if err != nil || got.Location() != time.Local || got.Before(before) || got.After(after) {
		t.Fatalf("local defaults: %v, %v", got, err)
	}
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 3, 22, 23, 30, 0, 0, time.UTC)
	p := nl.Parser{Location: loc, Now: func() time.Time { return now }}
	got, err = p.Parse("tomorrow at midnight")
	want := time.Date(2026, 3, 24, 0, 0, 0, 0, loc)
	if err != nil || !got.Equal(want) || got.Location() != loc {
		t.Fatalf("configured location: %v, %v; want %v", got, err, want)
	}
	got, err = p.Parse("tomorrow at midnight UTC")
	want = time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	if err != nil || !got.Equal(want) || got.Location() != time.UTC {
		t.Fatalf("explicit input location: %v, %v; want %v", got, err, want)
	}
}
