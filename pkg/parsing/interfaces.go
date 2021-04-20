package parser

import (
	"github.com/go-git/go-git/v5"
)

// ResourceParser implementations should fetch the source using the CloneOptions and
// parse the resources in the path into a set of resource.Resource values.
type ResourceParser func(path string, opts *git.CloneOptions) ([]*Resource, error)
