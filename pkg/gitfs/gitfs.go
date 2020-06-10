package gitfs

import (
	"fmt"
	"path"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"sigs.k8s.io/kustomize/pkg/fs"
)

// gitFS is an internal implementation of the Kustomize
// filesystem abstraction.
=======
type notSupported string

func (f notSupported) Error() string {
	return fmt.Sprintf("feature %#v not supported", string(f))
}

>>>>>>> Implement support for reading from a Git repo.
type gitFS struct {
	tree *object.Tree
}

// New creates and returns a go-git storage adapter.
func New(t *object.Tree) fs.FileSystem {
	return &gitFS{tree: t}
}

<<<<<<< HEAD
// Create implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) Create(name string) (fs.File, error) {
	return nil, errNotSupported("Create")
}

<<<<<<< HEAD
// MkDir implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) Mkdir(name string) error {
	return errNotSupported("MkDir")
}

<<<<<<< HEAD
// MkDirAll implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) MkdirAll(name string) error {
	return errNotSupported("MkdirAll")
}

<<<<<<< HEAD
// RemoveAll implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) RemoveAll(name string) error {
	return errNotSupported("RemoveAll")
}

<<<<<<< HEAD
// Open implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) Open(name string) (fs.File, error) {
	return nil, errNotSupported("Open")
}

<<<<<<< HEAD
// IsDir implements fs.FileSystem.
func (g gitFS) IsDir(name string) bool {
	// If it exists as a file, it's not a directory.
=======
func (g gitFS) IsDir(name string) bool {
	// If it exists as a file, it's not a directory, Git doesn't store files.
>>>>>>> Implement support for reading from a Git repo.
	_, err := g.tree.File(name)
	if err == nil {
		return false
	}
<<<<<<< HEAD
	// Git doesn't store directories.
	//
=======
>>>>>>> Implement support for reading from a Git repo.
	// If we can find a file with a prefix of the name we're looking for, then
	// the name is a directory.
	//
	// TODO: make this a bit more efficent, cache found dirs?
	isDir := false
<<<<<<< HEAD
	err = g.tree.Files().ForEach(func(f *object.File) error {
=======
	g.tree.Files().ForEach(func(f *object.File) error {
>>>>>>> Implement support for reading from a Git repo.
		if strings.HasPrefix(f.Name, name) {
			isDir = true
			return storer.ErrStop
		}
		return nil
	})
<<<<<<< HEAD
	// TODO: not a lot of choice here, there's no scope for returning an error.
	if err != nil {
		panic(err)
	}
	return isDir
}

// CleanedAbs implements fs.FileSystem.
=======
	return isDir
}

>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) CleanedAbs(p string) (fs.ConfirmedDir, string, error) {
	if g.IsDir(p) {
		return fs.ConfirmedDir(p), "", nil
	}
	d := path.Dir(p)
	f := path.Base(p)
	return fs.ConfirmedDir(d), f, nil
}

<<<<<<< HEAD
// Exists implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) Exists(name string) bool {
	return false
}

<<<<<<< HEAD
// Glob implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) Glob(pattern string) ([]string, error) {
	return nil, errNotSupported("Glob")
}

// ReadFile implements fs.FileSystem.
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

<<<<<<< HEAD
// WriteFile implements fs.FileSystem.
=======
>>>>>>> Implement support for reading from a Git repo.
func (g gitFS) WriteFile(name string, data []byte) error {
	return errNotSupported("WriteFile")
}

func errNotSupported(s string) error {
	return notSupported(s)
}
<<<<<<< HEAD

type notSupported string

func (f notSupported) Error() string {
	return fmt.Sprintf("feature %#v not supported", string(f))
}
=======
>>>>>>> Implement support for reading from a Git repo.
