package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseEmptyFile(t *testing.T) {
	app, err := Parse("testdata/app1")
	if err != nil {
		t.Fatal(err)
	}
	if app != nil {
		t.Fatalf("wanted no app, got %#v", app)
	}
}

func TestParseNoFile(t *testing.T) {
	app, err := Parse("testdata")

	if app != nil {
		t.Errorf("did not expect to parse an app: %#v", app)
	}
	if err == nil {
		t.Fatal("expected to get an error")
	}
}

func TestParseApplication(t *testing.T) {
	app, err := Parse("testdata/go-demo")
	if err != nil {
		t.Fatal(err)
	}

	want := &Config{
		AppsToServices: map[string][]string{
			"go-demo": {"go-demo", "redis"},
		},
		Services: map[string]*Service{
			"go-demo": {Name: "go-demo-http", Replicas: 1, Images: []string{"bigkevmcd/go-demo:876ecb3"}},
			"redis":   {Name: "redis", Replicas: 1, Images: []string{"redis:6-alpine"}},
		},
	}
	assertCmp(t, want, app, "failed to match app")
}

func TestExtractAppAndServices(t *testing.T) {
	redis := map[string]string{"app.kubernetes.io/name": "redis", "app.kubernetes.io/part-of": "go-demo"}
	state := map[string][]string{}

	extractAppAndServices(redis, state)

	want := map[string][]string{
		"go-demo": {"redis"},
	}
	assertCmp(t, want, state, "failed to match apps and services")
	goDemo := map[string]string{"app.kubernetes.io/name": "go-demo", "app.kubernetes.io/part-of": "go-demo"}

	extractAppAndServices(goDemo, state)

	want = map[string][]string{
		"go-demo": {"go-demo", "redis"},
	}
	assertCmp(t, want, state, "failed to match apps and services")
}

func TestExtractService(t *testing.T) {
	redisMap := map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"app.kubernetes.io/name":    "redis",
				"app.kubernetes.io/part-of": "go-demo",
			},
			"name":      "redis",
			"namespace": "test-env",
		},
		"spec": map[string]interface{}{
			"replicas": int64(1),
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"image": "redis:6-alpine",
							"name":  "redis",
							"ports": []interface{}{
								map[string]interface{}{
									"containerPort": 6379,
								},
							},
						},
					},
				},
			},
		},
	}

	svc := extractService(redisMap)
	want := &Service{
		Name:      "redis",
		Namespace: "test-env",
		Replicas:  1,
		Images:    []string{"redis:6-alpine"},
	}
	assertCmp(t, want, svc, "failed to match service")
}

func assertCmp(t *testing.T, want, got interface{}, msg string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf(msg+":\n%s", diff)
	}
}
