package registry

import "errors"

type Registry[K comparable, V any] map[K]V

var ErrExists = errors.New("exists")
var ErrNotFound = errors.New("not found")

func (r Registry[K, V]) Register(k K, v V) error {
	if r.Has(k) {
		return ErrExists
	}
	r[k] = v
	return nil
}

func (r Registry[K, V]) Has(k K) bool {
	_, found := r[k]
	return found
}

func (r Registry[K, V]) Get(k K) (V, error) {
	v, exists := r[k]
	if !exists {
		return v, ErrNotFound
	}
	return v, nil
}
