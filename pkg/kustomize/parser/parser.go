package parser

import (
	"log"

	"github.com/go-git/go-git/v5"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/konfig"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	kustypes "sigs.k8s.io/kustomize/api/types"

	"github.com/bigkevmcd/peanut/pkg/gitfs"
)

const (
	serviceLabel = "app.kubernetes.io/name"
	appLabel     = "app.kubernetes.io/part-of"
)

// Parse takes a path to a kustomization.yaml file and extracts the service
// configuration from the built resources.
//
// Currently assumes that the standard Kubernetes annotations are used
// (app.kubernetes.io) to identify apps and services (part-of is the app name,
// name is the service name)
//
// Also multi-Deployment services are not supported currently.
func Parse(path string) ([]runtime.Object, error) {
	fs := filesys.MakeFsOnDisk()
	return ParseConfig(path, fs)
}

// ParseFromGit takes a go-git CloneOptions struct and a filepath, and extracts
// the service configuration from there.
func ParseFromGit(path string, opts *git.CloneOptions) ([]runtime.Object, error) {
	gfs, err := gitfs.NewInMemoryFromOptions(opts)
	if err != nil {
		return nil, err
	}
	return ParseConfig(path, gfs)
}

// ParseConfig takes a path and an implementation of the kustomize fs.FileSystem
// and parses the configuration into apps.
func ParseConfig(path string, files filesys.FileSystem) ([]runtime.Object, error) {
	resMap, err := ParseTreeToResMap(path, files)
	if err != nil {
		return nil, err
	}
	if resMap.Size() == 0 {
		return nil, nil
	}
	conv, err := newUnstructuredConverter()
	if err != nil {
		return nil, err
	}
	results := make([]runtime.Object, resMap.Size())
	for _, k := range resMap.Resources() {
		extractResource(conv, k)
	}

	return results, nil
}

// ParseTreeToResMap is the main Kustomize parsing mechanism, it returns the raw
// objects parsed by Kustomize.
func ParseTreeToResMap(dirPath string, files filesys.FileSystem) (resmap.ResMap, error) {
	buildOptions := &krusty.Options{
		UseKyaml:               false,
		DoLegacyResourceSort:   true,
		LoadRestrictions:       kustypes.LoadRestrictionsNone,
		AddManagedbyLabel:      false,
		DoPrune:                false,
		PluginConfig:           konfig.DisabledPluginConfig(),
		AllowResourceIdChanges: false,
	}

	k := krusty.MakeKustomizer(files, buildOptions)
	return k.Run(dirPath)
}

// If this is an unknown type (to the converter) no images will be extracted.
func extractResource(conv *unstructuredConverter, res *resource.Resource) runtime.Object {
	obj := &unstructured.Unstructured{
		Object: res.Map(),
	}
	t, err := conv.fromUnstructured(obj)
	if err != nil {
		return nil
	}
	log.Printf("KEVIN!! %v\n", t)
	return nil
}
