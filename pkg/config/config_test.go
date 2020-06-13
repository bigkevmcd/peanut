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

}
