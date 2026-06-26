package scaffold_test

import (
	"os"
	"path/filepath"
	"testing"
)

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func TestLayerDirectoriesExist(t *testing.T) {
	root := moduleRoot(t)
	required := []string{
		"internal/domain",
		"internal/application",
		"internal/infrastructure",
		"internal/interfaces",
	}
	for _, rel := range required {
		path := filepath.Join(root, rel)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
		if !info.IsDir() {
			t.Fatalf("%s is not a directory", rel)
		}
	}
}

func TestLayerPackagesHaveSourceFiles(t *testing.T) {
	root := moduleRoot(t)
	layers := []string{
		"internal/domain",
		"internal/application",
		"internal/infrastructure",
		"internal/interfaces",
	}
	for _, rel := range layers {
		entries, err := os.ReadDir(filepath.Join(root, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		var goFiles int
		for _, e := range entries {
			if !e.IsDir() && filepath.Ext(e.Name()) == ".go" {
				goFiles++
			}
		}
		if goFiles == 0 {
			t.Fatalf("%s has no .go source files", rel)
		}
	}
}
