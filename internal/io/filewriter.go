package io

import (
	"fmt"
	stdio "io"
	"os"
	"path/filepath"
)

// CopyReplay writes a native replay copy to targetPath.
func CopyReplay(sourcePath, targetPath string) error {
	in, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source replay: %w", err)
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	out, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("create destination replay: %w", err)
	}
	defer out.Close()

	if _, err := stdio.Copy(out, in); err != nil {
		return fmt.Errorf("copy replay: %w", err)
	}
	return nil
}
