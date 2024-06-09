package batch

type ActionType int

const (
	ActionTypeDel = ActionType(0)
	ActionTypeSet = ActionType(1)
)

func New(cap int) *Batch {
	if cap < 2 {
		cap = 2
	}
	return &Batch{
		acts:   make([]ActionType, 0, cap),
		keys:   make([][]byte, 0, cap),
		values: make([][]byte, 0, cap),
	}
}

type Batch struct {
	acts   []ActionType
	keys   [][]byte
	values [][]byte
}

func (b *Batch) ForEach(h func(action ActionType, key, value []byte)) {
	for i := range b.acts {
		h(b.acts[i], b.keys[i], b.values[i])
	}
}

func (b *Batch) appendRec(action ActionType, key, value []byte) {
	b.acts = append(b.acts, action)
	b.keys = append(b.keys, key)
	b.values = append(b.values, value)
}

func (b *Batch) AddAction(action ActionType, key, value []byte) {
	b.appendRec(action, key, value)
}

func (b *Batch) AddActionSet(key, value []byte) {
	b.appendRec(ActionTypeSet, key, value)
}

func (b *Batch) AddActionDel(key []byte) {
	b.appendRec(ActionTypeDel, key, nil)
}
