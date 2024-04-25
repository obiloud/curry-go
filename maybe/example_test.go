package maybe_test

import (
	"fmt"
	"log"

	"github.com/obiloud/curry-go/maybe"
	"github.com/obiloud/curry-go/nub"
)

func Example_main() {

	// holds no value
	notInt := maybe.Nothing[int]()

	if notInt.IsJust() {
		log.Print("Holds a value")
	} else {
		log.Print("Holds no value")
	}

	// holds the value of 1
	justOne := maybe.Just(1)

	log.Printf("access value on nothing: %d", maybe.WithDefault(0, notInt))
	log.Printf("access value on just: %d", maybe.WithDefault(0, justOne))

	// increment by 1 function
	increment := func(x int) int { return x + 1 }

	// lifts a function over a maybe
	justTwo := maybe.Map(increment, justOne)

	log.Printf("1 + (Just 1) = %s", justTwo)

	// sum two numbers function
	sum := func(x int, y int) int { return x + y }

	// lifts a function over two maybes
	justThree := maybe.Map2(sum, justOne, justTwo)

	log.Printf("(Just 1) + (Just 2) = %s", justThree)

	// function that produces a value if the input is even
	isEven := func(x int) maybe.Maybe[int] {
		if x%2 == 0 {
			return maybe.Just(x)
		}

		return maybe.Nothing[int]()
	}

	// chains computations that may produce value or nothing
	maybeEven := maybe.Bind(isEven, justThree)

	log.Printf("isEven(Just 3) = %s", maybeEven)

	// Maybe is Functor

	log.Printf("Functor Identity %s == %s", maybe.Map(nub.Id[int], justOne), nub.Id(justOne))

	// stringify func
	stringify := func(x int) string { return fmt.Sprint(x) }

	log.Printf("Functor Composition %s == %s", maybe.Map(nub.Compose(stringify, increment), justOne), maybe.Map(stringify, maybe.Map(increment, justOne)))

	// Maybe is Appicative

	log.Printf("Applicative identity %s == %s", maybe.Apply[int, int](maybe.Just(nub.Id[int]), justOne), justOne)

	log.Printf(
		"Applicative Composition %s == %s",
		maybe.Apply[int, string](
			maybe.Just(stringify),
			maybe.Apply[int, int](
				maybe.Just(increment),
				justOne,
			),
		),
		maybe.Apply[int, string](
			maybe.Apply[int, string](
				maybe.Apply[int, string](
					maybe.Just(nub.Curry(nub.Compose[int, int, string])),
					maybe.Just(stringify),
				),
				maybe.Just(increment),
			),
			justOne,
		),
	)
}
