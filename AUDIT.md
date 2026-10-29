# Review for the safety fork

> [!WARNING]
> This was written by AI and may contain errors. The findings on the timestamps were interesting with day-light saving time, etc. Codex told me it fixed it in here, but I have not verified it. I just know that the library works for now.

Reviewed 2026-10-29. Upstream snapshot: [`31604653dcd7e8a3201d2b695a05870cd11567c4`](https://github.com/client9/nowandlater/commit/31604653dcd7e8a3201d2b695a05870cd11567c4).

## Scope and conclusion

Inspected the current tracked Go source, dependency manifest, CLI, Makefile, CI, tests, and documentation. Retrieved all 93 upstream commit details and scanned their available patches for risky execution, dependencies, credentials, and workflow changes. Every downloaded upstream file was checked against its pinned Git blob hash before editing.

No malicious behavior was identified in this scope. The runtime module has no third-party dependencies; imports are standard-library packages and this module's own packages. No runtime network access, shell execution, credential reads, unsafe code, cgo, binary payloads, or initialization hooks were found. The CLI reads stdin/arguments and writes stdout/stderr.

The requested functional review and fixes focus on English. This is a source review with bounded dynamic testing, not a formal proof, certification, or guarantee about future upstream versions. A fork preserves control over updates; it does not by itself establish safety. A dedicated vulnerability-database scan was not run because `govulncheck` was unavailable. Standard-library security also depends on the Go toolchain used to build the application.

## Confirmed English findings and disposition

1. **Relative offsets overflowed `time.Duration`.** At reference time `2026-03-22T10:15:30Z`, `in 300 years` returned `1741-06-20T10:40:56Z` without an error. Calendar components are now preserved and resolved using `time.AddDate`, while elapsed seconds are added without nanosecond multiplication. Three hundred years now resolves to `2326-03-22T10:15:30Z`. Relative results outside years 1..9999 return an error. Fractional calendar units use a calendar whole part followed by an approximate elapsed fraction, as documented in the README.

2. **Explicit midnight was treated as an absent clock time.** `tomorrow at midnight` retained the reference clock time; `tomorrow at 00:00:30` also ignored its explicit time. Time combinators now record whether a clock was explicitly supplied, including zero-valued hour/minute/second components. `midnight` was already an English keyword and continues to mean 00:00.

3. **Calendar days were flattened to 86,400 seconds.** In `America/New_York`, from `2026-03-07T23:30:00-05:00`, `tomorrow` skipped March 8. From `2026-11-01T00:30:00-04:00`, its interval still started on November 1. Days and larger whole calendar units now use local calendar arithmetic. Elapsed hours retain elapsed-time semantics. DST day intervals are verified to last 23 or 25 hours as appropriate.

4. **Invalid absolute input is normalized instead of rejected.** Examples: `2026-02-31` becomes March 3; `99:99` rolls forward several days. This behavior is deliberately retained at the user's request. Applications must not treat a successful parse as strict validation of an absolute timestamp.

5. **The CLI's date-only `-now` flag parsed UTC midnight, then converted to local time.** This shifted the intended calendar date in some locations. It now uses `time.ParseInLocation` with `time.Local`.

6. **Ten-digit integer tokens could overflow a 32-bit platform int and panic.** Tokenization now checks `strconv.Atoi` errors instead of using a panic-based conversion for these tokens. Native testing was on darwin/arm64; 32-bit behavior was checked by source inspection rather than executing a 32-bit binary.

## Timezone behavior

The default parser uses Go's `time.Now()` and `time.Local`. A custom `Now` supplies its own location unless `Parser.Location` overrides it. A timezone explicitly written in the input has highest priority. IANA locations loaded with `time.LoadLocation` use Go's timezone database and DST rules. Abbreviations such as `EST` are intentionally fixed offsets. These precedences have regression tests.

Go's `time.AddDate` normalizes month ends, leap days, and nonexistent local times according to Go's documented behavior. Ambiguous or nonexistent wall-clock times during DST transitions are not given a new application-specific disambiguation policy by this fork.

## Build and supply chain

- Retained the original MIT license and attribution; changed the module and internal imports to `github.com/Unbreathable/nowandlater` so users actually depend on the fork.
- The original upstream CI is retained: it tests Go 1.26 and stable, but uses mutable action tags, latest lint tooling, and pull-request write permission. These remain developer/CI supply-chain hardening opportunities. The publishing credential lacks GitHub workflow scope, so the proposed workflow hardening could not be committed. Race tests, vet, and the English fuzz target were run locally.
- Removed the `curl | sh` developer setup recipe. `make env` instead installs a version-pinned Go linter from source; it is optional developer tooling, not a runtime dependency.
- No automated upstream synchronization or automatic releases were added.

## Validation

The published fork tree was compared with the locally tested snapshot: all 64 tracked files matched their Git blob hashes before release preparation.

- Upstream baseline: `go test -race ./...` and `go vet ./...` passed; the upstream fuzz target completed a 30-second run. Passing baseline tests did not detect the semantic issues above.
- Fork: full tests with the race detector; `go vet ./...`; root-package `golangci-lint run .`; English fuzzing for 30 seconds (135,463 executions in the recorded run).
- Added English regression tests for exact 300-year results, historical/strict policies, compound units, excessive offsets, midnight/seconds, spring/fall DST, 23/25-hour intervals, local defaults, explicit timezone precedence, and local CLI date parsing.
- Shared calendar arithmetic also changes four existing non-English test expectations from approximate months/years to calendar months/years. No additional non-English feature work was done.

## Remaining scope limits

Non-English parsing was not hardened in this release. The initial broad check found that malformed Japanese/Chinese time components such as `3時100分` / `3点100分` can panic; that work was explicitly deferred when the requested focus became English. The English-only fuzz target does not certify those tokenizers. Shared language configuration must also remain immutable while parsers use it concurrently.
