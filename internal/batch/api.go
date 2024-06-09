package batch

type ActionsList interface {
	ForEach(h func(action ActionType, key, value []byte))
	AddActionSet(key, value []byte)
	AddActionDel(key []byte)
}
