package db

type DBEngine interface {
	Open() (err error)
	Close() (err error)
	Has(key []byte) (b bool)
	Put(key, value []byte) (err error)
	Get(key []byte) (value []byte, err error)
	Del(key []byte) (err error)
	GetString(key []byte) (value string, err error)
	Batch(put map[string][]byte, del []string) (err error)
	GetLike(prefix []byte) (datamap map[string][]byte, err error)
	GetKeys() (bys []string, err error)
	GetKeysPrefix(prefix []byte) (bys []string, err error)
	GetKeysPrefixLimit(prefix []byte, limit int) (bys []string, err error)
	GetIterLimit(prefix string, limit string) (datamap map[string][]byte, err error)
}
