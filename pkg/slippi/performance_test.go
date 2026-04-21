package slippi

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func TestPerformanceGate(t *testing.T) {
	root := mustProjectRoot(t)
	fixturesRoot := filepath.Join(root, "testdata", "slp")
	fixtures := mustFixtures(t, fixturesRoot)
	limit := goOnlyPerformanceLimit(t)

	// Warmup
	_, _ = runGoIteration(fixtures)

	const measuredRuns = 5
	goRuns := make([]float64, 0, measuredRuns)

	for i := 0; i < measuredRuns; i++ {
		goElapsed, err := runGoIteration(fixtures)
		if err != nil {
			t.Fatalf("go benchmark run %d failed: %v", i+1, err)
		}
		goRuns = append(goRuns, goElapsed)
	}

	goMedian := median(goRuns)
	if goMedian > limit {
		t.Fatalf("performance gate failed: go median %.2fms > limit %.2fms", goMedian, limit)
	}

	t.Logf("performance gate passed: go median %.2fms <= limit %.2fms", goMedian, limit)
}

func BenchmarkGoSummary(b *testing.B) {
	root := mustProjectRoot(b)
	fixtures := mustFixtures(b, filepath.Join(root, "testdata", "slp"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := runGoIteration(fixtures); err != nil {
			b.Fatalf("go iteration failed: %v", err)
		}
	}
}

func runGoIteration(fixtures []string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	start := time.Now()
	for _, fixture := range fixtures {
		game := NewGame(fixture)
		if _, err := game.Summary(ctx); err != nil {
			return 0, fmt.Errorf("summary %s: %w", fixture, err)
		}
	}
	return float64(time.Since(start).Microseconds()) / 1000.0, nil
}

func median(values []float64) float64 {
	cp := slices.Clone(values)
	slices.Sort(cp)
	return cp[len(cp)/2]
}

func goOnlyPerformanceLimit(t *testing.T) float64 {
	t.Helper()
	// Temporary static limit while removing JS runtime dependency.
	// This is intentionally conservative and can be tightened as parser/stats mature.
	return 120000.0
}
