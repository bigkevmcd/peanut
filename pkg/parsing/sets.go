package parser

import "sort"

type stringSet map[string]bool

func newStringSet() stringSet {
	return stringSet{}
}

// Add adds a string to the set.
func (s stringSet) Add(n string) {
	s[n] = true
}

// Elements returns a sorted list of elements in the set.
func (s stringSet) Elements() []string {
	e := []string{}
	for k := range s {
		e = append(e, k)
	}
	sort.Strings(e)
	return e
}
