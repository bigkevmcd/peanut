package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestApp(t *testing.T) {
	goDemo := &App{
		Name: "go-demo",
		Environments: []*Environment{
			{Name: "dev"},
			{Name: "staging"},
			{Name: "production"},
		},
	}

	cfg := &Config{
		Apps: []*App{
			goDemo,
		},
	}

	got := cfg.App("go-demo")

	if diff := cmp.Diff(goDemo, got); diff != "" {
		t.Fatalf("app didn't match:\n%s", diff)
	}

	unknown := cfg.App("unknown")
	if unknown != nil {
		t.Fatalf("got %#v, want nil", unknown)
	}
}

func TestEnvironment(t *testing.T) {
	dev := &Environment{Name: "dev"}
	goDemo := &App{
		Name: "go-demo",
		Environments: []*Environment{
			dev,
			{Name: "staging"},
			{Name: "production"},
		},
	}
	got := goDemo.Environment("dev")

	if diff := cmp.Diff(dev, got); diff != "" {
		t.Fatalf("env didn't match:\n%s", diff)
	}
	if got.App != goDemo {
		t.Fatalf("got %#v, want %#v", got.App, goDemo)
	}

	unknown := goDemo.Environment("unknown")
	if unknown != nil {
		t.Fatalf("got %#v, want nil", unknown)
	}
}

func TestEnvironmentPath(t *testing.T) {
	dev := &Environment{Name: "dev", RelPath: "../dev"}
	goDemo := &App{
		Name: "go-demo",
		Path: "/deploy/environments/base",
		Environments: []*Environment{
			dev,
		},
	}
	dev.App = goDemo

	if v := dev.Path(); v != "/deploy/environments/dev" {
		t.Fatalf("Path() got %#v, want %#v", v, "deploy/environments/dev")
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

	all, err := goDemo.ParseManifests()
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
