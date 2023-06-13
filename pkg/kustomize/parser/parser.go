package parser

import (
	"fmt"
	"sort"

	"github.com/go-git/go-git/v5"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kyaml/filesys"

	"github.com/bigkevmcd/peanut/pkg/gitfs"
)

const (
	serviceLabel = "app.kubernetes.io/name"
	appLabel     = "app.kubernetes.io/part-of"
)

// Config is a representation of the apps and services, and configurations for
// the services.
type Config struct {
	Apps []*App
}

// App gets the named app from the config, or returns nil if none exist.
func (c *Config) App(s string) *App {
	for _, v := range c.Apps {
		if v.Name == s {
			return v
		}
	}
	return nil
}

// App is a component with multiple services, and in multiple environments.
type App struct {
	Name     string
	Services []*Service
}

// Service is a representation of a component within the Apps/Services model.
type Service struct {
	Name      string
	Namespace string
	Replicas  int64
	Images    []string
}

// Parse takes a path to a kustomization.yaml file and extracts the service
// configuration from the built resources.
//
// Currently assumes that the standard Kubernetes annotations are used
// (app.kubernetes.io) to identify apps and services (part-of is the app name,
// name is the service name)
//
// Also multi-Deployment services are not supported currently.
func Parse(path string) (*Config, error) {
	fs := filesys.MakeFsOnDisk()
	return ParseConfig(path, fs)
}

// ParseFromGit takes a go-git CloneOptions struct and a filepath, and extracts
// the service configuration from there.
func ParseFromGit(path string, opts *git.CloneOptions) (*Config, error) {
	gfs, err := gitfs.NewInMemoryFromOptions(opts)
	if err != nil {
		return nil, err
	}
	return ParseConfig(path, gfs)
}

// ParseConfig takes a path and an implementation of the kustomize fs.FileSystem
// and parses the configuration into apps.
func ParseConfig(path string, files filesys.FileSystem) (*Config, error) {
	cfg := &Config{Apps: []*App{}}
	resMap, err := ParseTreeToResMap(path, files)
	if err != nil {
		return nil, err
	}
	if resMap.Size() == 0 {
		return nil, nil
	}
	for _, k := range resMap.AllIds() {
		switch k.Gvk.Kind {
		case "Deployment":
			r, err := resMap.GetById(k)
			if err != nil {
				return nil, fmt.Errorf("failed to get resource %v: %w", k, err)
			}
			name := appName(r.GetLabels())
			if name == "" {
				continue
			}
			app := cfg.App(name)
			if app == nil {
				app = &App{Name: name}
				cfg.Apps = append(cfg.Apps, app)
			}
			data, err := r.Map()
			if err != nil {
				return nil, fmt.Errorf("failed to get object data: %w", err)
			}
			svc := extractService(data)
			app.Services = append(app.Services, svc)
			sort.Slice(app.Services, func(i, j int) bool { return app.Services[i].Name < app.Services[j].Name })
		}
	}

	return cfg, nil
}

// ParseTreeToResMap is the main Kustomize parsing mechanism, it returns the raw
// objects parsed by Kustomize.
func ParseTreeToResMap(dirPath string, files filesys.FileSystem) (resmap.ResMap, error) {
	k := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	return k.Run(files, dirPath)
}

func appName(r map[string]string) string {
	return r[appLabel]
}

// TODO: write a generic dotted path walker for the map[string]interface{}
// (again).
func extractService(v map[string]interface{}) *Service {
	meta := v["metadata"].(map[string]interface{})
	spec := v["spec"].(map[string]interface{})
	templateSpec := spec["template"].(map[string]interface{})["spec"].(map[string]interface{})
	svc := &Service{
		Name:      mapString("name", meta),
		Namespace: mapString("namespace", meta),
		Replicas:  mapInt64("replicas", spec),
		Images:    []string{},
	}
	for _, v := range templateSpec["containers"].([]interface{}) {
		svc.Images = append(svc.Images, mapString("image", v.(map[string]interface{})))
	}
	return svc
}

func mapInt64(k string, v map[string]interface{}) int64 {
	switch i := v[k].(type) {
	case int:
		return int64(i)
	case int64:
		return i
	default:
		return 0
	}
}

func mapString(k string, v map[string]interface{}) string {
	s, ok := v[k].(string)
	if !ok {
		return ""
	}
	return s
}
