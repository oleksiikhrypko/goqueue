package batch

type BatchActionType int

const (
	BatchActionTypeDel = BatchActionType(0)
	BatchActionTypePut = BatchActionType(1)
)

func New(cap int) *Batch {
	if cap < 2 {
		cap = 2
	}
	return &Batch{
		acts:   make([]BatchActionType, 0, cap),
		keys:   make([][]byte, 0, cap),
		values: make([][]byte, 0, cap),
	}
}

type Batch struct {
	acts   []BatchActionType
	keys   [][]byte
	values [][]byte
}

func (b *Batch) ForEach(h func(action BatchActionType, key, value []byte)) {
	for i := range b.acts {
		h(b.acts[i], b.keys[i], b.values[i])
	}
}

func (b *Batch) appendRec(action BatchActionType, key, value []byte) {
	b.acts = append(b.acts, action)
	b.keys = append(b.keys, key)
	b.values = append(b.values, value)
}

func (b *Batch) Put(key, value []byte) {
	b.appendRec(BatchActionTypePut, key, value)
}

func (b *Batch) Delete(key []byte) {
	b.appendRec(BatchActionTypeDel, key, nil)
}
