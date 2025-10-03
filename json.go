package scuter

type JSON[T any] struct {
	JSONRequest
	JSONResponder[T]
}
