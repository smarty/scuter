package scuter

import "net/http"

type Logger interface {
	Printf(string, ...any)
}

type PooledModelFramework[T any] struct {
	logger Logger
	pool   *Pool[T]
	reset  func(T)
}

func NewPooledModelFramework[T any](logger Logger, create func() T, reset func(T)) *PooledModelFramework[T] {
	return &PooledModelFramework[T]{logger: logger, pool: NewPool[T](create), reset: reset}
}

func (this *PooledModelFramework[T]) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	model := this.pool.Get()
	this.reset(model)
	defer this.pool.Put(model)
	err := Flush(response, this.serveHTTP(request, model))
	if err != nil {
		this.logger.Printf("[WARN] JSON serialization error: %v", err)
	}
}
func (this *PooledModelFramework[T]) serveHTTP(_ *http.Request, _ T) ResponseOption {
	panic("must be overridden by end-user implementation")
}

type Framework struct {
	logger Logger
}

func NewFramework(logger Logger) *Framework {
	return &Framework{logger: logger}
}

func (this *Framework) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	err := Flush(response, this.serveHTTP(request))
	if err != nil {
		this.logger.Printf("[WARN] JSON serialization error: %v", err)
	}
}
func (this *Framework) serveHTTP(_ *http.Request) ResponseOption {
	panic("must be overridden by end-user implementation")
}
