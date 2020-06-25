package kustomize

import (
	"testing"

	"sigs.k8s.io/kustomize/pkg/image"
	"sigs.k8s.io/kustomize/pkg/types"

	"github.com/google/go-cmp/cmp"
)

func TestAddImageOverride(t *testing.T) {
	k := createKustomizer()

	err := k.AddImageOverride("test/built-image", "ef7bebf8bdb1919d947afe46ab4b2fb4278039b3")
	fatalIfError(t, err)

	want := &types.Kustomization{
		Images: []image.Image{
			{Name: "test/built-image", NewTag: "ef7bebf8bdb1919d947afe46ab4b2fb4278039b3"},
		},
	}
	if diff := cmp.Diff(want, k.Kustomization()); diff != "" {
		t.Fatalf("Kustomization didn't match:\n%s", diff)
	}
}

func TestOverrideImageUpdatesExistingOverride(t *testing.T) {
	k := createKustomizer()
	fatalIfError(t, k.AddImageOverride("test/built-image", "ef7bebf8bdb1919d947afe46ab4b2fb4278039b3"))

	fatalIfError(t, k.AddImageOverride("test/built-image", "6d56c1750f65b4f648040959313c5004e5c351cb"))

	want := &types.Kustomization{
		Images: []image.Image{
			{Name: "test/built-image", NewTag: "6d56c1750f65b4f648040959313c5004e5c351cb"},
		},
	}
	if diff := cmp.Diff(want, k.Kustomization()); diff != "" {
		t.Fatalf("Kustomization didn't match:\n%s", diff)
	}
}

func createKustomizer() *Kustomizer {
	k := &types.Kustomization{}
	return NewKustomizer(k)
}

func fatalIfError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
