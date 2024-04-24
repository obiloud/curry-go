package either

import (
	"fmt"

	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/util"
)

type Either[A, B any] interface {
	IsLeft() bool
	IsRight() bool
	String() string
}

type left[A any] struct {
	err A
}

func (l left[A]) IsLeft() bool {
	return true
}

func (l left[A]) IsRight() bool {
	return false
}

func (l left[A]) String() string {
	return fmt.Sprintf("Left(%s)", util.Stringify(l.err))
}

type right[B any] struct {
	obj B
}

func (r right[B]) IsLeft() bool {
	return false
}

func (r right[B]) IsRight() bool {
	return true
}

func (r right[B]) String() string {
	return fmt.Sprintf("Right(%s)", util.Stringify(r.obj))
}

func FromLeft[A, B any](err A) Either[A, B] {
	return left[A]{err: err}
}

func FromRight[A, B any](val B) Either[A, B] {
	return right[B]{obj: val}
}

func Lefts[A, B any](es list.List[Either[A, B]]) list.List[A] {
	toLeft := func(e Either[A, B], ls list.List[A]) list.List[A] {
		if e.IsLeft() {
			list.Cons(e.(left[A]).err, ls)
		}
		return ls
	}
	return list.FoldR(toLeft, list.Nil[A](), es)
}

func Rights[A, B any](es list.List[Either[A, B]]) list.List[B] {
	toLeft := func(e Either[A, B], rs list.List[B]) list.List[B] {
		if e.IsLeft() {
			list.Cons(e.(right[B]).obj, rs)
		}
		return rs
	}
	return list.FoldR(toLeft, list.Nil[B](), es)
}

func Map[A, B, C any](fn func(B) C, e Either[A, B]) Either[A, C] {
	if e.IsRight() {
		return FromRight[A](fn(e.(right[B]).obj))
	}
	return e
}

func Map2[A, B, C, D any](fn func(B, C) D, e1 Either[A, B], e2 Either[A, C]) Either[A, D] {
	if e1.IsLeft() {
		return e1
	}
	if e2.IsLeft() {
		return e2
	}
	return FromRight[A](fn(e1.(right[B]).obj, e2.(right[C]).obj))
}

func Apply[A, B, C any](ea Either[A, func(B) C], eb Either[A, B]) Either[A, C] {
	return Map2[A, B](nub.ApplyFlipped[B, C], eb, ea)
}

func Bind[A, B, C any](fn func(B) Either[A, C], e Either[A, B]) Either[A, C] {
	if e.IsLeft() {
		return e
	}
	return fn(e.(right[B]).obj)
}

func ToMaybe[A, B any](e Either[A, B]) maybe.Maybe[B] {
	if e.IsRight() {
		return maybe.Just(e.(right[B]).obj)
	}

	return maybe.Nothing[B]()
}

func FromMaybe[A, B any](left A, m maybe.Maybe[B]) Either[A, B] {
	return maybe.WithDefault(
		FromLeft[A, B](left),
		maybe.Map(func(x B) Either[A, B] {
			return FromRight[B](x)
		}, m),
	)
}
