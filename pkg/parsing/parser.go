package parser

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/resource"

	"github.com/bigkevmcd/peanut/pkg/gitfs"
	"github.com/bigkevmcd/peanut/pkg/kustomize/parser"
)

// ParseFromGit takes a go-git CloneOptions struct and a filepath, and extracts
// the service configuration from there.
func ParseFromGit(path string, opts *git.CloneOptions) ([]runtime.Object, error) {
	gfs, err := gitfs.NewInMemoryFromOptions(opts)
	if err != nil {
		return nil, err
	}
	return ParseConfig(path, gfs)
}

// ParseConfig accepts a path within a repository to extract objects from.
func ParseConfig(dirPath string, fs filesys.FileSystem) ([]runtime.Object, error) {
	r, err := parser.ParseTreeToResMap(dirPath, fs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from %q: %w", dirPath, err)
	}
	conv, err := newUnstructuredConverter()
	if err != nil {
		return nil, err
	}
	var resources []runtime.Object
	for _, v := range r.Resources() {
		converted, err := convert(conv, v)
		if err != nil {
			return nil, err
		}
		resources = append(resources, converted)
	}
	return resources, nil
}

// convert converts a Kustomize resource into a generic Unstructured resource
// which which the unstructured converter uses to create resources from.
func convert(conv *unstructuredConverter, r *resource.Resource) (runtime.Object, error) {
	uns := &unstructured.Unstructured{
		Object: r.Map(),
	}
	t, err := conv.fromUnstructured(uns)
	if err != nil {
		return nil, err
	}
	return t, nil

}
