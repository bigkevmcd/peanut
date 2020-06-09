package parser

import (
	"sort"

	"sigs.k8s.io/kustomize/k8sdeps"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/target"
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
	cfg := &Config{AppsToServices: map[string][]string{}, Services: map[string]*Service{}}
	k8sfactory := k8sdeps.NewFactory()
	fs := fs.MakeRealFS()
	ldr, err := loader.NewLoader(path, fs)
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

func extractService(v map[string]interface{}) *Service {
	meta := v["metadata"].(map[string]interface{})
	spec := v["spec"].(map[string]interface{})
	templateSpec := spec["template"].(map[string]interface{})["spec"].(map[string]interface{})
	svc := &Service{
		Name:      meta["name"].(string),
		Namespace: mapString("namespace", meta),
		Replicas:  spec["replicas"].(int64),
		Images:    []string{},
	}
	for _, v := range templateSpec["containers"].([]interface{}) {
		svc.Images = append(svc.Images, v.(map[string]interface{})["image"].(string))
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
