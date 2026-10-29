# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-10-29

### Added

- Independently controlled fork of [client9/nowandlater](https://github.com/client9/nowandlater), based on upstream commit `31604653dcd7e8a3201d2b695a05870cd11567c4`.
- English regression tests for local timezone defaults, DST, midnight, calendar units, and large offsets; an English-only fuzz target.
- Audit scope and findings in `AUDIT.md`.

### Fixed

- Calendar days, weeks, months, and years now use Go's `time.AddDate`; hours, minutes, and seconds remain elapsed units.
- Large relative offsets no longer overflow `time.Duration`; relative results outside years 1..9999 return errors.
- Explicit midnight and zero-hour times retain their intended clock time in relative expressions.
- The CLI's date-only `-now` option now uses local midnight through `time.ParseInLocation`.
- Numeric tokens that overflow the platform int now return a parse error rather than panicking.

### Changed

- Go module path is `github.com/Unbreathable/nowandlater`.
- Developer setup installs a pinned Go tool instead of piping a downloaded shell script to `sh`.
- Invalid absolute dates and clock values retain upstream's permissive normalization.

- Reorder structs to save memory using [betteralign](https://github.com/dkorunic/betteralign/)

## [0.9.0] - 2026-04-28

Initial public pre-release. The library API is considered stable; documentation
and internal review are still in progress ahead of v1.0.0.

### Added

- `Parser` type with `Parse` and `ParseInterval` methods for converting
  natural-language date/time strings into `time.Time` values and
  `[start, end)` calendar intervals.
- Zero-dependency implementation using only the Go standard library (no
  regular expressions, no parser generators).
- Configurable `Parser` fields: `Lang`, `Location`, `Now`, and `DateOrder`.
  The zero value is valid and defaults to English, `time.Local`, and
  `time.Now`.
- `Parser` is safe for concurrent use.
- Language support for English, Spanish, French, German, Italian,
  Portuguese, Russian, Japanese, and Chinese, exposed as `LangEn`,
  `LangEs`, `LangFr`, `LangDe`, `LangIt`, `LangPt`, `LangRu`, `LangJa`,
  and `LangZh`.
- Package-level language lookup helpers in `lookup.go`.
- Absolute date formats: ISO 8601, alternative separators
  (`/`, `.`), month-name forms, and locale-sensitive `MM/DD/YYYY` /
  `DD/MM/YYYY` parsing via `DateOrder`.
- Relative expressions such as `next Monday`, `in 3 days`, `3 days ago`,
  `this week`, and `last month`.
- `cmd/nldate` CLI for ad-hoc parsing and debugging
  (`go run ./cmd/nldate "next Monday"`).
- Benchmarks (`bench_test.go`) and a fuzz/test suite with >95% coverage.
- Package documentation in `doc.go` and contributor docs under `docs/`.

[Unreleased]: https://github.com/Unbreathable/nowandlater/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Unbreathable/nowandlater/compare/31604653dcd7e8a3201d2b695a05870cd11567c4...v1.0.0
[0.9.0]: https://github.com/client9/nowandlater/releases/tag/v0.9.0
