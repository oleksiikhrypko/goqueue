package topic

import (
	"strings"
)

func buildKey(name []byte) []byte {
	const (
		pfx = "t:"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(name))
	s.WriteString(pfx)
	s.Write(name)
	return []byte(s.String())
}

func buildSequenceListName(topic, key []byte) []byte {
	const (
		pfx = "t:"
		sep = ":>"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(topic) + len(sep) + len(key))
	s.WriteString(pfx)
	s.Write(topic)
	s.WriteString(sep)
	s.Write(key)
	return []byte(s.String())
}
