package dict

import (
	"github.com/obiloud/curry-go/list"
	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
	"github.com/obiloud/curry-go/tuple"
)

// A dictionary mapping unique keys to values. The keys can be any comparable type.
type Dict[A nub.Ord, B any] struct {
	dict list.List[tuple.Tuple[A, B]]
}

// Convert a dictionary into a string.
func (dict Dict[A, B]) String() string {
	return dict.dict.String()
}

// Convert an association list into a dictionary.
func FromList[A nub.Ord, B any](ls list.List[tuple.Tuple[A, B]]) Dict[A, B] {
	return Dict[A, B]{dict: list.SortBy(tuple.First[A, B], ls)}
}

// Convert a dictionary into an association list of key-value pairs, sorted by keys.
func ToList[A nub.Ord, B any](dict Dict[A, B]) list.List[tuple.Tuple[A, B]] {
	return dict.dict
}

// Get all of the keys in a dictionary, sorted from lowest to highest.
func Keys[A nub.Ord, B any](dict Dict[A, B]) list.List[A] {
	return list.Map(tuple.First[A, B], dict.dict)
}

// Get all of the values in a dictionary, in the order of their keys.
func Values[A nub.Ord, B any](dict Dict[A, B]) list.List[B] {
	return list.Map(tuple.Second[A, B], dict.dict)
}

// Convert a golang map into a dictionary.
func FromGoMap[A nub.Ord, B any](gomap map[A]B) Dict[A, B] {
	dict := Empty[A, B]()
	for k, v := range gomap {
		dict = Insert(k, v, dict)
	}
	return dict
}

// Convert a dictionary into a golang map.
func ToGoMap[A nub.Ord, B any](dict Dict[A, B]) map[A]B {
	return FoldL(func(key A, value B, acc map[A]B) map[A]B {
		acc[key] = value
		return acc
	}, map[A]B{}, dict)
}

func Empty[A nub.Ord, B any]() Dict[A, B] {
	return Dict[A, B]{dict: list.Nil[tuple.Tuple[A, B]]()}
}

func IsEmpty[A nub.Ord, B any](dict Dict[A, B]) bool {
	return list.IsEmpty(dict.dict)
}

func Singleton[A nub.Ord, B any](key A, value B) Dict[A, B] {
	return Dict[A, B]{dict: list.Singleton(tuple.Pair(key, value))}
}

func Get[A nub.Ord, B any](key A, dict Dict[A, B]) maybe.Maybe[B] {
	return maybe.Map(tuple.Second[A, B],
		list.Head(list.Filter(func(pair tuple.Tuple[A, B]) bool {
			return tuple.First(pair) == key
		}, dict.dict)),
	)
}

func Insert[A nub.Ord, B any](key A, value B, dict Dict[A, B]) Dict[A, B] {
	return FoldR(func(k A, v B, acc Dict[A, B]) Dict[A, B] {
		if !list.Member(k, Keys(acc)) {
			return Dict[A, B]{dict: list.SortBy(tuple.First[A, B], list.Cons(tuple.Pair(k, v), acc.dict))}
		}
		return acc
	}, Singleton(key, value), dict)
}

func Remove[A nub.Ord, B any](key A, dict Dict[A, B]) Dict[A, B] {
	return Dict[A, B]{dict: list.Filter(func(pair tuple.Tuple[A, B]) bool { return tuple.First(pair) != key }, dict.dict)}
}

func Update[A nub.Ord, B any](key A, updateFunc func(maybe.Maybe[B]) maybe.Maybe[B], dict Dict[A, B]) Dict[A, B] {
	return maybe.WithDefault(
		Remove(key, dict),
		maybe.Map(nub.Curry(func(key A, value B) Dict[A, B] {
			return Insert(key, value, dict)
		})(key), updateFunc(Get(key, dict))),
	)
}

func Member[A nub.Ord, B any](key A, dict Dict[A, B]) bool {
	return list.Member(key, Keys(dict))
}

func Size[A nub.Ord, B any](dict Dict[A, B]) int {
	return list.Length(dict.dict)
}

func Union[A nub.Ord, B any](a Dict[A, B], b Dict[A, B]) Dict[A, B] {
	return FoldR(Insert[A, B], a, b)
}

func Intersect[A nub.Ord, B any](a Dict[A, B], b Dict[A, B]) Dict[A, B] {
	return Filter(func(key A, _ B) bool {
		return Member(key, b)
	}, a)
}

func Diff[A nub.Ord, B any](a Dict[A, B], b Dict[A, B]) Dict[A, B] {
	return FoldR(func(key A, _ B, acc Dict[A, B]) Dict[A, B] {
		return Remove(key, acc)
	}, a, b)
}

func Filter[A nub.Ord, B any](isGood func(A, B) bool, dict Dict[A, B]) Dict[A, B] {
	return FoldR(func(key A, value B, acc Dict[A, B]) Dict[A, B] {
		if isGood(key, value) {
			return Insert(key, value, acc)
		}
		return acc
	}, Empty[A, B](), dict)
}

func Partition[A nub.Ord, B any](isGood func(A, B) bool, dict Dict[A, B]) tuple.Tuple[Dict[A, B], Dict[A, B]] {
	return FoldL(func(key A, value B, acc tuple.Tuple[Dict[A, B], Dict[A, B]]) tuple.Tuple[Dict[A, B], Dict[A, B]] {
		insert := func(dict Dict[A, B]) Dict[A, B] {
			return Insert(key, value, dict)
		}
		if isGood(key, value) {
			return tuple.MapFirst(insert, acc)
		}
		return tuple.MapSecond(insert, acc)
	}, tuple.Pair(Empty[A, B](), Empty[A, B]()), dict)
}

func Merge[A nub.Ord, B, C, D any](insertLeft func(A, B, D) D, insertBoth func(A, B, C, D) D, insertRight func(A, C, D) D, left Dict[A, B], right Dict[A, C], result D) D {
	var step func(rKey A, rVal C, acc tuple.Tuple[list.List[tuple.Tuple[A, B]], D]) tuple.Tuple[list.List[tuple.Tuple[A, B]], D]

	step = func(rKey A, rVal C, acc tuple.Tuple[list.List[tuple.Tuple[A, B]], D]) tuple.Tuple[list.List[tuple.Tuple[A, B]], D] {
		xs := tuple.First(acc)

		if list.IsEmpty(xs) {
			return tuple.MapSecond(func(r D) D {
				return insertRight(rKey, rVal, r)
			}, acc)
		}

		return maybe.WithDefault(
			acc,
			maybe.Map2(func(head tuple.Tuple[A, B], tail list.List[tuple.Tuple[A, B]]) tuple.Tuple[list.List[tuple.Tuple[A, B]], D] {
				lKey := tuple.First(head)
				lValue := tuple.Second(head)

				if lKey < rKey {
					return step(rKey, rVal, tuple.Pair(tail, insertLeft(lKey, lValue, tuple.Second(acc))))
				} else if lKey > rKey {
					return tuple.MapSecond(func(r D) D {
						return insertRight(rKey, rVal, r)
					}, acc)
				}

				return tuple.Pair(tail, insertBoth(lKey, lValue, rVal, tuple.Second(acc)))

			}, list.Head(xs), list.Tail(xs)))
	}

	intermediate := FoldL(step, tuple.Pair(left.dict, result), right)

	return list.FoldL(func(pair tuple.Tuple[A, B], acc D) D {
		return insertLeft(tuple.First(pair), tuple.Second(pair), acc)
	}, tuple.Second(intermediate), tuple.First(intermediate))
}

func Map[A nub.Ord, B, C any](fn func(B) C, dict Dict[A, B]) Dict[A, C] {
	return Dict[A, C]{dict: list.Map(func(pair tuple.Tuple[A, B]) tuple.Tuple[A, C] {
		return tuple.MapSecond(fn, pair)
	}, dict.dict)}
}

func FoldL[A nub.Ord, B any, C any](fn func(A, B, C) C, acc C, dict Dict[A, B]) C {
	return list.FoldL(func(pair tuple.Tuple[A, B], acc C) C {
		return fn(tuple.First(pair), tuple.Second(pair), acc)
	}, acc, dict.dict)
}

func FoldR[A nub.Ord, B any, C any](fn func(A, B, C) C, acc C, dict Dict[A, B]) C {
	return list.FoldR(func(pair tuple.Tuple[A, B], acc C) C {
		return fn(tuple.First(pair), tuple.Second(pair), acc)
	}, acc, dict.dict)
}
