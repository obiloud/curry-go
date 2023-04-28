package list

import (
	"fmt"
	"strings"

	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/tuple"
	"github.com/obiloud/curry-go/util"
)

type List[T any] interface {
	isCons() bool
	String() string
}

type nilList[T any] struct{}

type consList[T any] struct {
	head T
	tail List[T]
}

func (n nilList[T]) isCons() bool {
	return false
}

func (n nilList[T]) String() string {
	return "Nil"
}

func (c consList[T]) isCons() bool {
	return true
}

func (c consList[T]) String() string {
	toString := func(x T) string {
		return util.Stringify(x)
	}
	strs := Map[T](toString, c)
	slice := ToSlice(strs)
	return fmt.Sprintf("[%s]", strings.Join(slice, ", "))
}

// CREATE

func Nil[T any]() List[T] {
	return nilList[T]{}
}

func Singleton[T any](x T) List[T] {
	return Cons(x, Nil[T]())
}

func Cons[T any](head T, tail List[T]) List[T] {
	return consList[T]{
		head: head,
		tail: tail,
	}
}

func Range(lo int, hi int) List[int] {
	return rangeHelp(lo, hi, Nil[int]())
}

func rangeHelp(lo int, hi int, list List[int]) List[int] {
	if lo <= hi {
		return rangeHelp(lo, hi-1, Cons(hi, list))
	} else {
		return list
	}
}

func Repeat[T any](n int, value T) List[T] {
	return repeatHelp(Nil[T](), n, value)
}

func repeatHelp[T any](result List[T], n int, value T) List[T] {
	if n <= 0 {
		return result
	}
	return repeatHelp(Cons(value, result), n-1, value)
}

// UTILITIES

func Length[T any](list List[T]) int {
	if list.isCons() {
		return 1 + Length(list.(consList[T]).tail)
	}
	return 0
}

func Reverse[T any](list List[T]) List[T] {
	return FoldL(Cons[T], Nil[T](), list)
}

func FromSlice[T any](slice []T) List[T] {
	ls := Nil[T]()
	for i := len(slice) - 1; i >= 0; i-- {
		ls = Cons(slice[i], ls)
	}
	return ls
}

func ToSlice[T any](list List[T]) []T {
	con := func(h T, acc []T) []T {
		return append([]T{h}, acc...)
	}
	return FoldR(con, []T{}, list)
}

func Member[T comparable](x T, list List[T]) bool {
	eq := nub.Curry(nub.Eq[T])
	return Any(eq(x), list)
}

func All[T any](isOK func(T) bool, ls List[T]) bool {
	return nub.Not(Any(nub.Compose(nub.Not, isOK), ls))
}

func Any[T any](isOK func(T) bool, ls List[T]) bool {
	switch ls.(type) {
	case nilList[T]:
		return false
	}
	if isOK(ls.(consList[T]).head) {
		return true
	} else {
		return Any(isOK, ls.(consList[T]).tail)
	}
}

func Maximum[T nub.Num](xs List[T]) maybe.Maybe[T] {
	switch xs.(type) {
	case nilList[T]:
		return maybe.Nothing[T]()
	}
	return maybe.Just(FoldL(nub.Max[T], xs.(consList[T]).head, xs.(consList[T]).tail))
}

func Minimum[T nub.Num](xs List[T]) maybe.Maybe[T] {
	switch xs.(type) {
	case nilList[T]:
		return maybe.Nothing[T]()
	}
	return maybe.Just(FoldL(nub.Min[T], xs.(consList[T]).head, xs.(consList[T]).tail))
}

func Sum[T nub.Num](xs List[T]) T {
	add := func(x T, y T) T {
		return x + y
	}
	return FoldL(add, 0, xs)
}

func Product[T nub.Num](xs List[T]) T {
	mul := func(x T, y T) T {
		return x * y
	}
	return FoldL(mul, 1, xs)
}

// TRANSFORM

func Map[A, B any](fn func(A) B, ls List[A]) List[B] {
	fold := func(x A, acc List[B]) List[B] {
		return Cons(fn(x), acc)
	}
	return FoldR(fold, Nil[B](), ls)
}

func IndexedMap[A, B any](fn func(int, A) B, ls List[A]) List[B] {
	return Map2(fn, Range(0, Length(ls)-1), ls)
}

func FoldL[A, B any](fn func(A, B) B, acc B, ls List[A]) B {
	switch ls.(type) {
	case nilList[A]:
		return acc
	}
	return FoldL(fn, fn(ls.(consList[A]).head, acc), ls.(consList[A]).tail)
}

func FoldR[A, B any](fn func(A, B) B, acc B, ls List[A]) B {
	return foldRHelper(fn, acc, 0, ls)
}

func foldRHelper[A, B any](fn func(A, B) B, acc B, ctr int, a List[A]) B {
	switch a.(type) {
	case nilList[A]:
		return acc
	}

	b := a.(consList[A]).tail
	switch b.(type) {
	case nilList[A]:
		return fn(a.(consList[A]).head, acc)
	}

	c := b.(consList[A]).tail
	switch c.(type) {
	case nilList[A]:
		return fn(a.(consList[A]).head, fn(b.(consList[A]).head, acc))
	}

	d := c.(consList[A]).tail
	switch d.(type) {
	case nilList[A]:
		return fn(a.(consList[A]).head, fn(b.(consList[A]).head, fn(c.(consList[A]).head, acc)))
	}

	e := d.(consList[A]).tail
	var rest B
	if ctr > 500 {
		rest = FoldL(fn, acc, Reverse(e))
	} else {
		rest = foldRHelper(fn, acc, ctr+1, e)
	}
	return fn(a.(consList[A]).head, fn(b.(consList[A]).head, fn(c.(consList[A]).head, fn(d.(consList[A]).head, rest))))
}

func Filter[T any](fn func(T) bool, ls List[T]) List[T] {
	isGood := func(x T, acc List[T]) List[T] {
		if fn(x) {
			return Cons(x, acc)
		}
		return acc
	}
	return FoldR(isGood, Nil[T](), ls)
}

func FilterMap[A, B any](fn func(A) maybe.Maybe[B], ls List[A]) List[B] {
	maybeCons := func(x A, xs List[B]) List[B] {
		return maybe.WithDefault(xs, maybe.Map(func(v B) List[B] {
			return Cons(v, xs)
		}, fn(x)))
	}
	return FoldR(maybeCons, Nil[B](), ls)
}

// COMBINE

func Append[T any](xs List[T], ys List[T]) List[T] {
	switch ys.(type) {
	case nilList[T]:
		return xs
	}
	return FoldR(Cons[T], ys, xs)
}

func Concat[T any](lists List[List[T]]) List[T] {
	return FoldR(Append[T], Nil[T](), lists)
}

func ConcatMap[A, B any](fn func(A) List[B], ls List[A]) List[B] {
	return Concat(Map(fn, ls))
}

func Intersperse[T any](sep T, xs List[T]) List[T] {
	switch xs.(type) {
	case nilList[T]:
		return xs
	}

	step := func(x T, rest List[T]) List[T] {
		return Cons(sep, Cons(x, rest))
	}

	spersed := FoldR(step, Nil[T](), xs.(consList[T]).tail)

	return Cons(xs.(consList[T]).head, spersed)
}

func Map2[A, B, C any](fn func(A, B) C, xs List[A], ys List[B]) List[C] {
	slice1 := ToSlice(xs)
	slice2 := ToSlice(ys)
	results := []C{}
	for i := 0; i < len(slice1) && i < len(slice2); i++ {
		a := slice1[i]
		b := slice2[i]
		results = append(results, fn(a, b))
	}
	return FromSlice(results)
}

// SORT

func Sort[T nub.Ord](ls List[T]) List[T] {
	return SortBy(nub.Id[T], ls)
}

func SortBy[A any, B nub.Ord](fn func(A) B, ls List[A]) List[A] {
	slice := ToSlice(ls)

	sortFn := func(x A, y A) nub.Order {
		return nub.Compare(fn(x), fn(y))
	}

	quicksort(sortFn, slice)

	return FromSlice(slice)
}

func SortWith[T any](sortFn func(T, T) nub.Order, ls List[T]) List[T] {
	slice := ToSlice(ls)

	quicksort(sortFn, slice)

	return FromSlice(slice)
}

func quicksort[T any](sortFn func(T, T) nub.Order, xs []T) {
	quicksortHelp(sortFn, xs, 0, len(xs))
}

func quicksortHelp[T any](sortFn func(T, T) nub.Order, xs []T, lo int, hi int) {
	if lo < hi {
		pivotIndex := partition(sortFn, xs, lo, hi)
		quicksortHelp(sortFn, xs, lo, pivotIndex)
		quicksortHelp(sortFn, xs, pivotIndex+1, hi)
	}
}

func partition[T any](sortFn func(T, T) nub.Order, slice []T, lo int, hi int) int {
	swap := func(i int, j int) {
		tmp := slice[i]
		slice[i] = slice[j]
		slice[j] = tmp
	}

	pivot := slice[lo]
	x := lo + 1

	for i := x; i < hi; i++ {
		if sortFn(slice[i], pivot) == nub.LT {
			swap(i, x)
			x++
		}
	}
	swap(lo, x-1)
	return x - 1
}

// DECONSTRUCT

func IsEmpty[T any](list List[T]) bool {
	return Length(list) == 0
}

func Head[T any](list List[T]) maybe.Maybe[T] {
	if list.isCons() {
		return maybe.Just(list.(consList[T]).head)
	}
	return maybe.Nothing[T]()
}

func Tail[T any](list List[T]) maybe.Maybe[List[T]] {
	if list.isCons() {
		return maybe.Just(list.(consList[T]).tail)
	}
	return maybe.Nothing[List[T]]()
}

func Take[T any](n int, ls List[T]) List[T] {
	if n <= 0 {
		return Nil[T]()
	} else {
		switch ls.(type) {
		case nilList[T]:
			return ls
		}

		slice := ToSlice(ls)

		if n < len(slice) {
			slice = slice[:n]
		}

		return FromSlice(slice)
	}
}

func Drop[T any](n int, ls List[T]) List[T] {
	if n <= 0 {
		return ls
	} else {
		switch ls.(type) {
		case nilList[T]:
			return ls
		}

		return Drop(n-1, ls.(consList[T]).tail)
	}
}

func Partition[T any](predicate func(T) bool, ls List[T]) tuple.Tuple[List[T], List[T]] {
	step := func(x T, pair tuple.Tuple[List[T], List[T]]) tuple.Tuple[List[T], List[T]] {
		if predicate(x) {
			return tuple.MapFirst(nub.Curry(Cons[T])(x), pair)
		}
		return tuple.MapSecond(nub.Curry(Cons[T])(x), pair)
	}
	return FoldR(step, tuple.Pair(Nil[T](), Nil[T]()), ls)
}

func Unzip[A, B any](xs List[tuple.Tuple[A, B]]) tuple.Tuple[List[A], List[B]] {
	step := func(p1 tuple.Tuple[A, B], p2 tuple.Tuple[List[A], List[B]]) tuple.Tuple[List[A], List[B]] {
		return tuple.MapBoth(
			nub.Curry(Cons[A])(tuple.First(p1)),
			nub.Curry(Cons[B])(tuple.Second(p1)),
			p2,
		)
	}
	return FoldR(step, tuple.Pair(Nil[A](), Nil[B]()), xs)
}
