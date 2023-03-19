package klist

import (
	"bytes"
	"strings"

	models "goqueue/pkg/proto/models"
)

func buildItemKey(name, item []byte) []byte {
	const (
		pfx = "qi:"
		sep = ":>"
	)
	var s strings.Builder
	s.Grow(len(pfx) + len(name) + len(sep) + len(item))
	s.WriteString(pfx)
	s.Write(name)
	s.WriteString(sep)
	s.Write(item)
	return []byte(s.String())
}

func buildStateKey(name []byte) []byte {
	const (
		pfx = "q:"
	)
	s := strings.Builder{}
	s.Grow(len(pfx) + len(name))
	s.WriteString(pfx)
	s.Write(name)
	return []byte(s.String())
}

func isEqual(item1, item2 []byte) bool {
	return bytes.Compare(item1, item2) == 0
}

func isItemFirst(state *models.KList, item []byte) bool {
	return isEqual(state.GetFirstItem(), item)
}

func isItemLast(state *models.KList, item []byte) bool {
	return isEqual(state.GetLastItem(), item)
}

func isEmpty(state *models.KList) bool {
	return state.Count == 0
}
