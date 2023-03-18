package batch

type List interface {
	ForEach(h func(action BatchActionType, key, value []byte))
	Put(key, value []byte)
	Delete(key []byte)
}
