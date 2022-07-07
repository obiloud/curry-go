package either

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	add1 := func(x int) int {
		return x + 1
	}
	if FromRight[string](3) != Map(add1, FromRight[string](2)) {
		t.Error("map OK")
	}

	if FromLeft[string, int]("error") != Map(add1, FromLeft[string, int]("error")) {
		t.Error("map Err")
	}
}

func TestMap2(t *testing.T) {
	sum := func(x, y int) int {
		return x + y
	}
	if FromRight[string](3) != Map2(sum, FromRight[string](1), FromRight[string](2)) {
		t.Error("map2 OK")
	}
	if FromLeft[string, int]("x") != Map2(sum, FromRight[string](1), FromLeft[string, int]("x")) {
		t.Error("map2 Err")
	}
}

func TestApply(t *testing.T) {
	add1 := func(x int) int {
		return x + 1
	}
	if FromRight[string](3) != Apply(FromRight[string](add1), FromRight[string](2)) {
		t.Error("apply OK")
	}
	if FromLeft[string, int]("x") != Apply(FromRight[string](add1), FromLeft[string, int]("x")) {
		t.Error("apply Err")
	}
}

func toInt(x string) Either[string, int] {
	n, err := strconv.Atoi(x)
	if err != nil {
		return FromLeft[string, int](fmt.Sprintf("%v", err))
	}

	return FromRight[string](n)
}

func isEven(n int) Either[string, int] {
	if n%2 == 0 {
		return FromRight[string](n)
	}
	return FromLeft[string, int]("number is odd")
}

func TestBind(t *testing.T) {
	if FromRight[string](42) != Bind(isEven, toInt("42")) {
		t.Error("bind OK")
	}
	if FromLeft[string, int]("strconv.Atoi: parsing \"4.2\": invalid syntax") != Bind(isEven, toInt("4.2")) {
		t.Errorf("bind first")
	}
	if FromLeft[string, int]("number is odd") != Bind(isEven, toInt("41")) {
		t.Error("bind second")
	}
}
