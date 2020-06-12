package config

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"sigs.k8s.io/yaml"
)

// Parse decodes YAML describing an environment manifest.
func Parse(in io.Reader) (*Config, error) {
	m := &Config{}
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// ParseFile is a wrapper around Parse that accepts a filename, it opens and
// parses the file, and closes it.
func ParseFile(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open: %s", filename)
	}
	defer f.Close()
	return Parse(f)
}
