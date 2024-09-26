package result

import (
	"errors"
	"fmt"
)

// Result is a generic interface that can represent either a success or an error.
type Result[T any] interface {
	// IsSuccess returns true if the result is a success, false otherwise.
	IsSuccess() bool
	// IsError returns true if the result is an error, false otherwise.
	IsError() bool
	// Value returns the success value and error value if there is error any.
	Value() (*T, error)
	// Error returns the error value.
	Error() error
}

// success represents a successful result.
type Success[T any] struct {
	Val *T
}

// NewSuccess creates a new success result.
func NewSuccess[T any](value *T) Result[T] {
	return &Success[T]{Val: value}
}

// IsSuccess implements Result.
func (s *Success[T]) IsSuccess() bool {
	return true
}

// Value implements Result.
func (s *Success[T]) Value() (*T, error) {
	return s.Val, nil
}

// IsError implements Result.
func (s *Success[T]) IsError() bool {
	return false
}

// Error implements Result.
func (s *Success[T]) Error() error {
	return errors.New("Cannot get error from a success result")
}

// error represents an error result.
type Failure[T any] struct {
	Err error
}

// NewError creates a new error result.
func NewError[T any](err error) Result[T] {
	return &Failure[T]{Err: err}
}

// IsSuccess implements Result.
func (f *Failure[T]) IsSuccess() bool {
	return false
}

// IsError implements Result.
func (f *Failure[T]) IsError() bool {
	return true
}

// Value implements Result.
func (f *Failure[T]) Value() (*T, error) {
	return nil, f.Error()
}

// Error implements Result.
func (f *Failure[T]) Error() error {
	return f.Err
}
func ToResult[T any](value *T, err error) Result[T] {
	if err != nil {
		return NewError[T](err)
	}
	return NewSuccess(value)
}

// CastError attempts to cast the error to a specific type E and returns a new Result.
func CastError[T, E any](r Result[T]) Result[E] {
	if r.IsSuccess() {
		return NewError[E](fmt.Errorf("cannot cast a success result"))
	}
	_, err := r.Value()
	return NewError[E](err)
}

// MapError maps the error of a Result to a new Result with a different type.
func MapError[T, R any](r Result[T], mapFn func(error) error) Result[R] {
	if r.IsSuccess() {
		return NewError[R](fmt.Errorf("cannot map a success result"))
	}
	return NewError[R](mapFn(r.Error()))
}
