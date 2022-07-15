package nub

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Float interface {
	~float32 | ~float64
}

type Num interface {
	Int | Float
}

type Ord interface {
	Num | ~string
}

//----------------------------------------------------------------
// curry
//----------------------------------------------------------------

func Curry[A any, B any, C any](fn func(A, B) C) func(A) func(B) C {
	return func(x A) func(B) C {
		return func(y B) C {
			return fn(x, y)
		}
	}
}

//----------------------------------------------------------------
// Compose
//----------------------------------------------------------------

func Compose[A, B, C any](fn func(B) C, fn1 func(A) B) func(A) C {
	return func(x A) C {
		return fn(fn1(x))
	}
}

//----------------------------------------------------------------
// Id
//----------------------------------------------------------------

func Id[T any](x T) T {
	return x
}

//----------------------------------------------------------------
// constant
//----------------------------------------------------------------

func Const[A, B any](a A) func(B) A {
	return func(_ B) A {
		return a
	}
}

//----------------------------------------------------------------
// flip
//----------------------------------------------------------------

func Flip[A, B, C any](fn func(A, B) C) func(B, A) C {
	return func(x B, y A) C {
		return fn(y, x)
	}
}

//----------------------------------------------------------------
// apply flipped
//----------------------------------------------------------------

func ApplyFlipped[A, B any](a A, fa func(A) B) B {
	return fa(a)
}

//----------------------------------------------------------------
// compare
//----------------------------------------------------------------

type Order int

const (
	LT Order = iota
	EQ
	GT
)

func Compare[T Ord](x T, y T) Order {
	if x < y {
		return LT
	}
	if x > y {
		return GT
	}
	return EQ
}

//----------------------------------------------------------------
// equal
//----------------------------------------------------------------

func Eq[T comparable](x T, y T) bool {
	return x == y
}

//----------------------------------------------------------------
// negate
//----------------------------------------------------------------

func Negate[T Num](x T) T {
	return -x
}

//----------------------------------------------------------------
// min / max
//----------------------------------------------------------------

func Min[T Num](x T, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T Num](x T, y T) T {
	if x > y {
		return x
	}
	return y
}

//----------------------------------------------------------------
// boolean not
//----------------------------------------------------------------

func Not(x bool) bool {
	return !x
}
