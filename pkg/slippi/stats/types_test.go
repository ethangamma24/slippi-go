package stats_test

import (
	"encoding/json"
	"strings"
	"testing"

	stats "github.com/ethangamma24/slippi-go/pkg/slippi/stats"
)

func TestStatsZeroValue(t *testing.T) {
	s := stats.Stats{
		ActionCounts: []stats.ActionCounts{{}},
	}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"attackCount"`) {
		t.Fatalf("expected attackCount in JSON, got %s", string(b))
	}
	if !strings.Contains(string(b), `"jab1"`) {
		t.Fatalf("expected jab1 in JSON, got %s", string(b))
	}
}
