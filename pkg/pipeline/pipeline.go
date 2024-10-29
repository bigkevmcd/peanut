package pipeline

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strings"

	"github.com/go-git/go-billy/v5"
)

var numericRe = regexp.MustCompile("^[0-9]+_")

// ListStages parses a directory within a filesystem, and identifies the
// stages of a pipeline and the correct order.
//
// For ordering the stages in a pipeline, you can opt to add a numeric prefix,
// in this case, the numeric prefix will be stripped.
//
// e.g. 01_staging 02_production will be returned as staging, production.
func ListStages(fs billy.Filesystem, dir string) ([]string, error) {
	dirs, err := fs.ReadDir(dir)

	if err != nil && !isPathError(err) {
		return nil, fmt.Errorf("failed to read directory %q: %w", dir, err)
	}

	stages := []string{}
	for _, v := range dirs {
		if !v.IsDir() {
			continue
		}
		stages = append(stages, v.Name())
	}
	sort.Strings(stages)
	return trimNumericPrefixes(stages), nil
}

func isPathError(err error) bool {
	var pathErr *fs.PathError

	return errors.As(err, &pathErr)
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
