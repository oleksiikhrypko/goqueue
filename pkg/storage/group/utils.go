package group

import (
	"strings"
)

func buildStateKey(name []byte) []byte {
	const (
		pfx = "g:"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(name))
	s.WriteString(pfx)
	s.Write(name)
	return []byte(s.String())
}
