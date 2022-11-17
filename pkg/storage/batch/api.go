package batch

type List interface {
	ForEach(h func(action BatchActionType, key, value []byte))
}
