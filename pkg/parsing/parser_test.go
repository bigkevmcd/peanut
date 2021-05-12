package parser

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kustomize/api/filesys"
)

const (
	nameLabel   = "app.kubernetes.io/name"
	partOfLabel = "app.kubernetes.io/part-of"
)

func TestParseNoFile(t *testing.T) {
	res, err := ParseFromGit(
		"testdata",
		&git.CloneOptions{
			URL:   "../..",
			Depth: 1,
		})

	if res != nil {
		t.Errorf("did not expect to parse resources: %#v", res)
	}
	if err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestParse(t *testing.T) {
	res, err := ParseConfig(
		"../../testdata/examples/environments/dev",
		filesys.MakeFsOnDisk())

	if err != nil {
		t.Fatal(err)
	}
	// This is only comparing the keys of the objects, as it's assumed that it's
	// Kustomize is parsing the YAML correctly.
	want := []string{
		"dev--v1-Namespace",
		"go-demo-config-dev-v1-ConfigMap",
		"go-demo-http-dev-apps/v1-Deployment",
		"go-demo-http-dev-v1-Service",
		"redis-dev-apps/v1-Deployment",
		"redis-dev-v1-Service",
	}
	got := resKeys(t, res...)
	sort.Strings(got)
	assertCmp(t, want, got, "failed to match parsed resources")
}

func assertCmp(t *testing.T, want, got interface{}, msg string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(msg+":\n%s", diff)
	}
}

func makeCloneOptions() *git.CloneOptions {
	o := &git.CloneOptions{
		URL:   "../..",
		Depth: 1,
	}
	if b := os.Getenv("GITHUB_BASE_REF"); b != "" {
		o.ReferenceName = plumbing.NewBranchReferenceName(b)
		o.URL = "https://github.com/bigkevmcd/go-demo.git"
	}
	return o
}

func resKeys(t *testing.T, objs ...runtime.Object) []string {
	res := make([]string, len(objs))
	for i, o := range objs {
		res[i] = resKey(t, o)
	}
	return res
}

func resKey(t *testing.T, o runtime.Object) string {
	t.Helper()
	oa, err := meta.Accessor(o)
	if err != nil {
		t.Fatalf("failed to get an accessor for %#v: %s", o, err)
	}
	tm, err := meta.TypeAccessor(o)
	if err != nil {
		t.Fatalf("failed to get a type accessor for %#v: %s", o, err)
	}
	return strings.Join([]string{oa.GetName(), oa.GetNamespace(), tm.GetAPIVersion(), tm.GetKind()}, "-")
}
