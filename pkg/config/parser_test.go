package config

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	parseTests := []struct {
		filename string
		want     *Config
	}{
		{"testdata/example1.yaml",
			&Config{
				Apps: []*App{
					{
						Name:    "go-demo",
						RepoURL: "https://github.com/bigkevmcd/go-demo.git",
						Path:    "/examples/kustomize/base",
						Environments: []*Environment{
							{Name: "dev", RelPath: "../overlays/dev"},
							{Name: "staging", RelPath: "../overlays/staging"},
							{Name: "production", RelPath: "../overlays/production"},
						},
					},
				},
			},
		},
	}

	for _, tt := range parseTests {
		t.Run(fmt.Sprintf("parsing %s", tt.filename), func(rt *testing.T) {
			got, err := ParseFile(tt.filename)
			if err != nil {
				rt.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				rt.Errorf("Parse(%s) failed diff\n%s", tt.filename, diff)
			}
		})
	}
}

func TestAppParseManifests(t *testing.T) {
	goDemo := &App{
		Name:    "go-demo",
		RepoURL: "../..",
		Path:    "pkg/config/testdata/go-demo/base",
		Environments: []*Environment{
			{Name: "dev", RelPath: "../overlays/dev"},
			{Name: "production", RelPath: "../overlays/production"},
			{Name: "staging", RelPath: "../overlays/staging"},
		},
	}

	all, err := ParseManifests(goDemo)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]map[string]map[string][]string{
		"go-demo": {
			"dev":        {"go-demo-http": {"bigkevmcd/go-demo:latest"}, "redis": {"redis:6-alpine"}},
			"production": {"go-demo-http": {"bigkevmcd/go-demo:production"}, "redis": {"redis:6-alpine"}},
			"staging":    {"go-demo-http": {"bigkevmcd/go-demo:staging"}, "redis": {"redis:6-alpine"}},
		},
	}
	assertCmp(t, want, all, "failed to parse manifests")
}

func assertCmp(t *testing.T, want, got interface{}, msg string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(msg+":\n%s", diff)
	}
}
