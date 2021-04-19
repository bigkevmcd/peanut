package pipeline

import (
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/google/go-cmp/cmp"
)

func TestListStages_with_no_dir(t *testing.T) {
	fs := memfs.New()

	got, err := ListStages(fs, "pipeline")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{}, got); diff != "" {
		t.Fatalf("failed:\n%s", diff)
	}
}

func TestListStages_with_dir(t *testing.T) {
	fs := memfs.New()
	if err := fs.MkdirAll("pipeline", 0644); err != nil {
		t.Fatal(err)
	}

	got, err := ListStages(fs, "pipeline")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{}, got); diff != "" {
		t.Fatalf("failed:\n%s", diff)
	}
}

func TestListStages_with_stages(t *testing.T) {
	fs := memfs.New()
	writeFiles(t, fs, "pipeline", "dev", "staging", "production")

	got, err := ListStages(fs, "pipeline")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"dev", "production", "staging"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("failed:\n%s", diff)
	}
}

func TestListStages_with_numbered_stages(t *testing.T) {
	fs := memfs.New()
	writeFiles(t, fs, "pipeline", "01_dev", "03_production", "02_qa")

	got, err := ListStages(fs, "pipeline")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"dev", "qa", "production"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("failed:\n%s", diff)
	}
}

func writeFiles(t *testing.T, fs billy.Filesystem, base string, stages ...string) {
	if err := fs.MkdirAll(base, 0644); err != nil {
		t.Fatal(err)
	}
	for _, s := range stages {
		fullname := filepath.Join(base, s, "config.yaml")
		f, err := fs.Create(fullname)
		if err != nil {
			t.Fatalf("failed to create %q: %s", fullname, err)
		}

		if _, err = f.Write([]byte(s + "\n")); err != nil {
			t.Fatalf("failed to write to %q: %s", fullname, err)
		}

		if err := f.Close(); err != nil {
			t.Fatalf("failed to close %q: %s", fullname, err)
		}
	}
}
