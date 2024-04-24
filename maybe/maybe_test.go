package maybe

import "testing"

func TestWithDefault(t *testing.T) {
	if WithDefault(5, Just(0)) != 0 {
		t.Error("no default used")
	}

	if WithDefault(5, Nothing[int]()) != 5 {
		t.Error("default used")
	}
}

func TestMap(t *testing.T) {
	add1 := func(x int) int {
		return x + 1
	}

	if Map(add1, Just(0)) != Just(1) {
		t.Error("on Just")
	}

	if Map(add1, Nothing[int]()) != Nothing[int]() {
		t.Error("on Nothing")
	}
}

func TestMap2(t *testing.T) {
	sum := func(x int, y int) int {
		return x + y
	}

	if Map2(sum, Just(0), Just(1)) != Just(1) {
		t.Error("on (Just, Just)")
	}

	if Map2(sum, Just(0), Nothing[int]()) != Nothing[int]() {
		t.Error("on (Just, Nothing)")
	}
	if Map2(sum, Nothing[int](), Just(0)) != Nothing[int]() {
		t.Error("on (Nothing, Just)")
	}
}

func TestApply(t *testing.T) {
	add1 := func(x int) int {
		return x + 1
	}

	if Apply[int, int](Just(add1), Just(0)) != Just(1) {
		t.Error("on (Just, Just)")
	}

	if Apply[int, int](Just(add1), Nothing[int]()) != Nothing[int]() {
		t.Error("on (Just, Nothing)")
	}

	if Apply[int, int](Nothing[func(int) int](), Just(0)) != Nothing[int]() {
		t.Error("on (Nothing, Just)")
	}
}

func TestBind(t *testing.T) {
	succeed := func(x int) Maybe[int] {
		return Just(x)
	}

	fail := func(_ int) Maybe[int] {
		return Nothing[int]()
	}

	if Bind(succeed, Just(1)) != Just(1) {
		t.Error("succeeding chain")
	}

	if Bind(succeed, Nothing[int]()) != Nothing[int]() {
		t.Error("original maybe failed")
	}

	if Bind(fail, Just(1)) != Nothing[int]() {
		t.Error("chained function failed")
	}
}
