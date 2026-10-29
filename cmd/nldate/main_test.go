package main

import (
	"testing"
	"time"
)

func TestDateOnlyNowUsesLocalMidnight(t *testing.T) {
	oldLocal := time.Local
	defer func() { time.Local = oldLocal }()
	for _, offset := range []int{-7 * 3600, 9 * 3600} {
		time.Local = time.FixedZone("test-local", offset)
		got, err := parseNowFlag("2026-03-22")
		want := time.Date(2026, 3, 22, 0, 0, 0, 0, time.Local)
		if err != nil || !got.Equal(want) || got.Location() != time.Local {
			t.Fatalf("offset %d: got %v, %v; want %v", offset, got, err, want)
		}
	}
}
