package parser

import (
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"

	"sigs.k8s.io/kustomize/k8sdeps"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/target"

	"github.com/bigkevmcd/peanut/pkg/gitfs"
)

const (
	serviceLabel = "app.kubernetes.io/name"
	appLabel     = "app.kubernetes.io/part-of"
)

// Config is a representation of the apps and services, and configurations for
// the services.
type Config struct {
	AppsToServices map[string][]string
	Services       map[string]*Service
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
	fs := fs.MakeRealFS()
	return parseConfig(path, fs)
}

// ParseFromGit takes a go-git CloneOptions struct and a filepath, and extracts
// the service configuration from there.
func ParseFromGit(path string, opts *git.CloneOptions) (*Config, error) {
	clone, err := git.Clone(memory.NewStorage(), nil, opts)
	if err != nil {
		return nil, err
	}
	ref, err := clone.Head()
	if err != nil {
		return nil, err
	}
	commit, err := clone.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	gfs := gitfs.New(tree)
	return parseConfig(path, gfs)
}

func parseConfig(path string, files fs.FileSystem) (*Config, error) {
	cfg := &Config{AppsToServices: map[string][]string{}, Services: map[string]*Service{}}
	k8sfactory := k8sdeps.NewFactory()
	ldr, err := loader.NewLoader(path, files)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = ldr.Cleanup()
		if err != nil {
			panic(err)
		}
	}()
	kt, err := target.NewKustTarget(ldr, k8sfactory.ResmapF, k8sfactory.TransformerF)
	if err != nil {
		return nil, err
	}
	r, err := kt.MakeCustomizedResMap()
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		return nil, nil
	}
	for k, v := range r {
		gvk := k.Gvk()
		switch gvk.Kind {
		case "Deployment":
			svc := extractAppAndServices(v.GetLabels(), cfg.AppsToServices)
			cfg.Services[svc] = extractService(v.Map())
		}
	}
	return cfg, nil
}

// ParseFromGit takes a go-git CloneOptions struct and a filepath, and extracts
// the service configuration from there.
func ParseFromGit(path string, opts *git.CloneOptions) (*Config, error) {
	cfg := &Config{AppsToServices: map[string][]string{}, Services: map[string]*Service{}}
	k8sfactory := k8sdeps.NewFactory()
	clone, err := git.Clone(memory.NewStorage(), nil, opts)
	if err != nil {
		return nil, err
	}
	ref, err := clone.Head()
	if err != nil {
		return nil, err
	}
	commit, err := clone.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	gfs := gitfs.New(tree)
	ldr, err := loader.NewLoader(path, gfs)
	if err != nil {
		return nil, err
	}
	defer ldr.Cleanup()
	kt, err := target.NewKustTarget(ldr, k8sfactory.ResmapF, k8sfactory.TransformerF)
	if err != nil {
		return nil, err
	}
	r, err := kt.MakeCustomizedResMap()
	if err != nil {
		return nil, err
	}
	if len(r) == 0 {
		return nil, nil
	}
	for k, v := range r {
		gvk := k.Gvk()
		switch gvk.Kind {
		case "Deployment":
			svc := extractAppAndServices(v.GetLabels(), cfg.AppsToServices)
			cfg.Services[svc] = extractService(v.Map())
		}
	}
	return cfg, nil
}

func extractAppAndServices(meta map[string]string, state map[string][]string) string {
	app, svc := appAndService(meta)
	appSvcs, ok := state[app]
	if !ok {
		appSvcs = []string{}
	}
	appSvcs = append(appSvcs, svc)
	sort.Strings(appSvcs)
	state[app] = appSvcs
	return svc
}

func appAndService(v map[string]string) (string, string) {
	return v[appLabel], v[serviceLabel]
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
		Replicas:  spec["replicas"].(int64),
		Images:    []string{},
	}
	for _, v := range templateSpec["containers"].([]interface{}) {
		svc.Images = append(svc.Images, mapString("image", v.(map[string]interface{})))
	}
	return svc
}

func mapString(k string, v map[string]interface{}) string {
	s, ok := v[k].(string)
	if !ok {
		return ""
	}
	return s
}
