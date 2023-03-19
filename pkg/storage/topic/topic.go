package topic

type DB interface {
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	IsNotFoundErr(err error) bool
}

type Topic struct {
	db   DB
	name string
}

func New(db DB, name string) *Topic {
	return &Topic{
		db:   db,
		name: name,
	}
}
