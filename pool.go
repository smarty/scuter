package scuter

import "sync"

// Pond is a generic wrapper over *sync.Pool, and sounds like an inviting place for creatures with scutes to hang out.
type Pond[T any] struct{ pool *sync.Pool }

func NewPond[T any](create func() T) *Pond[T] {
	return &Pond[T]{pool: &sync.Pool{New: func() any { return create() }}}
}
func (this *Pond[T]) Get() T  { return this.pool.Get().(T) }
func (this *Pond[T]) Put(t T) { this.pool.Put(t) }
