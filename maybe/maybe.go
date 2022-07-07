package maybe

import (
	"errors"
	"fmt"

	"github.com/obiloud/curry-go/debug"
	"github.com/obiloud/curry-go/nub"
)

type Maybe[T any] interface {
	IsJust() bool
	IsNothing() bool
	Unwrap() (T, error)
	String() string
}

type just[T any] struct {
	obj T
}

type nothing[T any] struct{}

func (j just[T]) IsJust() bool {
	return true
}

func (j just[T]) IsNothing() bool {
	return false
}

func (j just[T]) Unwrap() (T, error) {
	return j.obj, nil
}

func (j just[T]) String() string {
	return fmt.Sprintf("Just (%s)", debug.Stringify(j.obj))
}

func (n nothing[T]) IsJust() bool {
	return false
}

func (n nothing[T]) IsNothing() bool {
	return true
}

func (n nothing[T]) Unwrap() (T, error) {
	return *new(T), errors.New("Nothing has no value")
}

func (n nothing[T]) String() string {
	return "Nothing"
}

func Just[T any](x T) Maybe[T] {
	return just[T]{obj: x}
}

func Nothing[T any]() Maybe[T] {
	return nothing[T]{}
}

func WithDefault[T any](x T, maybe Maybe[T]) T {
	if maybe.IsNothing() {
		return x
	}
	return maybe.(just[T]).obj
}

func Map[A, B any](fn func(A) B, maybe Maybe[A]) Maybe[B] {
	switch val := maybe.(type) {
	case just[A]:
		return just[B]{obj: fn(val.obj)}
	}
	return Nothing[B]()
}

func Map2[A, B, C any](fn func(A, B) C, m1 Maybe[A], m2 Maybe[B]) Maybe[C] {
	if m1.IsJust() && m2.IsJust() {
		x, _ := m1.Unwrap()
		y, _ := m2.Unwrap()

		return Just(fn(x, y))
	}
	return Nothing[C]()
}

func Apply[A, B any](ma Maybe[func(A) B], mb Maybe[A]) Maybe[B] {
	return Map2(nub.ApplyFlipped[A, B], mb, ma)
}

func Bind[A, B any](fn func(A) Maybe[B], m Maybe[A]) Maybe[B] {
	if m.IsJust() {
		return fn(m.(just[A]).obj)
	}
	return Nothing[B]()
}
