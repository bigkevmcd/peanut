package peanut

import (
	"testing"

	"github.com/bigkevmcd/peanut/pkg/parser"
)

func TestPromoteImage(t *testing.T) {
	_ = parser.Config{
		Apps: []*parser.App{
			{
				Name: "go-demo",
				Services: []*parser.Service{
					{Name: "go-demo-http", Replicas: 1, Images: []string{"bigkevmcd/go-demo:876ecb3"}},
					{Name: "redis", Replicas: 1, Images: []string{"redis:6-alpine"}},
				},
			},
		},
	}

}
