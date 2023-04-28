package tuple

import (
	"fmt"

	"github.com/obiloud/curry-go/util"
)

type Tuple[A, B any] struct {
	first  A
	second B
}

func Pair[A, B any](first A, second B) Tuple[A, B] {
	return Tuple[A, B]{
		first:  first,
		second: second,
	}
}

func (pair Tuple[A, B]) String() string {
	return fmt.Sprintf("(%s, %s)", util.Stringify(pair.first), util.Stringify(pair.second))
}

func First[A, B any](tuple Tuple[A, B]) A {
	return tuple.first
}

func Second[A, B any](tuple Tuple[A, B]) B {
	return tuple.second
}

func MapFirst[A, B, C any](fn func(A) C, pair Tuple[A, B]) Tuple[C, B] {
	first := fn(pair.first)

	return Pair(first, pair.second)
}

func MapSecond[A, B, C any](fn func(B) C, pair Tuple[A, B]) Tuple[A, C] {
	second := fn(pair.second)

	return Pair(pair.first, second)
}

func MapBoth[A, B, C, D any](f1 func(A) C, f2 func(B) D, pair Tuple[A, B]) Tuple[C, D] {
	first := f1(pair.first)
	second := f2(pair.second)

	return Pair(first, second)
}
