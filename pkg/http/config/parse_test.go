package config

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	parseTests := []struct {
		filename string
		want     *Config
	}{
		{"testdata/example1.yaml", &Config{
			Apps: []App{
				{
					Name: "go-demo",
					Environments: []Environment{
						{Name: "dev"},
						{Name: "staging"},
						{Name: "production"}},
				},
			},
		},
		},
	}

	for _, tt := range parseTests {
		t.Run(fmt.Sprintf("parsing %s", tt.filename), func(rt *testing.T) {
			got, err := ParseFile(tt.filename)
			if err != nil {
				rt.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				rt.Errorf("Parse(%s) failed diff\n%s", tt.filename, diff)
			}
		})
	}
}
