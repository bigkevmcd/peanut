package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bigkevmcd/peanut/pkg/gitfs"
	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
	"github.com/go-git/go-git/v5"
	"sigs.k8s.io/yaml"
)

// ParseManifests parses the configuration's manifests into overall picture of
// the repository's applications.
// TODO: this should probably accept a fs.FileSystem to allow reusing the Git
// clone.
// TODO: This should also not be a map[string]map[string]map[string][]string :-)
func ParseManifests(a *App) (map[string]map[string]map[string][]string, error) {
	result := map[string]map[string]map[string][]string{}
	gfs, err := gitfs.NewInMemoryFromOptions(&git.CloneOptions{
		URL: a.RepoURL,
	})
	if err != nil {
		return nil, err
	}
	// TODO: This should probably reject data if the app is not the same as
	// a.Name.
	err = a.EachEnvironment(func(e *Environment) error {
		parsed, err := parser.ParseConfig(e.Path(), gfs)
		if err != nil {
			return err
		}
		for _, app := range parsed.Apps {
			envs, ok := result[app.Name]
			if !ok {
				envs = map[string]map[string][]string{}
			}
			appSvcs := map[string][]string{}
			for _, svc := range app.Services {
				appSvcs[svc.Name] = svc.Images[:]
			}
			envs[e.Name] = appSvcs
			result[app.Name] = envs
		}
		return nil
	})
	return result, err
}

// Parse decodes YAML describing an environment manifest.
func Parse(in io.Reader) (*Config, error) {
	m := &Config{}
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(buf, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ParseFile is a wrapper around Parse that accepts a filename, it opens and
// parses the file, and closes it.
func ParseFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open: %s", filename)
	}
	defer f.Close()
	return Parse(f)
}
