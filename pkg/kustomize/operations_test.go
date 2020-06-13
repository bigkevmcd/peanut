package kustomize

import (
	"testing"

	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/types"

	"github.com/bigkevmcd/peanut/pkg/parser"
)

func TestOverrideImageUpdatesExistingOverride(t *testing.T) {
}

func TestOverrideImageRemovesOverrideIfMatchese(t *testing.T) {
}

func TestOverrideImage(t *testing.T) {
	testFs := fs.MakeFakeFS()
	testFs.WriteTestKustomization()

	_ = types.Kustomization{}

	err := OverrideImage(testFs)

	if err != nil {
		t.Fatal(err)
	}
}

func testApp() *parser.App {
	return &parser.App{
		Name: "go-demo",
		Services: []*parser.Service{
			{Name: "go-demo-http", Replicas: 1, Images: []string{"bigkevmcd/go-demo:876ecb3"}},
			{Name: "redis", Replicas: 1, Images: []string{"redis:6-alpine"}},
		},
	}
}
