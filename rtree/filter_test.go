package rtree

import (
	"testing"

	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
)

func TestFilter(t *testing.T) {
	if Filter(nub.Const[bool, string](true), interestingTree) != maybe.Just(interestingTree) {
		t.Error("Filtering a Tree with a predicate that always returns true returns the same tree")
	}

	beforeE := func(elem string) bool {
		return elem < "e"
	}

	if Filter(beforeE, interestingTree) != maybe.Just(multiChildTree) {
		t.Error("Filtering a Tree with a predicate returns a filtered Tree")
	}

	isK := func(elem string) bool {
		return elem == "k"
	}

	if Filter(isK, interestingTree) != maybe.Nothing[RTree[string]]() {
		t.Error("If a subtree contains an element which would evaluate the predicate to True it is still not in the result Tree if the parent datum evaluates to false")
	}
}

func TestFilterWithChildPresedence(t *testing.T) {
	if FilterWithChildPresedence(nub.Const[bool, string](true), interestingTree) != maybe.Just(interestingTree) {
		t.Error("Filtering a Tree with a predicate that always returns true returns the same tree")
	}

	beforeE := func(elem string) bool {
		return elem < "e"
	}

	if FilterWithChildPresedence(beforeE, interestingTree) != maybe.Just(multiChildTree) {
		t.Error("Filtering a Tree with a predicate returns a filtered Tree")
	}

	foo := func(elem string) bool {
		return elem == "fooo"
	}

	if FilterWithChildPresedence(foo, interestingTree) != maybe.Nothing[RTree[string]]() {
		t.Error("If an element is no where to be found in the tree returns Nothing")
	}

	isK := func(elem string) bool {
		return elem == "k"
	}

	expected := maybe.Just(RTree[string]{
		Data: "a",
		Children: list.Singleton(
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
	})

	if FilterWithChildPresedence(isK, interestingTree) != expected {
		t.Error("If a predicate evaluates to False for a Node but True for one of it's children then the Node will remain in the Tree")
	}
}
