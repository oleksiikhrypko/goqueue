package klist

import (
	"bytes"
	"strings"
)

func buildItemKey(name, item []byte) []byte {
	var s strings.Builder
	s.Write(name)
	s.WriteString(":")
	s.Write(item)
	s.WriteString(":")
	return []byte(s.String())
}

func buildStateKey(name []byte) []byte {
	s := strings.Builder{}
	s.Write(name)
	s.WriteString(":state")
	return []byte(s.String())
}

func isEqual(item1, item2 []byte) bool {
	return bytes.Compare(item1, item2) == 0
}
