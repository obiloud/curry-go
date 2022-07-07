package rtree

import (
	"testing"

	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/nub"
)

var unorderedTree RTree[string] = RTree[string]{
	Data: "a",
	Children: list.Cons(
		RTree[string]{
			Data: "c",
			Children: list.Cons(
				RTree[string]{
					Data:     "g",
					Children: list.Nil[RTree[string]](),
				},
				list.Singleton(RTree[string]{
					Data:     "f",
					Children: list.Nil[RTree[string]](),
				}),
			),
		},
		list.Cons(
			RTree[string]{
				Data: "b",
				Children: list.Singleton(
					RTree[string]{
						Data: "e",
						Children: list.Singleton(
							RTree[string]{
								Data:     "k",
								Children: list.Nil[RTree[string]](),
							},
						),
					},
				),
			},
			list.Singleton(
				RTree[string]{
					Data: "d",
					Children: list.Cons(
						RTree[string]{
							Data:     "i",
							Children: list.Nil[RTree[string]](),
						},
						list.Cons(
							RTree[string]{
								Data:     "h",
								Children: list.Nil[RTree[string]](),
							},
							list.Singleton(
								RTree[string]{
									Data:     "j",
									Children: list.Nil[RTree[string]](),
								},
							),
						),
					),
				},
			),
		),
	),
}

var reverseSortedTree RTree[string] = RTree[string]{
	Data: "a",
	Children: list.Cons(
		RTree[string]{
			Data: "d",
			Children: list.Cons(
				RTree[string]{
					Data:     "j",
					Children: list.Nil[RTree[string]](),
				},
				list.Cons(
					RTree[string]{
						Data:     "i",
						Children: list.Nil[RTree[string]](),
					},
					list.Singleton(
						RTree[string]{
							Data:     "h",
							Children: list.Nil[RTree[string]](),
						},
					),
				))},
		list.Cons(
			RTree[string]{
				Data: "c",
				Children: list.Cons(
					RTree[string]{
						Data:     "g",
						Children: list.Nil[RTree[string]](),
					},
					list.Singleton(
						RTree[string]{
							Data:     "f",
							Children: list.Nil[RTree[string]](),
						},
					),
				)},
			list.Singleton(
				RTree[string]{
					Data: "b",
					Children: list.Singleton(
						RTree[string]{
							Data: "e",
							Children: list.Singleton(
								RTree[string]{
									Data:     "k",
									Children: list.Nil[RTree[string]](),
								},
							),
						},
					),
				},
			),
		),
	),
}

func flippedComparison[T nub.Ord](a T, b T) nub.Order {
	compared := nub.Compare(a, b)

	if compared == nub.LT {
		return nub.GT
	}
	if compared == nub.GT {
		return nub.LT
	}
	return nub.EQ
}

func TestSort(t *testing.T) {
	if SortBy(nub.Id[string], deepTree) != deepTree {
		t.Error("Sorting a Tree with only one child per levels yields the same Tree")
	}

	if SortBy(nub.Id[string], interestingTree) != interestingTree {
		t.Error("Sorting a sorted Tree returns the same Tree")
	}

	if SortBy(nub.Id[string], unorderedTree) != interestingTree {
		t.Error("Sorting an unsorted Tree returns a sorted Tree")
	}

	if SortWith(flippedComparison[string], interestingTree) != reverseSortedTree {
		t.Error("Sorting with a Tree with a reversed comperator reverse-sorts a Tree")
	}
}
