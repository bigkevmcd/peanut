package kustomize

import (
	"sigs.k8s.io/kustomize/v3/pkg/image"
	"sigs.k8s.io/kustomize/v3/pkg/types"
)

type Kustomizer struct {
	imageOverrides map[string]image.Image
	src            *types.Kustomization
}

// NewKustomizer creates and returns a new Kustomizer for manipulating
// Kustomization files.
func NewKustomizer(k *types.Kustomization) *Kustomizer {
	imageOverrides := map[string]image.Image{}
	for _, v := range k.Images {
		imageOverrides[v.Name] = v
	}
	return &Kustomizer{
		src:            k,
		imageOverrides: imageOverrides,
	}
}

// AddImageOverride adds an override for a specific image.
//
// Existing overrides for the same image are replaced.
func (k *Kustomizer) AddImageOverride(srcImage, newTag string) error {
	k.imageOverrides[srcImage] = image.Image{Name: srcImage, NewTag: newTag}
	return nil
}

// Kustomization gets the updated configuration.
func (k *Kustomizer) Kustomization() *types.Kustomization {
	images := []image.Image{}
	for _, v := range k.imageOverrides {
		images = append(images, v)
	}
	k.src.Images = images
	return k.src
}
