package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigkevmcd/peanut/pkg/http/config"
	"github.com/google/go-cmp/cmp"
)

func TestListApps(t *testing.T) {
	ts := httptest.NewTLSServer(NewRouter(makeConfig()))
	t.Cleanup(ts.Close)

	res, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	assertJSONResponse(t, res, map[string]interface{}{
		"apps": []interface{}{
			map[string]interface{}{
				"name": "go-demo",
			},
		},
	})
}

func TestGetApp(t *testing.T) {
	ts := httptest.NewTLSServer(NewRouter(makeConfig()))
	t.Cleanup(ts.Close)

	res, err := ts.Client().Get(ts.URL + "/apps/go-demo")
	if err != nil {
		t.Fatal(err)
	}
	assertJSONResponse(t, res, map[string]interface{}{
		"name":     "go-demo",
		"repo_url": "https://github.com/bigkevmcd/go-demo.git",
		"path":     "/examples/kustomize/base",
		"environments": []interface{}{
			map[string]interface{}{"name": "dev", "rel_path": "../overlays/dev"},
			map[string]interface{}{"name": "staging", "rel_path": "../overlays/staging"},
			map[string]interface{}{"name": "production", "rel_path": "../overlays/production"},
		},
	})
}

func TestGetEnvironment(t *testing.T) {
	ts := httptest.NewTLSServer(NewRouter(makeConfig()))
	t.Cleanup(ts.Close)

	res, err := ts.Client().Get(ts.URL + "/apps/go-demo/envs/dev")
	if err != nil {
		t.Fatal(err)
	}
	assertJSONResponse(t, res, map[string]interface{}{
		"environment": map[string]interface{}{
			"name":     "dev",
			"rel_path": "../overlays/dev",
		},
	})
}

func makeConfig() *config.Config {
	return &config.Config{
		Apps: []*config.App{
			{
				Name:    "go-demo",
				RepoURL: "https://github.com/bigkevmcd/go-demo.git",
				Path:    "/examples/kustomize/base",
				Environments: []*config.Environment{
					{Name: "dev", RelPath: "../overlays/dev"},
					{Name: "staging", RelPath: "../overlays/staging"},
					{Name: "production", RelPath: "../overlays/production"},
				},
			},
		},
	}
}

// TODO: assert the content-type.
func assertJSONResponse(t *testing.T, res *http.Response, want map[string]interface{}) {
	t.Helper()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]interface{}{}
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("failed to parse %s: %s", b, err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("JSON response failed:\n%s", diff)
	}
}
