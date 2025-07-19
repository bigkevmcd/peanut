package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bigkevmcd/peanut/pkg/config"
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
		"path":     "pkg/config/testdata/go-demo/base",
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

func TestGetDesiredState(t *testing.T) {
	cfg := makeConfig()
	cfg.Apps[0].RepoURL = "../../"

	// TODO: This should be mocked out, by decoupling the behaviour from the
	// App model.
	ts := httptest.NewTLSServer(NewRouter(cfg))
	t.Cleanup(ts.Close)

	res, err := ts.Client().Get(ts.URL + "/apps/go-demo/desired")
	if err != nil {
		t.Fatal(err)
	}
	assertJSONResponse(t, res, map[string]interface{}{
		"name":     "go-demo",
		"path":     "pkg/config/testdata/go-demo/base",
		"repo_url": "../../",
		"environments": []interface{}{
			map[string]interface{}{
				"name":     "dev",
				"rel_path": "../overlays/dev",
				"services": []interface{}{
					map[string]interface{}{"images": []interface{}{"bigkevmcd/go-demo:latest"}, "name": "go-demo-http"},
					map[string]interface{}{"images": []interface{}{"redis:6-alpine"}, "name": "redis"},
				},
			},
			map[string]interface{}{
				"name":     "staging",
				"rel_path": "../overlays/staging",
				"services": []interface{}{
					map[string]interface{}{"images": []interface{}{"bigkevmcd/go-demo:staging"}, "name": "go-demo-http"},
					map[string]interface{}{"images": []interface{}{"redis:6-alpine"}, "name": "redis"},
				},
			},
			map[string]interface{}{
				"name":     "production",
				"rel_path": "../overlays/production",
				"services": []interface{}{
					map[string]interface{}{"images": []interface{}{"bigkevmcd/go-demo:production"}, "name": "go-demo-http"},
					map[string]interface{}{"images": []interface{}{"redis:6-alpine"}, "name": "redis"},
				},
			},
		},
	})
}

func makeConfig() *config.Config {
	return &config.Config{
		Apps: []*config.App{
			{
				Name:    "go-demo",
				RepoURL: "https://github.com/bigkevmcd/go-demo.git",
				Path:    "pkg/config/testdata/go-demo/base",
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
	if res.StatusCode != http.StatusOK {
		t.Fatalf("didn't get a successful response: %v", res.StatusCode)
	}
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
		t.Errorf("JSON response failed:\n%s", diff)
	}
}
