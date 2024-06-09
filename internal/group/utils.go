package group

import (
	"strings"
)

func buildStateKey(topic, name string) []byte {
	const (
		pfx = "g:"
		sep = ":>"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(topic) + len(sep) + len(name))
	s.WriteString(pfx)
	s.WriteString(name)
	return []byte(s.String())
}
