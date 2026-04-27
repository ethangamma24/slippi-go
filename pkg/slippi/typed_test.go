package slippi

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func jsonEqual(a, b any) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	var aMap, bMap map[string]any
	json.Unmarshal(aJSON, &aMap)
	json.Unmarshal(bJSON, &bMap)
	return reflect.DeepEqual(aMap, bMap)
}

func fixturePath(name string) string {
	root, _ := projectRoot()
	return filepath.Join(root, "testdata", "slp", name)
}

func TestNewGameFromBytes_Parsed_MatchesAnyFacade(t *testing.T) {
	path := fixturePath("test.slp")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	ctx := context.Background()
	gFile := NewGame(path)
	gBytes := NewGameFromBytes("test.slp", data)

	// Compare Settings
	anySettings, err := gFile.GetSettings(ctx)
	if err != nil {
		t.Fatalf("GetSettings: %v", err)
	}
	typedSettings, err := gBytes.SettingsTyped(ctx)
	if err != nil {
		t.Fatalf("SettingsTyped: %v", err)
	}
	if !jsonEqual(anySettings, typedSettings) {
		anySettingsJSON, _ := json.Marshal(anySettings)
		typedSettingsJSON, _ := json.Marshal(typedSettings)
		t.Fatalf("settings JSON mismatch\nany:  %s\ntyped:%s", anySettingsJSON, typedSettingsJSON)
	}

	// Compare Metadata
	anyMeta, err := gFile.GetMetadata(ctx)
	if err != nil {
		t.Fatalf("GetMetadata: %v", err)
	}
	typedMeta, err := gBytes.MetadataTyped(ctx)
	if err != nil {
		t.Fatalf("MetadataTyped: %v", err)
	}
	if !jsonEqual(anyMeta, typedMeta) {
		anyMetaJSON, _ := json.Marshal(anyMeta)
		typedMetaJSON, _ := json.Marshal(typedMeta)
		t.Fatalf("metadata JSON mismatch\nany:  %s\ntyped:%s", anyMetaJSON, typedMetaJSON)
	}

	// Compare GameEnd
	anyEnd, err := gFile.GetGameEnd(ctx)
	if err != nil {
		t.Fatalf("GetGameEnd: %v", err)
	}
	typedEnd, err := gBytes.GameEndTyped(ctx)
	if err != nil {
		t.Fatalf("GameEndTyped: %v", err)
	}
	if !jsonEqual(anyEnd, typedEnd) {
		anyEndJSON, _ := json.Marshal(anyEnd)
		typedEndJSON, _ := json.Marshal(typedEnd)
		t.Fatalf("gameEnd JSON mismatch\nany:  %s\ntyped:%s", anyEndJSON, typedEndJSON)
	}
}

func TestStatsTyped_MatchesAnyFacade(t *testing.T) {
	path := fixturePath("lCancel.slp")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	ctx := context.Background()
	gFile := NewGame(path)
	gBytes := NewGameFromBytes("lCancel.slp", data)

	anyStats, err := gFile.GetStats(ctx)
	if err != nil {
		t.Fatalf("GetStats: %v", err)
	}
	typedStats, err := gBytes.StatsTyped(ctx)
	if err != nil {
		t.Fatalf("StatsTyped: %v", err)
	}

	if !jsonEqual(anyStats, typedStats) {
		anyStatsJSON, _ := json.Marshal(anyStats)
		typedStatsJSON, _ := json.Marshal(typedStats)
		t.Fatalf("stats JSON mismatch\nany:  %s\ntyped:%s", anyStatsJSON, typedStatsJSON)
	}
}

func TestNewGameFromReader_HandlesIOErrors(t *testing.T) {
	r := &errorReader{err: io.ErrUnexpectedEOF}
	g, err := NewGameFromReader("bad", r)
	if err == nil {
		t.Fatal("expected error from NewGameFromReader")
	}
	if g != nil {
		t.Fatal("expected nil Game on error")
	}
}

type errorReader struct {
	err error
}

func (e *errorReader) Read(_ []byte) (int, error) {
	return 0, e.err
}

func TestGame_NoSource_ErrorsCleanly(t *testing.T) {
	g := &Game{}
	ctx := context.Background()
	_, err := g.Parsed(ctx)
	if err == nil {
		t.Fatal("expected error for Game with no source")
	}
	if err.Error() != "parse game : slippi: Game has no source (use NewGame, NewGameFromBytes, or NewGameFromReader)" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestParseMetaFromBytes_OnlyTouchesMetadata(t *testing.T) {
	path := fixturePath("test.slp")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	ctx := context.Background()
	g := NewGameFromBytes("test.slp", data)
	fullMeta, err := g.MetadataTyped(ctx)
	if err != nil {
		t.Fatalf("MetadataTyped: %v", err)
	}

	fastMeta, err := ParseMetaFromBytes("test.slp", data)
	if err != nil {
		t.Fatalf("ParseMetaFromBytes: %v", err)
	}

	if fullMeta.StartAt != fastMeta.StartAt {
		t.Fatalf("StartAt mismatch: full=%q fast=%q", fullMeta.StartAt, fastMeta.StartAt)
	}
	if fullMeta.LastFrame != fastMeta.LastFrame {
		t.Fatalf("LastFrame mismatch: full=%d fast=%d", fullMeta.LastFrame, fastMeta.LastFrame)
	}
}

func TestNewGameFromBytes_ContextCancellation(t *testing.T) {
	path := fixturePath("test.slp")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	g := NewGameFromBytes("test.slp", data)
	_, err = g.Parsed(ctx)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
