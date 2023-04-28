package debug

import (
	"log"

	"github.com/obiloud/curry-go/util"
)

func Debug[T any](msg string, a T) T {
	log.Printf("%s: %s", msg, util.Stringify(a))
	return a
}
