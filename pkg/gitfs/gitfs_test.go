package gitfs

import (
	"io/ioutil"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

var _ filesys.FileSystem = gitFS{}

func TestUnsupportedFeatures(t *testing.T) {
	gfs := gitFS{}

	_, err := gfs.Create("testing/file")
	assertIsUnsupported(t, err)
	assertIsUnsupported(t, gfs.Mkdir("testing"))
	assertIsUnsupported(t, gfs.MkdirAll("testing/testing"))
	assertIsUnsupported(t, gfs.RemoveAll("testing/testing"))
	_, err = gfs.Open("testing/file")
	assertIsUnsupported(t, err)
	_, err = gfs.Glob("test*")
	assertIsUnsupported(t, err)
	err = gfs.WriteFile("testing", []byte("testing"))
	assertIsUnsupported(t, err)

	_, err = gfs.ReadDir("")
	assertIsUnsupported(t, err)
}

func TestReadFile(t *testing.T) {
	gfs := makeClonedGFS(t)

	remote, err := gfs.ReadFile("README.md")
	assertNoError(t, err)

	local, err := ioutil.ReadFile("../../README.md")
	assertNoError(t, err)

	if diff := cmp.Diff(local, remote); diff != "" {
		t.Fatalf("failed to fetch correct file:\n%s", diff)
	}
}

func TestIsDir(t *testing.T) {
	gfs := makeClonedGFS(t)

	if gfs.IsDir("README.md") {
		t.Fatal("IsDir() returned true for a file")
	}

	if !gfs.IsDir("pkg") {
		t.Fatal("IsDir() returned false for a directory")
	}
}

func TestCleanedAbs(t *testing.T) {
	gfs := makeClonedGFS(t)

	dir, _, err := gfs.CleanedAbs("pkg/gitfs/gitfs_test.go")
	assertNoError(t, err)

	if dir != filesys.ConfirmedDir("pkg/gitfs") {
		t.Fatalf("got %#v", dir)
	}
}

func TestErrNotSupported(t *testing.T) {
	err := errNotSupported("Created")
	if s := err.Error(); s != "feature \"Created\" not supported" {
		t.Fatalf("got %s, want %v", s, "testing")
	}
}

func assertIsUnsupported(t *testing.T, err error) {
	t.Helper()
	if _, ok := err.(notSupported); !ok {
		t.Fatalf("got %#v, want ErrNotSupported", err)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewInMemoryFromOptions(t *testing.T) {
	gfs := makeClonedGFS(t)

	got, err := gfs.ReadFile("LICENSE")
	if err != nil {
		t.Fatal(err)
	}
	want, err := ioutil.ReadFile("../../LICENSE")
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("failed to read file:\n%s", diff)
	}
}

func makeClonedGFS(t *testing.T) filesys.FileSystem {
	t.Helper()
	gfs, err := NewInMemoryFromOptions(
		&git.CloneOptions{
			URL: "../../",
		})
	assertNoError(t, err)
	return gfs
}
