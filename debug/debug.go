package debug

import (
	"log"
	"os"

	"github.com/obiloud/curry-go/util"
)

func Debug[T any](msg string, a T) T {
	if os.Getenv("GLOBAL_DEBUG_ENABLED") != "" {
		log.Printf("%s: %s", msg, util.Stringify(a))
	}
	return a
}
