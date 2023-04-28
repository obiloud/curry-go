package rtree

import (
	"fmt"

	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/util"
)

// The necessary information needed to reconstruct a MultiwayTree as it is
// navigated with a Zipper. This context includes the datum that was at the
// previous node, a list of children that came before the node, and a list of
// children that came after the node.

type Context[T any] struct {
	Previous T
	Before   list.List[RTree[T]]
	After    list.List[RTree[T]]
}

func (c Context[T]) String() string {
	return fmt.Sprintf("Context %s (%s) (%s)", util.Stringify(c.Previous), c.Before.String(), c.After.String())
}

// A list of Contexts that is contructed as a MultiwayTree is navigated.
// Breadcrumbs are used to retain information about parts of the tree that move out
// of focus. As the tree is navigated, the needed Context is pushed onto the list
// Breadcrumbs, and they are maintained in the reverse order in which they are
// visited

// type Breadcrumbs[T any] list.List[Context[T]]

// A structure to keep track of the current Tree, as well as the Breadcrumbs to
// allow us to continue navigation through the rest of the tree.

// type Zipper[T any] tuple.Tuple[RTree[T], Breadcrumbs[T]]

type Zipper[T any] struct {
	Tree        RTree[T]
	Breadcrumbs list.List[Context[T]]
}

func (z Zipper[T]) String() string {
	return fmt.Sprintf("Zipper (%s) (%s)", z.Tree.String(), z.Breadcrumbs.String())
}

// Separate a list into three groups. This function is unique to MultiwayTree
// needs. In order to navigate to children of any Tree, a way to break the children
// into pieces is needed.
// The pieces are:
//   - before: The list of children that come before the desired child
//   - focus: The desired child Tree
//   - after: The list of children that come after the desired child
// These pieces help create a Context, which assist the Zipper

type split[T any] struct {
	before list.List[RTree[T]]
	focus  RTree[T]
	after  list.List[RTree[T]]
}

func splitOnIndex[T any](n int, xs list.List[RTree[T]]) maybe.Maybe[split[T]] {
	before := list.Take(n, xs)

	focus := list.Head(list.Drop(n, xs))

	after := list.Drop(n+1, xs)

	toSplit := func(f RTree[T]) split[T] {
		return split[T]{
			before: before,
			focus:  f,
			after:  after,
		}
	}

	return maybe.Map(toSplit, focus)
}

// Move up relative to the current Zipper focus. This allows navigation from a
// child to it's parent.

func GoUp[T any](zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	if list.IsEmpty(zipper.Breadcrumbs) {
		return maybe.Nothing[Zipper[T]]()
	}

	zip := func(context Context[T]) Zipper[T] {
		return Zipper[T]{
			Tree: RTree[T]{
				Data:     context.Previous,
				Children: list.Append(context.Before, list.Cons(zipper.Tree, context.After)),
			},
			Breadcrumbs: maybe.WithDefault(list.Nil[Context[T]](), list.Tail(zipper.Breadcrumbs)),
		}
	}
	return maybe.Map(zip, list.Head(zipper.Breadcrumbs))
}

// Move down relative to the current Zipper focus. This allows navigation from
// a parent to it's children.

func GoToChild[T any](n int, zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	changeCtx := func(split split[T]) Zipper[T] {
		return Zipper[T]{
			Tree: split.focus,
			Breadcrumbs: list.Cons(Context[T]{
				Previous: zipper.Tree.Data,
				Before:   split.before,
				After:    split.after,
			}, zipper.Breadcrumbs),
		}
	}

	return maybe.Map(changeCtx, splitOnIndex(n, zipper.Tree.Children))
}

// Move down and as far right as possible relative to the current Zipper focus.
// This allows navigation from a parent to it's last child.

func GoToRightMostChild[T any](zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return GoToChild(list.Length(zipper.Tree.Children)-1, zipper)
}

// Move left relative to the current Zipper focus. This allows navigation from
// a child to it's previous sibling.

func GoLeft[T any](zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	if list.IsEmpty(zipper.Breadcrumbs) {
		return maybe.Nothing[Zipper[T]]()
	}

	zip := func(context Context[T]) maybe.Maybe[Zipper[T]] {
		if list.IsEmpty(context.Before) {
			return maybe.Nothing[Zipper[T]]()
		}
		reversed := list.Reverse(context.Before)

		newCtx := func(t RTree[T], rest list.List[RTree[T]]) Zipper[T] {
			return Zipper[T]{
				Tree: t,
				Breadcrumbs: list.Cons(
					Context[T]{
						Previous: context.Previous,
						Before:   list.Reverse(rest),
						After:    list.Cons(zipper.Tree, context.After),
					},
					maybe.WithDefault(list.Nil[Context[T]](), list.Tail(zipper.Breadcrumbs)),
				),
			}
		}
		return maybe.Map2(newCtx, list.Head(reversed), list.Tail(reversed))
	}
	return maybe.Bind(zip, list.Head(zipper.Breadcrumbs))
}

// Move right relative to the current Zipper focus. This allows navigation from
// a child to it's next sibling.

func GoRight[T any](zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	if list.IsEmpty(zipper.Breadcrumbs) {
		return maybe.Nothing[Zipper[T]]()
	}

	zip := func(context Context[T]) maybe.Maybe[Zipper[T]] {
		if list.IsEmpty(context.After) {
			return maybe.Nothing[Zipper[T]]()
		}

		newCtx := func(t RTree[T], rest list.List[RTree[T]]) Zipper[T] {
			return Zipper[T]{
				Tree: t,
				Breadcrumbs: list.Cons(
					Context[T]{
						Previous: context.Previous,
						Before:   list.Append(context.Before, list.Singleton(zipper.Tree)),
						After:    rest,
					},
					maybe.WithDefault(list.Nil[Context[T]](), list.Tail(zipper.Breadcrumbs)),
				),
			}
		}
		return maybe.Map2(newCtx, list.Head(context.After), list.Tail(context.After))
	}
	return maybe.Bind(zip, list.Head(zipper.Breadcrumbs))
}

// Moves to the previous node in the hierarchy, depth-first.

func GoToPrevious[T any](zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return maybe.WithDefault(GoUp(zipper),
		maybe.Map(recurseDownAndRight[T], GoLeft(zipper)),
	)
}

func recurseDownAndRight[T any](z Zipper[T]) maybe.Maybe[Zipper[T]] {
	return maybe.WithDefault(
		maybe.Just(z),
		maybe.Map(recurseDownAndRight[T], GoToRightMostChild(z)),
	)
}

// Moves to the next node in the hierarchy, depth-first. If already
// at the end, stays there.

func GoToNext[T any](zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	first := GoToChild(0, zipper)
	if first.IsJust() {
		return first
	}

	right := GoRight(zipper)
	if right.IsJust() {
		return right
	}

	return upAndOver(zipper)
}

func upAndOver[T any](z Zipper[T]) maybe.Maybe[Zipper[T]] {
	up := GoUp(z)
	if up.IsNothing() {
		return up
	}

	right := maybe.Bind(GoRight[T], up)
	if right.IsJust() {
		return right
	}

	return maybe.Bind(upAndOver[T], up)
}

// Move to the root of the current Zipper focus. This allows navigation from
// any part of the tree back to the root.

func GoToRoot[T any](zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	if list.IsEmpty(zipper.Breadcrumbs) {
		return maybe.Just(zipper)
	}
	return maybe.Bind(GoToRoot[T], GoUp(zipper))
}

// Move the focus to the first element for which the predicate is True. If no
// such element exists returns Nothing. Starts searching at the root of the tree.

func GoTo[T any](predicate func(T) bool, zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return maybe.Bind(nub.Curry(goToElementOrNext[T])(predicate), GoToRoot(zipper))
}

func goToElementOrNext[T any](predicate func(T) bool, z Zipper[T]) maybe.Maybe[Zipper[T]] {
	if predicate(z.Tree.Data) {
		return maybe.Just(z)
	} else {
		return maybe.Bind(nub.Curry(goToElementOrNext[T])(predicate), GoToNext(z))
	}
}

// Update the datum at the current Zipper focus. This allows changes to be made
// to a part of a node's datum information, given the previous state of the node.

func UpdateDatum[T any](fn func(T) T, zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return maybe.Just(Zipper[T]{
		Tree: RTree[T]{
			Data:     fn(zipper.Tree.Data),
			Children: zipper.Tree.Children,
		},
		Breadcrumbs: zipper.Breadcrumbs,
	})
}

// Replace the datum at the current Zipper focus. This allows complete
// replacement of a node's datum information, ignoring the previous state of the
// node.

func ReplaceDatum[T any](newDatum T, zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return UpdateDatum(nub.Const[T, T](newDatum), zipper)
}

func UpdateChildren[T any](newChildren list.List[RTree[T]], zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return maybe.Just(Zipper[T]{
		Tree: RTree[T]{
			Data:     zipper.Tree.Data,
			Children: newChildren,
		},
		Breadcrumbs: zipper.Breadcrumbs,
	})
}

// Inserts a Tree as the first child of the Tree at the current focus. Does not move the focus.

func InsertChildTree[T any](child RTree[T], zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return maybe.Just(Zipper[T]{
		Tree:        InsertChild(child, zipper.Tree),
		Breadcrumbs: zipper.Breadcrumbs,
	})
}

// Inserts a Tree as the last child of the Tree at the current focus. Does not move the focus.

func AppendChildTree[T any](child RTree[T], zipper Zipper[T]) maybe.Maybe[Zipper[T]] {
	return maybe.Just(Zipper[T]{
		Tree:        AppendChild(child, zipper.Tree),
		Breadcrumbs: zipper.Breadcrumbs,
	})
}

// Access the datum at the current Zipper focus.

func Datum[T any](zipper Zipper[T]) T {
	return zipper.Tree.Data
}
