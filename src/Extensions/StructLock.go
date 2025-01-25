package Extensions

import "sync"

type StructLock[T any] struct {
	data *T
	lock sync.Mutex
}

func NewStructLock[T any](data T) *StructLock[T] {
	return &StructLock[T]{data: &data}
}

func (lock *StructLock[T]) Get() T {
	lock.lock.Lock()
	defer lock.lock.Unlock()
	return *lock.data
}

func (lock *StructLock[T]) Set(data T) {
	lock.lock.Lock()
	defer lock.lock.Unlock()
	*lock.data = data
}

func (lock *StructLock[T]) Update(f func(data *T)) {
	lock.lock.Lock()
	defer lock.lock.Unlock()
	f(lock.data)
}
