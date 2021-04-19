package parser

import (
	"sort"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestParseNoFile(t *testing.T) {
	app, err := Parse("testdata")

	if app != nil {
		t.Errorf("did not expect to parse an app: %#v", app)
	}
	if err == nil {
		t.Fatal("expected to get an error")
	}
}

// func TestParseApplication(t *testing.T) {
// 	parseTests := []struct {
// 		filename    string
// 		description string
// 		want        []runtime.Object
// 	}{
// 		{
// 			"testdata/app1",
// 			"empty kustomization",
// 			nil,
// 		},
// 		{
// 			"testdata/go-demo",
// 			"completely local - paths refer to relative paths",
// 			nil,
// 		},
// 		{
// 			"testdata/app2",
// 			"local file refers to a remote path",
// 			nil,
// 		},
// 		{
// 			"testdata/app3",
// 			"Kustomize 3 configuration - remote path",
// 			nil,
// 		},
// 	}

// 	for _, tt := range parseTests {
// 		app, err := Parse(tt.filename)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if diff := cmp.Diff(tt.want, app); diff != "" {
// 			t.Errorf("%s failed to parse:\n%s", tt.filename, diff)
// 		}
// 	}
// }

func TestParseFromGit(t *testing.T) {
	res, err := ParseFromGit(
		"pkg/kustomize/parser/testdata/go-demo",
		&git.CloneOptions{
			URL:   "../../..",
			Depth: 1,
		})

	if err != nil {
		t.Fatal(err)
	}
	sort.SliceStable(res, func(i, j int) bool { return resKey(t, res[i]) < resKey(t, res[j]) })

	want := []runtime.Object{}
	// sort.SliceStable(want, func(i, j int) bool { return resKey(want[i]) < resKey(want[j]) })
	assertCmp(t, want, res, "failed to match parsed resources")
}

func TestParseApplicationFromGit(t *testing.T) {
	app, err := ParseFromGit(
		"pkg/kustomize/parser/testdata/go-demo",
		&git.CloneOptions{
			URL:   "../../..",
			Depth: 1,
		})
	if err != nil {
		t.Fatal(err)
	}

	assertCmp(t, "testing", app, "failed to match app")
}

func assertCmp(t *testing.T, want, got interface{}, msg string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(msg+":\n%s", diff)
	}
}

func resKey(t *testing.T, o runtime.Object) string {
	oa, err := meta.Accessor(o)
	if err != nil {
		t.Fatalf("failed to get the object meta for object %#v: %s", o, err)
	}
	ta, err := meta.TypeAccessor(o)
	if err != nil {
		t.Fatalf("failed to get the type meta for object %#v: %s", o, err)
	}

	return strings.Join([]string{oa.GetName(), oa.GetNamespace(), ta.GetKind(), ta.GetAPIVersion()}, "-")
}
