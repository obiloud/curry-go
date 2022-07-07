package rtree

import (
	"testing"

	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
)

func TestNavigationSingleChild(t *testing.T) {
	expected := maybe.Just(Zipper[string]{
		Tree: RTree[string]{Data: "b", Children: list.Nil[RTree[string]]()},
		Breadcrumbs: list.Singleton(Context[string]{
			Previous: "a",
			Before:   list.Nil[RTree[string]](),
			After:    list.Nil[RTree[string]](),
		}),
	})
	result := maybe.Bind(
		nub.Curry(GoToChild[string])(0),
		maybe.Just(Zipper[string]{
			Tree:        singleChildTree,
			Breadcrumbs: list.Nil[Context[string]](),
		}),
	)
	if result != expected {
		t.Error("Navigate to child (only child)")
	}
}

func TestNavigationOneOfMany(t *testing.T) {
	expected := maybe.Just(Zipper[string]{
		Tree: RTree[string]{Data: "c", Children: list.Nil[RTree[string]]()},
		Breadcrumbs: list.Singleton(Context[string]{
			Previous: "a",
			Before:   list.Singleton(RTree[string]{Data: "b", Children: list.Nil[RTree[string]]()}),
			After:    list.Singleton(RTree[string]{Data: "d", Children: list.Nil[RTree[string]]()}),
		}),
	})
	result := maybe.Bind(
		nub.Curry(GoToChild[string])(1),
		maybe.Just(Zipper[string]{
			Tree:        multiChildTree,
			Breadcrumbs: list.Nil[Context[string]](),
		}),
	)
	if result != expected {
		t.Error("Navigate to child (one of many)")
	}
}

func TestNavigateChildDeep(t *testing.T) {
	expected := maybe.Just(Zipper[string]{
		Tree: RTree[string]{Data: "d", Children: list.Nil[RTree[string]]()},
		Breadcrumbs: list.Cons(Context[string]{
			Previous: "c",
			Before:   list.Nil[RTree[string]](),
			After:    list.Nil[RTree[string]](),
		}, list.Cons(Context[string]{
			Previous: "b",
			Before:   list.Nil[RTree[string]](),
			After:    list.Nil[RTree[string]](),
		}, list.Cons(Context[string]{
			Previous: "a",
			Before:   list.Nil[RTree[string]](),
			After:    list.Nil[RTree[string]](),
		}, list.Nil[Context[string]]()))),
	})

	in := maybe.Just(Zipper[string]{
		Tree:        deepTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})
	r1 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)
	r2 := maybe.Bind(nub.Curry(GoToChild[string])(0), r1)
	r3 := maybe.Bind(nub.Curry(GoToChild[string])(0), r2)

	if r3 != expected {
		t.Error("Navigate to child (deep)")
	}
}

func TestNavigateToLast(t *testing.T) {
	rempty := maybe.Bind(GoToRightMostChild[string], maybe.Just(Zipper[string]{
		Tree:        noChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	}))

	if rempty != maybe.Nothing[Zipper[string]]() {
		t.Error("Navigate to last child of an empty tree returns Nothing")
	}

	rleft := maybe.Bind(nub.Curry(GoToChild[string])(0), maybe.Just(Zipper[string]{
		Tree:        singleChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	}))

	rright := maybe.Bind(GoToRightMostChild[string], maybe.Just(Zipper[string]{
		Tree:        singleChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	}))

	if rright != rleft {
		t.Error("Navigate to last child of a tree with just one child moves to that child")
	}

	rleft1 := maybe.Bind(nub.Curry(GoToChild[string])(2), maybe.Just(Zipper[string]{
		Tree:        multiChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	}))

	rright1 := maybe.Bind(GoToRightMostChild[string], maybe.Just(Zipper[string]{
		Tree:        multiChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	}))

	if rright1 != rleft1 {
		t.Error("Navigate to last child of a tree with multiple children moves to the last child")
	}

	interesting := maybe.Just(Zipper[string]{
		Tree:        interestingTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	rl1 := maybe.Bind(nub.Curry(GoToChild[string])(2), interesting)
	rl2 := maybe.Bind(nub.Curry(GoToChild[string])(2), rl1)

	rr1 := maybe.Bind(GoToRightMostChild[string], interesting)
	rr2 := maybe.Bind(GoToRightMostChild[string], rr1)

	if rl2 != rr2 {
		t.Error("Navigate to last child of an interestingTree")
	}
}

func TestNavigateUp(t *testing.T) {
	e1 := maybe.Just(Zipper[string]{
		Tree:        RTree[string]{Data: "a", Children: list.Singleton(RTree[string]{Data: "b", Children: list.Nil[RTree[string]]()})},
		Breadcrumbs: list.Nil[Context[string]](),
	})
	r1 := maybe.Bind(GoUp[string],
		maybe.Bind(nub.Curry(GoToChild[string])(0), maybe.Just(Zipper[string]{
			Tree:        singleChildTree,
			Breadcrumbs: list.Nil[Context[string]](),
		})),
	)
	if r1 != e1 {
		t.Error("Navigate up (single level)")
	}

	e2 := maybe.Just(Zipper[string]{
		Tree: RTree[string]{
			Data: "a",
			Children: list.Cons(RTree[string]{
				Data:     "b",
				Children: list.Nil[RTree[string]]()},
				list.Cons(
					RTree[string]{
						Data:     "c",
						Children: list.Nil[RTree[string]]()},
					list.Singleton(RTree[string]{
						Data:     "d",
						Children: list.Nil[RTree[string]]()}),
				),
			),
		},
		Breadcrumbs: list.Nil[Context[string]](),
	})
	r2 := maybe.Bind(GoUp[string],
		maybe.Bind(nub.Curry(GoToChild[string])(1), maybe.Just(Zipper[string]{
			Tree:        multiChildTree,
			Breadcrumbs: list.Nil[Context[string]](),
		})),
	)
	if r2 != e2 {
		t.Error("Navigate up (single level with many children)")
	}

	e3 := maybe.Just(Zipper[string]{
		Tree:        deepTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})
	in := maybe.Just(Zipper[string]{
		Tree:        deepTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})
	r3 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)
	r4 := maybe.Bind(nub.Curry(GoToChild[string])(0), r3)
	r5 := maybe.Bind(nub.Curry(GoToChild[string])(0), r4)
	r6 := maybe.Bind(GoUp[string], r5)
	r7 := maybe.Bind(GoUp[string], r6)
	r8 := maybe.Bind(GoUp[string], r7)

	if r8 != e3 {
		t.Error("Navigate up from a child (deep)")
	}

	r9 := maybe.Just(Zipper[string]{
		Tree:        singleChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})
	r10 := maybe.Bind(nub.Curry(GoToChild[string])(0), r9)
	r11 := maybe.Bind(nub.Curry(GoToChild[string])(0), r10)

	if maybe.Nothing[Zipper[string]]() != r11 {
		t.Error("Navigate beyond the tree (only child)")
	}

	if maybe.Nothing[Zipper[string]]() != maybe.Bind(GoUp[string], r9) {
		t.Error("Navigate beyond the tree (up past root)")
	}
}

func TestNavigateLeft(t *testing.T) {
	if maybe.Nothing[Zipper[string]]() != maybe.Bind(GoLeft[string], maybe.Just(Zipper[string]{
		Tree:        noChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})) {
		t.Error("Navigate to left sibling on no child tree does not work")
	}

	in := maybe.Just(Zipper[string]{
		Tree:        multiChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	left1 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)
	left2 := maybe.Bind(GoLeft[string], left1)

	right1 := maybe.Bind(nub.Curry(GoToChild[string])(2), in)
	right2 := maybe.Bind(GoRight[string], right1)

	if left2 != right2 {
		t.Error("Navigate to left child")
	}

	left3 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)

	right3 := maybe.Bind(nub.Curry(GoToChild[string])(2), in)
	right4 := maybe.Bind(GoLeft[string], right3)
	right5 := maybe.Bind(GoLeft[string], right4)

	if left3 != right5 {
		t.Error("Navigate to left child twice")
	}
}

func TestNavigateRight(t *testing.T) {
	if maybe.Nothing[Zipper[string]]() != maybe.Bind(GoRight[string], maybe.Just(Zipper[string]{
		Tree:        noChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})) {
		t.Error("Navigate to right sibling on no child tree does not work")
	}

	e1 := maybe.Just(Zipper[string]{
		Tree: RTree[string]{
			Data:     "c",
			Children: list.Nil[RTree[string]](),
		},
		Breadcrumbs: list.Singleton(Context[string]{
			Previous: "a",
			Before:   list.Singleton(RTree[string]{Data: "b", Children: list.Nil[RTree[string]]()}),
			After:    list.Singleton(RTree[string]{Data: "d", Children: list.Nil[RTree[string]]()}),
		}),
	})

	in := maybe.Just(Zipper[string]{
		Tree:        multiChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r1 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)
	r2 := maybe.Bind(GoRight[string], r1)

	if r2 != e1 {
		t.Error("Navigate to right child")
	}

	e2 := maybe.Bind(nub.Curry(GoToChild[string])(2), in)

	r3 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)
	r4 := maybe.Bind(GoRight[string], r3)
	r5 := maybe.Bind(GoRight[string], r4)

	if r5 != e2 {
		t.Error("Navigate to right child twice")
	}

	r6 := maybe.Bind(GoRight[string], e2)

	if r6 != maybe.Nothing[Zipper[string]]() {
		t.Error("Navigate to right child when there are no siblings left return Nothing")
	}
}

func TestNavigateNext(t *testing.T) {
	if maybe.Nothing[Zipper[string]]() != maybe.Bind(GoToNext[string], maybe.Just(Zipper[string]{
		Tree:        noChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})) {
		t.Error("Navigate to next child on Tree with just one node")
	}

	in := maybe.Just(Zipper[string]{
		Tree:        interestingTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r1 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)
	e1 := maybe.Bind(nub.Curry(GoToChild[string])(0), r1)

	r2 := maybe.Bind(GoToNext[string], r1)

	if r2 != e1 {
		t.Error("Navigate to next child on an interesting tree will select the next node")
	}

	e2 := maybe.Bind(nub.Curry(GoToChild[string])(1), in)

	r3 := maybe.Bind(nub.Curry(GoToChild[string])(0), in)
	r4 := maybe.Bind(nub.Curry(GoToChild[string])(0), r3)
	r5 := maybe.Bind(nub.Curry(GoToChild[string])(0), r4)
	r6 := maybe.Bind(GoToNext[string], r5)

	if r6 != e2 {
		t.Error("Navigate to next child when the end of a branch has been reached will perform backtracking until the next node down can be reached")
	}

	in2 := maybe.Just(Zipper[string]{
		Tree:        deepTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r7 := maybe.Bind(
		GoToNext[string],
		maybe.Bind(
			GoToNext[string],
			maybe.Bind(
				GoToNext[string],
				maybe.Bind(
					GoToNext[string],
					in2,
				),
			),
		),
	)

	if r7 != maybe.Nothing[Zipper[string]]() {
		t.Error("Navigating past the end of a Tree will return Nothing")
	}

	e3 := maybe.Bind(
		nub.Curry(GoToChild[string])(1),
		maybe.Bind(
			nub.Curry(GoToChild[string])(2),
			in,
		),
	)

	ra1 := maybe.Bind(GoToNext[string], in)
	ra2 := maybe.Bind(GoToNext[string], ra1)
	ra3 := maybe.Bind(GoToNext[string], ra2)
	ra4 := maybe.Bind(GoToNext[string], ra3)
	ra5 := maybe.Bind(GoToNext[string], ra4)
	ra6 := maybe.Bind(GoToNext[string], ra5)
	ra7 := maybe.Bind(GoToNext[string], ra6)
	ra8 := maybe.Bind(GoToNext[string], ra7)
	ra9 := maybe.Bind(GoToNext[string], ra8)

	if ra9 != e3 {
		t.Error("Consecutive goToNext on an interestingTree end up on the right Node")
	}
}

func TestNavigatePrevious(t *testing.T) {
	in := maybe.Just(Zipper[string]{
		Tree:        multiChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r1 := maybe.Bind(nub.Curry(GoToChild[string])(1), in)
	r2 := maybe.Bind(nub.Curry(GoToChild[string])(2), in)
	r3 := maybe.Bind(GoToPrevious[string], r2)

	if r1 != r3 {
		t.Error("Navigate to previous child when there are siblings will select the sibling")
	}

	in1 := maybe.Just(Zipper[string]{
		Tree:        interestingTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r4 := maybe.Bind(nub.Curry(GoToChild[string])(0), in1)

	r5 := maybe.Bind(nub.Curry(GoToChild[string])(0), in1)
	r6 := maybe.Bind(nub.Curry(GoToChild[string])(0), r5)
	r7 := maybe.Bind(GoToPrevious[string], r6)

	if r4 != r7 {
		t.Error("Navigate to previous child on an interesting tree will select the previous node")
	}

	ra1 := maybe.Bind(nub.Curry(GoToChild[string])(0), in1)
	ra2 := maybe.Bind(nub.Curry(GoToChild[string])(0), ra1)
	ra3 := maybe.Bind(nub.Curry(GoToChild[string])(0), ra2)

	ra4 := maybe.Bind(nub.Curry(GoToChild[string])(1), in1)
	ra5 := maybe.Bind(GoToPrevious[string], ra4)

	if ra3 != ra5 {
		t.Error("Navigate to previous child when the beginning of a branch has been reached will perform backtracking until the next node down can be reached")
	}

	in2 := maybe.Just(Zipper[string]{
		Tree:        singleChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	rb1 := maybe.Bind(nub.Curry(GoToChild[string])(0), in2)
	rb2 := maybe.Bind(GoToPrevious[string], rb1)
	rb3 := maybe.Bind(GoToPrevious[string], rb2)

	if rb3 != maybe.Nothing[Zipper[string]]() {
		t.Error("Navigating past the beginning of a Tree will return Nothing")
	}

	e := maybe.Bind(nub.Curry(GoToChild[string])(0), in1)

	rc1 := maybe.Bind(nub.Curry(GoToChild[string])(2), in1)
	rc2 := maybe.Bind(nub.Curry(GoToChild[string])(2), rc1)
	rc3 := maybe.Bind(GoToPrevious[string], rc2)
	rc4 := maybe.Bind(GoToPrevious[string], rc3)
	rc5 := maybe.Bind(GoToPrevious[string], rc4)
	rc6 := maybe.Bind(GoToPrevious[string], rc5)
	rc7 := maybe.Bind(GoToPrevious[string], rc6)
	rc8 := maybe.Bind(GoToPrevious[string], rc7)
	rc9 := maybe.Bind(GoToPrevious[string], rc8)
	rc10 := maybe.Bind(GoToPrevious[string], rc9)
	rc11 := maybe.Bind(GoToPrevious[string], rc10)

	if rc11 != e {
		t.Error("Consecutive goToPrevious on an interestingTree end up on the right Node")
	}
}
