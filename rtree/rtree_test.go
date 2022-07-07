package rtree

import (
	"testing"

	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/tuple"
)

var interestingTree RTree[string] = RTree[string]{
	Data: "a",
	Children: list.Cons(
		RTree[string]{
			Data: "b",
			Children: list.Cons(RTree[string]{
				Data: "e",
				Children: list.Cons(RTree[string]{
					Data:     "k",
					Children: list.Nil[RTree[string]](),
				}, list.Nil[RTree[string]]()),
			}, list.Nil[RTree[string]]()),
		},
		list.Cons(
			RTree[string]{
				Data: "c",
				Children: list.Cons(RTree[string]{
					Data:     "f",
					Children: list.Nil[RTree[string]](),
				}, list.Cons(RTree[string]{
					Data:     "g",
					Children: list.Nil[RTree[string]](),
				}, list.Nil[RTree[string]]())),
			},
			list.Cons(
				RTree[string]{
					Data: "d",
					Children: list.Cons(RTree[string]{
						Data:     "h",
						Children: list.Nil[RTree[string]](),
					}, list.Cons(RTree[string]{
						Data:     "i",
						Children: list.Nil[RTree[string]](),
					}, list.Cons(RTree[string]{
						Data:     "j",
						Children: list.Nil[RTree[string]](),
					}, list.Nil[RTree[string]]()))),
				}, list.Nil[RTree[string]]()),
		),
	),
}

var noChildTree RTree[string] = RTree[string]{
	Data:     "a",
	Children: list.Nil[RTree[string]](),
}

var multiChildTree RTree[string] = RTree[string]{
	Data: "a",
	Children: list.Cons(RTree[string]{
		Data:     "b",
		Children: list.Nil[RTree[string]](),
	}, list.Cons(RTree[string]{
		Data:     "c",
		Children: list.Nil[RTree[string]](),
	}, list.Cons(RTree[string]{
		Data:     "d",
		Children: list.Nil[RTree[string]](),
	}, list.Nil[RTree[string]]()))),
}

var singleChildTree RTree[string] = RTree[string]{
	Data:     "a",
	Children: list.Singleton(RTree[string]{Data: "b", Children: list.Nil[RTree[string]]()}),
}

var deepTree RTree[string] = RTree[string]{
	Data: "a",
	Children: list.Cons(RTree[string]{
		Data: "b",
		Children: list.Cons(RTree[string]{
			Data: "c",
			Children: list.Cons(RTree[string]{
				Data:     "d",
				Children: list.Nil[RTree[string]](),
			}, list.Nil[RTree[string]]()),
		}, list.Nil[RTree[string]]()),
	}, list.Nil[RTree[string]]()),
}

func TestFindingElement(t *testing.T) {
	in := maybe.Just(Zipper[string]{
		Tree:        interestingTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	foo := func(elem string) bool {
		return elem == "FOO"
	}

	r1 := maybe.Bind(nub.Curry(GoTo[string])(foo), in)

	if r1 != maybe.Nothing[Zipper[string]]() {
		t.Error("Trying to find a non existing element in a Tree returns Nothing")
	}

	e := maybe.Bind(
		nub.Curry(GoToChild[string])(0),
		maybe.Bind(
			nub.Curry(GoToChild[string])(2),
			in,
		),
	)

	findh := func(elem string) bool {
		return elem == "h"
	}

	r2 := maybe.Bind(nub.Curry(GoTo[string])(findh), in)

	if r2 != e {
		t.Error("Trying to find an existing element in a Tree moves the focus to this element")
	}
}

func TestInsert(t *testing.T) {
	expected := maybe.Just(Zipper[string]{
		Tree:        interestingTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r1 := maybe.Bind(nub.Curry(GoToChild[string])(0), maybe.Just(Zipper[string]{
		Tree:        multiChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	}))

	r2 := maybe.Bind(nub.Curry(InsertChildTree[string])(RTree[string]{Data: "e", Children: list.Nil[RTree[string]]()}), r1)
	// debug.Debug("r2", r2)
	r3 := maybe.Bind(nub.Curry(GoToChild[string])(0), r2)
	// debug.Debug("r3", r3)
	r4 := maybe.Bind(nub.Curry(InsertChildTree[string])(RTree[string]{Data: "k", Children: list.Nil[RTree[string]]()}), r3)
	// debug.Debug("r4", r4)
	r5 := maybe.Bind(GoUp[string], r4)
	// debug.Debug("r5", r5)
	r6 := maybe.Bind(GoRight[string], r5)
	// debug.Debug("r6", r6)
	r7 := maybe.Bind(nub.Curry(InsertChildTree[string])(RTree[string]{Data: "g", Children: list.Nil[RTree[string]]()}), r6)
	// debug.Debug("r7", r7)
	r8 := maybe.Bind(nub.Curry(InsertChildTree[string])(RTree[string]{Data: "f", Children: list.Nil[RTree[string]]()}), r7)
	// debug.Debug("r8", r8)
	r9 := maybe.Bind(GoRight[string], r8)
	// debug.Debug("r9", r9)
	r10 := maybe.Bind(nub.Curry(InsertChildTree[string])(RTree[string]{Data: "j", Children: list.Nil[RTree[string]]()}), r9)
	// debug.Debug("r10", r10)
	r11 := maybe.Bind(nub.Curry(InsertChildTree[string])(RTree[string]{Data: "i", Children: list.Nil[RTree[string]]()}), r10)
	// debug.Debug("r11", r11)
	r12 := maybe.Bind(nub.Curry(InsertChildTree[string])(RTree[string]{Data: "h", Children: list.Nil[RTree[string]]()}), r11)
	// debug.Debug("r12", r12)
	r13 := maybe.Bind(GoToRoot[string], r12)
	// debug.Debug("r13", r13)

	if r13 != expected {
		t.Error("Inserting children can turn a multiChildTree into an interestingTree")
	}
}

func TestAppend(t *testing.T) {
	expected := maybe.Just(Zipper[string]{
		Tree:        interestingTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r1 := maybe.Bind(nub.Curry(GoToChild[string])(0), maybe.Just(Zipper[string]{
		Tree:        multiChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	}))

	r2 := maybe.Bind(nub.Curry(AppendChildTree[string])(RTree[string]{Data: "e", Children: list.Nil[RTree[string]]()}), r1)
	// debug.Debug("r2", r2)
	r3 := maybe.Bind(nub.Curry(GoToChild[string])(0), r2)
	// debug.Debug("r3", r3)
	r4 := maybe.Bind(nub.Curry(AppendChildTree[string])(RTree[string]{Data: "k", Children: list.Nil[RTree[string]]()}), r3)
	// debug.Debug("r4", r4)
	r5 := maybe.Bind(GoUp[string], r4)
	// debug.Debug("r5", r5)
	r6 := maybe.Bind(GoRight[string], r5)
	// debug.Debug("r6", r6)
	r7 := maybe.Bind(nub.Curry(AppendChildTree[string])(RTree[string]{Data: "f", Children: list.Nil[RTree[string]]()}), r6)
	// debug.Debug("r7", r7)
	r8 := maybe.Bind(nub.Curry(AppendChildTree[string])(RTree[string]{Data: "g", Children: list.Nil[RTree[string]]()}), r7)
	// debug.Debug("r8", r8)
	r9 := maybe.Bind(GoRight[string], r8)
	// debug.Debug("r9", r9)
	r10 := maybe.Bind(nub.Curry(AppendChildTree[string])(RTree[string]{Data: "h", Children: list.Nil[RTree[string]]()}), r9)
	// debug.Debug("r10", r10)
	r11 := maybe.Bind(nub.Curry(AppendChildTree[string])(RTree[string]{Data: "i", Children: list.Nil[RTree[string]]()}), r10)
	// debug.Debug("r11", r11)
	r12 := maybe.Bind(nub.Curry(AppendChildTree[string])(RTree[string]{Data: "j", Children: list.Nil[RTree[string]]()}), r11)
	// debug.Debug("r12", r12)
	r13 := maybe.Bind(GoToRoot[string], r12)
	// debug.Debug("r13", r13)

	if r13 != expected {
		t.Error("appending children can turn a multiChildTree into an interestingTree")
	}
}

func TestFlatten(t *testing.T) {
	exp1 := list.Cons("a", list.Cons("b", list.Cons("c", list.Singleton("d"))))
	exp2 := list.FromSlice([]string{"a", "b", "e", "k", "c", "f", "g", "d", "h", "i", "j"})

	if Flatten(multiChildTree) != exp1 {
		t.Error("Flatten multiChildTree")
	}

	if Flatten(deepTree) != exp1 {
		t.Error("Flatten deepTree")
	}

	if Flatten(interestingTree) != exp2 {
		t.Error("Flatten interestingTree")
	}
}

func TestFoldL(t *testing.T) {
	cons := func(x string, xs list.List[string]) list.List[string] {
		return list.Cons(x, xs)
	}
	if Flatten(interestingTree) != list.Reverse(FoldL(cons, list.Nil[string](), interestingTree)) {
		t.Error("Foldl interestingTree into List")
	}
}

func TestLength(t *testing.T) {
	if Length(interestingTree) != 11 {
		t.Error("Length of an interesting Tree")
	}

	if Length(noChildTree) != 1 {
		t.Error("Length of a noChildTree")
	}

	if Length(deepTree) != 4 {
		t.Error("Length of a deepTree")
	}

	if Length(interestingTree) != list.Length(Flatten(interestingTree)) {
		t.Error("Length of a Tree is equal to length of a flattened tree")
	}
}

func TestIndexedMap(t *testing.T) {
	result := maybe.WithDefault(
		list.Nil[int](),
		maybe.Map(
			Flatten[int],
			IndexedMap(func(i int, s string) int {
				return i
			}, interestingTree),
		),
	)
	if result != list.Range(0, 10) {
		t.Error("Maps a function with index over the Tree, transforms Tree")
	}
}

func TestTuplesOfDatumAndFlatChildren(t *testing.T) {
	e1 := list.FromSlice([]tuple.Tuple[string, list.List[string]]{
		tuple.Pair("a", list.Cons("b", list.Cons("c", list.Singleton("d")))),
		tuple.Pair("b", list.Nil[string]()),
		tuple.Pair("c", list.Nil[string]()),
		tuple.Pair("d", list.Nil[string]()),
	})

	if TuplesOfDatumAndFlatChildren(multiChildTree) != e1 {
		t.Error("TuplesOfDatumAndFlatChildren multiChildTree")
	}

	e2 := list.FromSlice([]tuple.Tuple[string, list.List[string]]{
		tuple.Pair("a", list.Cons("b", list.Cons("c", list.Singleton("d")))),
		tuple.Pair("b", list.Cons("c", list.Singleton("d"))),
		tuple.Pair("c", list.Singleton("d")),
		tuple.Pair("d", list.Nil[string]()),
	})

	if TuplesOfDatumAndFlatChildren(deepTree) != e2 {
		t.Error("TuplesOfDatumAndFlatChildren deepTree")
	}

	e3 := list.FromSlice([]tuple.Tuple[string, list.List[string]]{
		tuple.Pair("a", list.FromSlice([]string{"b", "e", "k", "c", "f", "g", "d", "h", "i", "j"})),
		tuple.Pair("b", list.FromSlice([]string{"e", "k"})),
		tuple.Pair("e", list.Singleton("k")),
		tuple.Pair("k", list.Nil[string]()),
		tuple.Pair("c", list.FromSlice([]string{"f", "g"})),
		tuple.Pair("f", list.Nil[string]()),
		tuple.Pair("g", list.Nil[string]()),
		tuple.Pair("d", list.FromSlice([]string{"h", "i", "j"})),
		tuple.Pair("h", list.Nil[string]()),
		tuple.Pair("i", list.Nil[string]()),
		tuple.Pair("j", list.Nil[string]()),
	})

	if TuplesOfDatumAndFlatChildren(interestingTree) != e3 {
		t.Error("TuplesOfDatumAndFlatChildren interestingTree")
	}
}
