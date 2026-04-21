package slippi

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func mustProjectRoot(tb testing.TB) string {
	tb.Helper()
	dir, err := os.Getwd()
	if err != nil {
		tb.Fatalf("getwd failed: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			tb.Fatalf("could not find go.mod from %s", dir)
		}
		dir = parent
	}
}

func mustFixtures(tb testing.TB, root string) []string {
	tb.Helper()
	fixtures := make([]string, 0, 64)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".slp") {
			if strings.EqualFold(d.Name(), "incomplete.slp") {
				return nil
			}
			fixtures = append(fixtures, path)
		}
		return nil
	})
	if err != nil {
		tb.Fatalf("walk fixtures failed: %v", err)
	}
	sort.Strings(fixtures)
	return fixtures
}
