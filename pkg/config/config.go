package config

import (
	"path"
)

// Environment is a k8s namespace/cluster that an application is deployed.
type Environment struct {
	Name    string `json:"name"`
	RelPath string `json:"rel_path"` // This is relative to the Path for the parent App.
	App     *App   `json:"-"`
}

// App represents a high-level application that is deployed across multiple
// environments, and configured through Kustomize.
type App struct {
	Name         string         `json:"name"`
	RepoURL      string         `json:"repo_url"`
	Path         string         `json:"path"`
	Environments []*Environment `json:"environments"`
}

// Config represents the managed apps.
type Config struct {
	Apps []*App `json:"apps,omitempty"`
}

// App returns the named app, or nil if not found.
func (c *Config) App(name string) *App {
	for _, v := range c.Apps {
		if v.Name == name {
			return v
		}
	}
	return nil
}

// Environment gets a named environment.
func (a *App) Environment(name string) *Environment {
	for _, v := range a.Environments {
		if v.Name == name {
			v.App = a
			return v
		}
	}
	return nil
}

// EachEnvironment iterates over each environment within the app, and calls it
// with an environment, the environment will have it's parent app linked
// correctly.
func (a *App) EachEnvironment(f func(e *Environment) error) error {
	for _, v := range a.Environments {
		err := f(a.Environment(v.Name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Path returns the app-relative path for the kustomize.yaml for this
// environment.
//
// For example, app in /test/base and environment in "../dev" would get
// "/test/dev".
func (e *Environment) Path() string {
	return path.Clean(path.Join(e.App.Path, e.RelPath))
}
