package pipeline

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-billy/v5"
)

var numericRe = regexp.MustCompile("^[0-9]+_")

// ListStages parses a directory withing a filesystem, and identifies the
// stages of a pipeline and the correct order.
//
// For ordering the stages in a pipeline, you can opt to add a numeric prefix,
// in this case, the numeric prefix will be stripped.
//
// e.g. 01_staging 02_production will be returned as staging, production.
func ListStages(fs billy.Filesystem, dir string) ([]string, error) {
	dirs, err := fs.ReadDir(dir)
	stages := []string{}
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %q: %w", dir, err)
	}
	for _, v := range dirs {
		if !v.IsDir() {
			continue
		}
		stages = append(stages, v.Name())
	}
	sort.Strings(stages)
	return trimNumericPrefixes(stages), nil
}

func trimNumericPrefixes(s []string) []string {
	trimmed := make([]string, len(s))
	for i := range s {

		p := numericRe.FindString(s[i])
		if p == "" {
			trimmed[i] = s[i]
			continue
		}
		trimmed[i] = strings.TrimPrefix(s[i], p)
	}
	return trimmed
}
