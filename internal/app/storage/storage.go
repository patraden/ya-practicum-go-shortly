package storage

type BasicKVStorage interface {
	Add(key string, value string) (string, error)
	Get(key string) (string, error)
	Del(key string) error
}
