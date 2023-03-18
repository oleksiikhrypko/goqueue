package topic

type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	IsNotFoundErr(err error) bool
}

type Topic struct {
	db   DB
	name []byte
}

func New(db DB, name []byte) *Topic {
	return &Topic{
		db:   db,
		name: name,
	}
}
