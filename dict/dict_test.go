package dict

import (
	"testing"

	"github.com/obiloud/curry-go/debug"
	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/tuple"
)

func TestBuild(t *testing.T) {
	if FromList(list.Nil[tuple.Tuple[string, string]]()) != Empty[string, string]() {
		t.Error("Empty")
	}

	if FromList(list.Singleton(tuple.Pair("k", "v"))) != Singleton("k", "v") {
		t.Error("Singleton")
	}

	if Singleton("k", "v") != Insert("k", "v", Empty[string, string]()) {
		t.Error("Insert")
	}

	if Singleton("k", "b") != Insert("k", "b", Singleton("k", "a")) {
		t.Error("Insert replace")
	}

	updateWithValueFunc := nub.Const[maybe.Maybe[string], maybe.Maybe[string]](maybe.Just("b"))

	if Singleton("k", "b") != Update("k", updateWithValueFunc, Singleton("k", "a")) {
		t.Error("Update")
	}

	updateWithNothingFunc := nub.Const[maybe.Maybe[string], maybe.Maybe[string]](maybe.Nothing[string]())

	if Empty[string, string]() != Update("k", updateWithNothingFunc, Singleton("k", "v")) {
		t.Error("Update Nothing")
	}

	if Empty[string, string]() != Remove("k", Singleton("k", "v")) {
		t.Error("Remove")
	}

	if Singleton("k", "v") != Remove("foo", Singleton("k", "v")) {
		t.Error("Remove not found")
	}

	if FromGoMap(map[string]string{"k": "v"}) != Singleton("k", "v") {
		t.Error("From Go Map")
	}

	if debug.Stringify(ToGoMap(Singleton("k", "v"))) != debug.Stringify(map[string]string{"k": "v"}) {
		t.Error("To Go Map")
	}
}

var animals = FromList(list.Cons(tuple.Pair("Tom", "cat"), list.Singleton(tuple.Pair("Jerry", "mouse"))))

func TestQuery(t *testing.T) {
	if !Member("Tom", animals) {
		t.Error("Member 1")
	}

	if Member("Spike", animals) {
		t.Error("Member 2")
	}

	if Get("Tom", animals) != maybe.Just("cat") {
		t.Error("Get 1")
	}

	if Get("Spike", animals) != maybe.Nothing[string]() {
		t.Error("Get 2")
	}

	if Size(Empty[string, string]()) != 0 {
		t.Error("Size of empty")
	}

	if Size(animals) != 2 {
		t.Error("Size of example")
	}
}

func TestCombine(t *testing.T) {
	if Union(Singleton("Jerry", "mouse"), Singleton("Tom", "cat")) != animals {
		t.Errorf("Union")
	}

	if Union(Singleton("Tom", "cat"), Singleton("Tom", "cat")) != Singleton("Tom", "cat") {
		t.Errorf("Union collision")
	}

	if Intersect(Singleton("Tom", "cat"), animals) != Singleton("Tom", "cat") {
		t.Error("Intersect")
	}

	if Diff(animals, Singleton("Tom", "cat")) != Singleton("Jerry", "mouse") {
		t.Error("Diff")
	}
}

func TestTransform(t *testing.T) {
	if Map(func(x int) int { return x + 1 }, FromList(list.Cons(tuple.Pair("a", 1), list.Singleton(tuple.Pair("b", 2))))) !=
		FromList(list.Cons(tuple.Pair("a", 2), list.Singleton(tuple.Pair("b", 3)))) {
		t.Error("Map")
	}

	if Filter(func(k string, _ string) bool { return k == "Tom" }, animals) != Singleton("Tom", "cat") {
		t.Error("Filter")
	}

	if Partition(func(k string, _ string) bool { return k == "Tom" }, animals) != tuple.Pair(Singleton("Tom", "cat"), Singleton("Jerry", "mouse")) {
		t.Error("Partition")
	}
}

func insertBoth[T nub.Ord](key T, valueLeft list.List[int], valueRight list.List[int], dict Dict[T, list.List[int]]) Dict[T, list.List[int]] {
	return Insert(key, list.Append(valueLeft, valueRight), dict)
}

func TestMerge(t *testing.T) {

	if Merge(
		Insert[int, list.List[int]],
		insertBoth[int],
		Insert[int, list.List[int]],
		Empty[int, list.List[int]](),
		Empty[int, list.List[int]](),
		Empty[int, list.List[int]](),
	) != Empty[int, list.List[int]]() {
		t.Error("Merge empties")
	}

	s1 := Insert("u1", list.Singleton(1), Empty[string, list.List[int]]())

	s2 := Insert("u2", list.Singleton(2), Empty[string, list.List[int]]())

	s23 := Insert("u2", list.Singleton(3), Empty[string, list.List[int]]())

	b1 := FromList(list.Map(func(x int) tuple.Tuple[int, list.List[int]] {
		return tuple.Pair(x, list.Singleton(x))
	}, list.Range(1, 10)))

	b2 := FromList(list.Map(func(x int) tuple.Tuple[int, list.List[int]] {
		return tuple.Pair(x, list.Singleton(x))
	}, list.Range(5, 15)))

	bExpected := list.Map(func(x int) tuple.Tuple[int, list.List[int]] {
		if x > 4 && x < 11 {
			return tuple.Pair(x, list.Cons(x, list.Singleton(x)))
		}
		return tuple.Pair(x, list.Singleton(x))
	}, list.Range(1, 15))

	if ToList(Merge(
		Insert[string, list.List[int]],
		insertBoth[string],
		Insert[string, list.List[int]],
		s1,
		s2,
		Empty[string, list.List[int]](),
	)) != list.FromSlice([]tuple.Tuple[string, list.List[int]]{tuple.Pair("u1", list.Singleton(1)), tuple.Pair("u2", list.Singleton(2))}) {
		t.Error("Merge singletons in order")
	}

	if ToList(Merge(
		Insert[string, list.List[int]],
		insertBoth[string],
		Insert[string, list.List[int]],
		s2,
		s1,
		Empty[string, list.List[int]](),
	)) != list.FromSlice([]tuple.Tuple[string, list.List[int]]{tuple.Pair("u1", list.Singleton(1)), tuple.Pair("u2", list.Singleton(2))}) {
		t.Error("Merge singletons out of order")
	}

	if ToList(Merge(
		Insert[string, list.List[int]],
		insertBoth[string],
		Insert[string, list.List[int]],
		s2,
		s23,
		Empty[string, list.List[int]](),
	)) != list.Singleton(tuple.Pair("u2", list.FromSlice([]int{2, 3}))) {
		t.Error("Merge with duplicate key")
	}

	if ToList(Merge(
		Insert[int, list.List[int]],
		insertBoth[int],
		Insert[int, list.List[int]],
		b1,
		b2,
		Empty[int, list.List[int]](),
	)) != bExpected {
		t.Error("Partially overlapping")
	}
}
