package compliance

import (
	"fmt"

	sariflib "github.com/owenrumney/go-sarif/v3/pkg/report/v210/sarif"
)

// SummarizeSARIF parses SARIF and returns the scanners that produced it plus a
// suppression-aware summary. Counts are bucketed by the standard SARIF result
// level (error / warning / note); tool-specific severity scores are not
// interpreted here.
func SummarizeSARIF(data []byte) ([]Scanner, Summary, error) {
	report, err := sariflib.FromBytes(data)
	if err != nil {
		return nil, Summary{}, fmt.Errorf("failed to parse SARIF: %w", err)
	}

	var scanners []Scanner
	summary := Summary{LevelCounts: map[string]int{}}

	for _, run := range report.Runs {
		if run.Tool != nil && run.Tool.Driver != nil {
			scanners = append(scanners, scannerFromDriver(run.Tool.Driver))
		}
		for _, result := range run.Results {
			switch {
			case isSuppressed(result):
				summary.Skipped++
			case isPass(result):
				summary.Passed++
			case isSkip(result):
				summary.Skipped++
			default: // a finding
				summary.Failed++
				summary.LevelCounts[resultLevel(result)]++
			}
		}
	}

	if len(summary.LevelCounts) == 0 {
		summary.LevelCounts = nil
	}
	return scanners, summary, nil
}

// isSuppressed reports whether a result carries an effective suppression. A
// suppression with status "rejected" does not suppress.
func isSuppressed(r *sariflib.Result) bool {
	for _, s := range r.Suppressions {
		if s.Status == nil || *s.Status != "rejected" {
			return true
		}
	}
	return false
}

func isPass(r *sariflib.Result) bool {
	return r.Kind == "pass"
}

// isSkip covers SARIF kinds that are neither pass nor an actionable finding.
func isSkip(r *sariflib.Result) bool {
	switch r.Kind {
	case "notApplicable", "informational", "open", "review":
		return true
	default:
		return false
	}
}

// resultLevel returns the SARIF level, defaulting to "warning" when unspecified
// (the SARIF default).
func resultLevel(r *sariflib.Result) string {
	if r.Level == "" {
		return "warning"
	}
	return r.Level
}

func scannerFromDriver(d *sariflib.ToolComponent) Scanner {
	s := Scanner{}
	if d.Name != nil {
		s.Name = *d.Name
	}
	switch {
	case d.SemanticVersion != nil:
		s.Version = *d.SemanticVersion
	case d.Version != nil:
		s.Version = *d.Version
	}
	return s
}
