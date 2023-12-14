package batch

type List interface {
	ForEach(h func(action ActionType, key, value []byte))
	AppendPut(key, value []byte)
	AppendDelete(key []byte)
}
