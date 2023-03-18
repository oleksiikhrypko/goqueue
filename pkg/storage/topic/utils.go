package topic

import (
	"strings"
)

func buildStateKey(name []byte) []byte {
	const (
		pfx = "t:"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(name))
	s.WriteString(pfx)
	s.Write(name)
	return []byte(s.String())
}
