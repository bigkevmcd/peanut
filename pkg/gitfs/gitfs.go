package gitfs

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/memory"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// gitFS is an internal implementation of the Kustomize
// filesystem abstraction.
type gitFS struct {
	tree *object.Tree
}

// New creates and returns a go-git storage adapter.
func New(t *object.Tree) filesys.FileSystem {
	return &gitFS{tree: t}
}

// NewInMemoryFromOptions clones a Git repository into memory.
func NewInMemoryFromOptions(opts *git.CloneOptions) (filesys.FileSystem, error) {
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
	return New(tree), nil
}

// ReadDir implements filesys.FileSystem.
func (g gitFS) ReadDir(path string) ([]string, error) {
	return nil, errNotSupported("Walk")
}

// IsDir implements filesys.FileSystem.
func (g gitFS) IsDir(name string) bool {
	// If it exists as a file, it's not a directory.
	_, err := g.tree.File(name)
	if err == nil {
		return false
	}
	// Git doesn't store directories.
	//
	// If we can find a file with a prefix of the name we're looking for, then
	// the name is a directory.
	//
	// TODO: make this a bit more efficent, cache found dirs?
	isDir := false
	err = g.tree.Files().ForEach(func(f *object.File) error {
		if strings.HasPrefix(f.Name, name) {
			isDir = true
			return storer.ErrStop
		}
		return nil
	})
	// TODO: not a lot of choice here, there's no scope for returning an error.
	if err != nil {
		panic(err)
	}
	return isDir
}

// CleanedAbs implements filesys.FileSystem.
func (g gitFS) CleanedAbs(p string) (filesys.ConfirmedDir, string, error) {
	if g.IsDir(p) {
		return filesys.ConfirmedDir(p), "", nil
	}
	d := path.Dir(p)
	f := path.Base(p)
	return filesys.ConfirmedDir(d), f, nil
}

// ReadFile implements filesys.FileSystem.
func (g gitFS) ReadFile(name string) ([]byte, error) {
	f, err := g.tree.File(name)
	if err != nil {
		return nil, err
	}
	b, err := f.Contents()
	if err != nil {
		return nil, err
	}
	return []byte(b), nil
}

// Walk implements filesys.FileSystem.
func (g gitFS) Walk(path string, walkFn filepath.WalkFunc) error {
	return errNotSupported("Walk")
}

// Create implements filesys.FileSystem.
func (g gitFS) Create(name string) (filesys.File, error) {
	return nil, errNotSupported("Create")
}

// MkDir implements filesys.FileSystem.
func (g gitFS) Mkdir(name string) error {
	return errNotSupported("MkDir")
}

// MkDirAll implements filesys.FileSystem.
func (g gitFS) MkdirAll(name string) error {
	return errNotSupported("MkdirAll")
}

// RemoveAll implements filesys.FileSystem.
func (g gitFS) RemoveAll(name string) error {
	return errNotSupported("RemoveAll")
}

// Open implements filesys.FileSystem.
func (g gitFS) Open(name string) (filesys.File, error) {
	return nil, errNotSupported("Open")
}

// Exists implements filesys.FileSystem.
func (g gitFS) Exists(name string) bool {
	return false
}

// Glob implements filesys.FileSystem.
func (g gitFS) Glob(pattern string) ([]string, error) {
	return nil, errNotSupported("Glob")
}

// WriteFile implements filesys.FileSystem.
func (g gitFS) WriteFile(name string, data []byte) error {
	return errNotSupported("WriteFile")
}

func errNotSupported(s string) error {
	return notSupported(s)
}

type notSupported string

func (f notSupported) Error() string {
	return fmt.Sprintf("feature %#v not supported", string(f))
}
