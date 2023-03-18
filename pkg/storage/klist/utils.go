package klist

import (
	"bytes"
	"strings"

	models "goqueue/pkg/proto/klist"
)

func buildItemKey(name, item []byte) []byte {
	var s strings.Builder
	s.Grow(len(name) + len(item) + 2)
	s.Write(name)
	s.WriteString(":")
	s.Write(item)
	s.WriteString(":")
	return []byte(s.String())
}

func buildStateKey(name []byte) []byte {
	s := strings.Builder{}
	s.Grow(len(name))
	s.Write(name)
	s.WriteString(":state")
	return []byte(s.String())
}

func isEqual(item1, item2 []byte) bool {
	return bytes.Compare(item1, item2) == 0
}

func isItemFirst(state *models.State, item []byte) bool {
	return isEqual(state.GetFirstItem(), item)
}

func isItemLast(state *models.State, item []byte) bool {
	return isEqual(state.GetLastItem(), item)
}

func isEmpty(state *models.State) bool {
	return state.Count == 0
}
