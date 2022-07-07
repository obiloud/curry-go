package rtree

import (
	"fmt"
	"testing"

	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
)

type record struct {
	Selected bool
	Expanded bool
}

var noChildRecord RTree[record] = RTree[record]{
	Data: record{
		Selected: false,
		Expanded: false,
	},
	Children: list.Nil[RTree[record]](),
}

func TestUpdateDatum(t *testing.T) {
	append := func(x string) string {
		return fmt.Sprintf("%sx", x)
	}
	e1 := maybe.Just(Zipper[string]{
		Tree: RTree[string]{
			Data:     "ax",
			Children: list.Nil[RTree[string]](),
		},
		Breadcrumbs: list.Nil[Context[string]](),
	})
	r1 := maybe.Bind(
		nub.Curry(UpdateDatum[string])(append),
		maybe.Just(Zipper[string]{
			Tree:        noChildTree,
			Breadcrumbs: list.Nil[Context[string]](),
		}),
	)

	if r1 != e1 {
		t.Error("Update datum (simple)")
	}

	e2 := maybe.Just(Zipper[record]{
		Tree: RTree[record]{
			Data: record{
				Selected: true,
				Expanded: false,
			},
			Children: list.Nil[RTree[record]](),
		},
		Breadcrumbs: list.Nil[Context[string]](),
	})

	setSelected := func(r record) record {
		r.Selected = true
		return r
	}

	r2 := maybe.Bind(
		nub.Curry(UpdateDatum[record])(setSelected),
		maybe.Just(Zipper[record]{
			Tree:        noChildRecord,
			Breadcrumbs: list.Nil[Context[string]](),
		}),
	)

	if r2 != e2 {
		t.Error("Update datum (record)")
	}
}

func TestReplaceDatum(t *testing.T) {
	e1 := maybe.Just(Zipper[string]{
		Tree: RTree[string]{
			Data:     "x",
			Children: list.Nil[RTree[string]](),
		},
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r1 := maybe.Bind(
		nub.Curry(ReplaceDatum[string])("x"),
		maybe.Just(Zipper[string]{
			Tree:        noChildTree,
			Breadcrumbs: list.Nil[Context[string]](),
		}),
	)

	if r1 != e1 {
		t.Error("Replace datum")
	}
}

var simpleForest list.List[RTree[string]] = list.Cons(
	RTree[string]{
		Data:     "foo",
		Children: list.Nil[RTree[string]](),
	},
	list.Singleton(RTree[string]{
		Data:     "bar",
		Children: list.Nil[RTree[string]](),
	}),
)

func TestReplaceChildren(t *testing.T) {
	e1 := maybe.Just(Zipper[string]{
		Tree:        noChildTree,
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r1 := maybe.Bind(
		nub.Curry(UpdateChildren[string])(list.Nil[RTree[string]]()),
		maybe.Just(Zipper[string]{
			Tree:        singleChildTree,
			Breadcrumbs: list.Nil[Context[string]](),
		}),
	)

	if r1 != e1 {
		t.Error("Replace children (replace with empty)")
	}

	e2 := maybe.Just(Zipper[string]{
		Tree: RTree[string]{
			Data:     "a",
			Children: simpleForest,
		},
		Breadcrumbs: list.Nil[Context[string]](),
	})

	r2 := maybe.Bind(
		nub.Curry(UpdateChildren[string])(simpleForest),
		maybe.Just(Zipper[string]{
			Tree:        interestingTree,
			Breadcrumbs: list.Nil[Context[string]](),
		}),
	)

	if r2 != e2 {
		t.Error("Replace children (replace with specific)")
	}
}
