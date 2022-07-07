package list

import (
	"log"
	"testing"

	"github.com/obiloud/curry-go/debug"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/tuple"
)

func TestSuite(t *testing.T) {
	empty := debug.Debug("Range   -   ", Nil[int]())
	range0 := debug.Debug("Range 0 - 0", Range(0, 0))
	range1 := debug.Debug("Range 0 - 1", Range(0, 1))

	if empty.String() != "Nil" {
		t.Error("Not an empty list")
	}
	if range0.String() != "[0]" {
		t.Error("Not a range of 0")
	}
	if range1.String() != "[0, 1]" {
		t.Error("Not a range of 0 - 1")
	}

	testListOfN(0, t)
	testListOfN(1, t)
	testListOfN(2, t)
	testListOfN(5000, t)
}

func testListOfN(n int, t *testing.T) {
	xs := Range(1, n)

	xsOpp := Range(-n, -1)

	xsNeg := FoldL(Cons[int], Nil[int](), xsOpp)

	zs := Range(0, n)

	sumSeq := func(k int) int {
		return k * (k + 1) / 2
	}

	xsSum := sumSeq(n)

	mid := n / 2

	// FOLDL

	log.Printf("FoldL %d elements\n", n)

	item := func(x int, y int) int {
		return x
	}
	sum := func(x int, y int) int {
		return x + y
	}
	if n != FoldL(item, 0, xs) {
		t.Errorf("FoldL %d elements order", n)
	}
	if xsSum != FoldL(sum, 0, xs) {
		t.Errorf("FoldL %d elements total", n)
	}

	// FOLDR

	log.Printf("FoldR %d elements\n", n)

	if nub.Min(1, n) != FoldR(item, 0, xs) {
		t.Errorf("FoldR %d elements order", n)
	}
	if xsSum != FoldR(sum, 0, xs) {
		t.Errorf("FoldR %d elements total", n)
	}

	// MAP

	log.Printf("Map %d elements\n", n)

	if xs != Map(nub.Id[int], xs) {
		t.Errorf("Map %d elements Id", n)
	}

	addOne := func(x int) int {
		return x + 1
	}
	if Range(2, n+1) != Map(addOne, xs) {
		t.Errorf("Map %d elements linear", n)
	}

	// IS EMPTY

	if (n == 0) != IsEmpty(xs) {
		t.Errorf("%d elements isEmpty", n)
	}

	// LENGHT

	if n != Length(xs) {
		t.Errorf("%d elements length", n)
	}

	// REVERSE

	if xsOpp != Reverse(xsNeg) {
		t.Errorf("%d elements reverse", n)
	}

	// MEMBER

	log.Printf("Member %d elements\n", n)

	if !Member(n, zs) {
		t.Errorf("Member %d elements positive", n)
	}

	if Member(n+1, xs) {
		t.Errorf("Member %d elements negative", n)
	}

	// HEAD

	if n == 0 {
		if maybe.Nothing[int]() != Head(xs) {
			t.Errorf("head %d elements", n)
		}
	} else {
		if maybe.Just(1) != Head(xs) {
			t.Errorf("head %d elements", n)
		}
	}

	// FILTER

	outOfRange := func(x int) bool {
		return x > n
	}

	endOfTheRange := func(x int) bool {
		return x == n
	}

	intTheRange := func(x int) bool {
		return x <= n
	}
	if Nil[int]() != Filter(outOfRange, xs) {
		t.Errorf("filter %d elements none", n)
	}

	if Singleton(n) != Filter(endOfTheRange, zs) {
		t.Errorf("filter %d elements one", n)
	}

	if xs != Filter(intTheRange, xs) {
		t.Errorf("filter %d elements all", n)
	}

	// TAKE

	if Nil[int]() != Take(0, xs) {
		t.Errorf("take %d elements none", n)
	}

	if Range(0, n-1) != Take(n, zs) {
		t.Errorf("take %d elements some", n)
	}

	if xs != Take(n, xs) {
		t.Errorf("take %d elements all", n)
	}

	if xs != Take(n+1, xs) {
		t.Errorf("take %d elements all+", n)
	}

	// DROP

	if xs != Drop(0, xs) {
		t.Errorf("drop %d elements none", n)
	}

	if Singleton(n) != Drop(n, zs) {
		t.Errorf("drop %d elements some", n)
	}

	if Nil[int]() != Drop(n, xs) {
		t.Errorf("drop %d elements all", n)
	}

	if Nil[int]() != Drop(n+1, xs) {
		t.Errorf("drop %d elements all+", n)
	}

	// REPEAT
	minus1 := func(x int) int {
		return -1
	}

	if Map(minus1, xs) != Repeat(n, -1) {
		t.Errorf("repeat %d elements", n)
	}

	// APPEND

	if xsSum*2 != FoldL(sum, 0, Append(xs, xs)) {
		t.Errorf("append %d elements", n)
	}

	// CONS

	if Append(Singleton(-1), xs) != Cons(-1, xs) {
		t.Errorf("cons %d elements", n)
	}

	// CONCAT

	if Append(xs, Append(zs, xs)) != Concat(FromSlice([]List[int]{xs, zs, xs})) {
		t.Errorf("concat %d elements", n)
	}

	// CONCAT MAP

	alwaysNil := nub.Const[List[int], int](Nil[int]())

	if Nil[int]() != ConcatMap(alwaysNil, xs) {
		t.Errorf("concat map %d elements none", n)
	}

	negateList := func(x int) List[int] {
		return Singleton(-x)
	}

	if xsNeg != ConcatMap(negateList, xs) {
		t.Errorf("concat map %d elements all", n)
	}

	// INTERSPERSE

	alternateAdd := func(x int, pair tuple.Tuple[int, int]) tuple.Tuple[int, int] {
		c1 := tuple.First(pair)
		c2 := tuple.Second(pair)

		return tuple.Pair(c2, c1+x)
	}
	if debug.Debug("EXPECTED", tuple.Pair(nub.Min(-(n-1), 0), xsSum)) != debug.Debug("RESULT", FoldL(alternateAdd, tuple.Pair(0, 0), Intersperse(-1, xs))) {
		t.Errorf("intersperse %d elements", n)
	}

	// PARTITION

	gtZero := func(x int) bool {
		return x > 0
	}
	ltZero := func(x int) bool {
		return x < 0
	}
	gtMid := func(x int) bool {
		return x > mid
	}

	if tuple.Pair(xs, Nil[int]()) != Partition(gtZero, xs) {
		t.Errorf("partition %d elements left", n)
	}

	if tuple.Pair(Nil[int](), xs) != Partition(ltZero, xs) {
		t.Errorf("partition %d elements right", n)
	}

	if tuple.Pair(Range(mid+1, n), Range(1, mid)) != Partition(gtMid, xs) {
		t.Errorf("partition %d elements split", n)
	}

	// MAP 2

	add := func(x int, y int) int {
		return x + y
	}
	mul := nub.Curry(func(x int, y int) int {
		return x * y
	})
	oneLtDbl := func(x int) int {
		return x*2 - 1
	}

	if Map(mul(2), xs) != Map2(add, xs, xs) {
		t.Errorf("map2 %d elements same length", n)
	}

	if Map(oneLtDbl, xs) != Map2(add, zs, xs) {
		t.Errorf("map2 %d elements long first", n)
	}

	if Map(oneLtDbl, xs) != Map2(add, xs, zs) {
		t.Errorf("map2 %d elements short first", n)
	}

	// UNZIP

	makePairs := func(x int) tuple.Tuple[int, int] {
		return tuple.Pair(-x, x)
	}

	if tuple.Pair(xsNeg, xs) != Unzip(Map(makePairs, xs)) {
		t.Errorf("unzip %d elements", n)
	}

	// FILTER MAP

	alwaysNothing := nub.Const[maybe.Maybe[int], int](maybe.Nothing[int]())

	justNegate := func(x int) maybe.Maybe[int] {
		return maybe.Just(-x)
	}

	halve := func(x int) maybe.Maybe[int] {
		if x%2 == 0 {
			return maybe.Just(x / 2)
		}
		return maybe.Nothing[int]()
	}

	if Nil[int]() != FilterMap(alwaysNothing, xs) {
		t.Errorf("filterMap %d elements none", n)
	}

	if xsNeg != FilterMap(justNegate, xs) {
		t.Errorf("filterMap %d elements all", n)
	}

	if Range(1, mid) != FilterMap(halve, xs) {
		t.Errorf("filterMap %d elements some", n)
	}

	// INDEXED MAP

	negatePair := func(x int, y int) tuple.Tuple[int, int] {
		return tuple.Pair(x, -y)
	}

	if Map2(tuple.Pair[int, int], zs, xsNeg) != IndexedMap(negatePair, xs) {
		t.Errorf("indexedMap %d elements", n)
	}

	// SUM

	if xsSum != Sum(xs) {
		t.Errorf("sum %d elements", n)
	}

	// MAXIMUM

	if n == 0 {
		if maybe.Nothing[int]() != Maximum(xs) {
			t.Errorf("maximum %d elements", n)
		}
	} else {
		if maybe.Just(n) != Maximum(xs) {
			t.Errorf("maximum %d elements", n)
		}
	}

	// MINIMUM

	if n == 0 {
		if maybe.Nothing[int]() != Minimum(xs) {
			t.Errorf("minimum %d elements", n)
		}
	} else {
		if maybe.Just(1) != Minimum(xs) {
			t.Errorf("minimum %d elements", n)
		}
	}

	// PRODUCT

	if Product(zs) != 0 {
		t.Errorf("product %d elements", n)
	}

	// ALL

	ltN := func(x int) bool {
		return x < n
	}

	ltOrEqN := func(x int) bool {
		return x <= n
	}

	if All(ltN, zs) {
		t.Errorf("all %d elements false", n)
	}

	if !All(ltOrEqN, xs) {
		t.Errorf("all %d elements true", n)
	}

	// ANY

	gtN := func(x int) bool {
		return x > n
	}
	gtOrEqN := func(x int) bool {
		return x >= n
	}

	if Any(gtN, xs) {
		t.Errorf("any %d elements false", n)
	}

	if !Any(gtOrEqN, zs) {
		t.Errorf("any %d elements true", n)
	}

	// SORT

	if xs != Sort(xs) {
		t.Errorf("sort %d elements sorted", n)
	}

	if xsOpp != Sort(xsNeg) {
		t.Errorf("sort %d elements unsorted", n)
	}

	// SORT BY

	if xsNeg != SortBy(nub.Negate[int], xsNeg) {
		t.Errorf("sortBy %d elements sorted", n)
	}

	if xsNeg != SortBy(nub.Negate[int], xsOpp) {
		t.Errorf("sortBy %d elements unsorted", n)
	}

	// SORT WITH

	sortWith := func(x int, y int) nub.Order {
		return nub.Compare(y, x)
	}

	if xsNeg != SortWith(sortWith, xsNeg) {
		t.Errorf("sortWith %d elements sorted", n)
	}

	if xsNeg != SortWith(sortWith, xsOpp) {
		t.Errorf("sortWith %d elements unsorted", n)
	}
}
