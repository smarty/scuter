package scuter

import "sync"

// Pool is a generic wrapper over *sync.Pool.
type Pool[T any] struct{ pool *sync.Pool }

func NewPool[T any](create func() T) *Pool[T] {
	return &Pool[T]{pool: &sync.Pool{New: func() any { return create() }}}
}
func (p *Pool[T]) Get() T  { return p.pool.Get().(T) }
func (p *Pool[T]) Put(t T) { p.pool.Put(t) }
