package topic

import (
	"strings"
)

func buildKey(name string) []byte {
	const (
		pfx = "t:"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(name))
	s.WriteString(pfx)
	s.WriteString(name)
	return []byte(s.String())
}

func BuildSubListName(topic string, subKey []byte) []byte {
	const (
		pfx = "t:"
		sep = ":>"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(topic) + len(sep) + len(subKey))
	s.WriteString(pfx)
	s.WriteString(topic)
	s.WriteString(sep)
	s.Write(subKey)
	return []byte(s.String())
}
