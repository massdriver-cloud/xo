package compliance

import "testing"

func TestSummarizeSARIF(t *testing.T) {
	data := []byte(`{
  "version": "2.1.0",
  "runs": [
    {
      "tool": { "driver": { "name": "checkov", "semanticVersion": "3.2.0" } },
      "results": [
        { "ruleId": "CKV_AWS_16", "level": "error", "message": { "text": "not encrypted" } },
        { "ruleId": "CKV_AWS_17", "level": "warning", "message": { "text": "warn" } },
        { "ruleId": "CKV_AWS_18", "kind": "pass", "message": { "text": "ok" } },
        { "ruleId": "CKV_AWS_19", "level": "error", "message": { "text": "suppressed" },
          "suppressions": [ { "kind": "external", "status": "accepted" } ] }
      ]
    }
  ]
}`)

	scanners, summary, err := SummarizeSARIF(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(scanners) != 1 || scanners[0].Name != "checkov" || scanners[0].Version != "3.2.0" {
		t.Fatalf("unexpected scanners: %+v", scanners)
	}
	if summary.Failed != 2 {
		t.Errorf("expected 2 failed (error + warning findings), got %d", summary.Failed)
	}
	if summary.Passed != 1 {
		t.Errorf("expected 1 passed, got %d", summary.Passed)
	}
	if summary.Skipped != 1 {
		t.Errorf("expected 1 skipped (suppressed), got %d", summary.Skipped)
	}
	if summary.LevelCounts["error"] != 1 {
		t.Errorf("expected 1 error-level finding (suppressed one excluded), got %d", summary.LevelCounts["error"])
	}
	if summary.LevelCounts["warning"] != 1 {
		t.Errorf("expected 1 warning-level finding, got %d", summary.LevelCounts["warning"])
	}
}

func TestSummarizeSARIF_InvalidInput(t *testing.T) {
	if _, _, err := SummarizeSARIF([]byte(`{not sarif`)); err == nil {
		t.Fatal("expected error for invalid input, got nil")
	}
}
