package rtree

import (
	"fmt"

	"github.com/obiloud/curry-go/debug"
	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/tuple"
)

type RTree[T any] struct {
	Data     T
	Children list.List[RTree[T]]
}

func (t RTree[T]) String() string {
	return fmt.Sprintf("Tree %s %s", debug.Stringify(t.Data), t.Children.String())
}

func InsertChild[T any](child RTree[T], tree RTree[T]) RTree[T] {
	tree.Children = list.Cons(child, tree.Children)
	return tree
}

func AppendChild[T any](child RTree[T], tree RTree[T]) RTree[T] {
	tree.Children = list.Append(tree.Children, list.Singleton(child))
	return tree
}

func FoldL[A, B any](fn func(A, B) B, acc B, tree RTree[A]) B {
	var treeUnwrap func(RTree[A], B) B
	treeUnwrap = func(t RTree[A], accu B) B {
		return list.FoldL(treeUnwrap, fn(t.Data, accu), t.Children)
	}
	return list.FoldL(treeUnwrap, fn(tree.Data, acc), tree.Children)
}

func FoldR[A, B any](fn func(A, B) B, acc B, tree RTree[A]) B {
	var treeUnwrap func(RTree[A], B) B
	treeUnwrap = func(t RTree[A], accu B) B {
		return fn(t.Data, list.FoldR(treeUnwrap, accu, t.Children))
	}
	return fn(tree.Data, list.FoldR(treeUnwrap, acc, tree.Children))
}

func Flatten[T any](tree RTree[T]) list.List[T] {
	cons := func(x T, xs list.List[T]) list.List[T] {
		return list.Cons(x, xs)
	}
	return FoldR(cons, list.Nil[T](), tree)
}

func TuplesOfDatumAndFlatChildren[T any](tree RTree[T]) list.List[tuple.Tuple[T, list.List[T]]] {
	return list.Append(
		list.Singleton(
			tuple.Pair(tree.Data, list.ConcatMap(Flatten[T], tree.Children)),
		),
		list.ConcatMap(TuplesOfDatumAndFlatChildren[T], tree.Children),
	)
}

func Length[T any](tree RTree[T]) int {
	count := func(_ T, acc int) int {
		return acc + 1
	}
	return FoldR(count, 0, tree)
}

func Map[A, B any](fn func(A) B, tree RTree[A]) RTree[B] {
	transformTree := func(c RTree[A]) RTree[B] {
		return Map(fn, c)
	}
	return RTree[B]{
		Data:     fn(tree.Data),
		Children: list.Map(transformTree, tree.Children),
	}
}

func MapListOverTree[A, B, C any](fn func(A, B) C, ls list.List[A], tree RTree[B]) maybe.Maybe[RTree[C]] {
	if list.IsEmpty(ls) {
		return maybe.Nothing[RTree[C]]()
	}
	if list.Length(ls) == 1 {
		mapData := func(head A) RTree[C] {
			return RTree[C]{
				Data:     fn(head, tree.Data),
				Children: list.Nil[RTree[C]](),
			}

		}
		return maybe.Map(mapData, list.Head(ls))
	}

	transform := func(head A, tail list.List[A]) RTree[C] {
		mappedDatum := fn(head, tree.Data)

		lengths := list.Map(Length[B], tree.Children)

		listGroupedByLengthOfChildren := splitByLength(lengths, tail)

		mappedChildren := list.Map2(func(l list.List[A], child RTree[B]) maybe.Maybe[RTree[C]] {
			return MapListOverTree(fn, l, child)
		}, listGroupedByLengthOfChildren, tree.Children)

		return RTree[C]{
			Data:     mappedDatum,
			Children: list.FilterMap(nub.Id[maybe.Maybe[RTree[C]]], mappedChildren),
		}
	}

	return maybe.Map2(transform, list.Head(ls), list.Tail(ls))
}

func splitByLength[T any](listOfLengths list.List[int], ls list.List[T]) list.List[list.List[T]] {
	return splitByLengthHelper(listOfLengths, ls, list.Nil[list.List[T]]())
}

func splitByLengthHelper[T any](listOfLengths list.List[int], ls list.List[T], acc list.List[list.List[T]]) list.List[list.List[T]] {
	if list.IsEmpty(listOfLengths) {
		return list.Reverse(acc)
	}

	return maybe.WithDefault(
		list.Reverse(acc),
		maybe.Map2(func(currentLength int, restLengths list.List[int]) list.List[list.List[T]] {
			if list.IsEmpty(ls) {
				return list.Reverse(acc)
			}
			return splitByLengthHelper(restLengths, list.Drop(currentLength, ls), list.Cons(list.Take(currentLength, ls), acc))

		}, list.Head(listOfLengths), list.Tail(listOfLengths)))
}

func IndexedMap[A, B any](fn func(int, A) B, tree RTree[A]) maybe.Maybe[RTree[B]] {
	return MapListOverTree(fn, list.Range(0, Length(tree)-1), tree)
}

func Filter[T any](predicate func(T) bool, tree RTree[T]) maybe.Maybe[RTree[T]] {
	if predicate(tree.Data) {
		return maybe.Just(RTree[T]{
			Data:     tree.Data,
			Children: list.FilterMap(nub.Curry(Filter[T])(predicate), tree.Children),
		})
	}

	return maybe.Nothing[RTree[T]]()
}

func FilterWithChildPresedence[T any](predicate func(T) bool, tree RTree[T]) maybe.Maybe[RTree[T]] {
	fChildren := list.FilterMap(nub.Curry(FilterWithChildPresedence[T])(predicate), tree.Children)
	if list.IsEmpty(fChildren) {
		if predicate(tree.Data) {
			return maybe.Just(RTree[T]{
				Data:     tree.Data,
				Children: list.Nil[RTree[T]](),
			})
		} else {
			return maybe.Nothing[RTree[T]]()
		}
	}

	return maybe.Just(RTree[T]{
		Data:     tree.Data,
		Children: fChildren,
	})
}

func SortBy[A any, B nub.Ord](fn func(A) B, tree RTree[A]) RTree[A] {
	sortedChildren := list.Map(nub.Curry(SortBy[A, B])(fn), list.SortBy(func(t RTree[A]) B {
		return fn(t.Data)
	}, tree.Children))

	return RTree[A]{
		Data:     tree.Data,
		Children: sortedChildren,
	}
}

func SortWith[T any](comperator func(T, T) nub.Order, tree RTree[T]) RTree[T] {
	sortedChildren := list.Map(nub.Curry(SortWith[T])(comperator), list.SortWith(func(a RTree[T], b RTree[T]) nub.Order {
		return comperator(a.Data, b.Data)
	}, tree.Children))

	return RTree[T]{
		Data:     tree.Data,
		Children: sortedChildren,
	}
}
