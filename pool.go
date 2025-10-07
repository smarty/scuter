package scuter

import "sync"

// Pool is a generic wrapper over *sync.Pool, and sounds like an inviting place for creatures with scutes to hang out.
type Pool[T any] struct{ pool *sync.Pool }

func NewPool[T any](create func() T) *Pool[T] {
	return &Pool[T]{pool: &sync.Pool{New: func() any { return create() }}}
}
func (this *Pool[T]) Get() T  { return this.pool.Get().(T) }
func (this *Pool[T]) Put(t T) { this.pool.Put(t) }
