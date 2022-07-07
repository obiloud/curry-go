package debug

import (
	"encoding/json"
	"fmt"
	"log"
)

func Debug[T any](msg string, a T) T {
	log.Printf("%s: %s", msg, Stringify(a))
	return a
}

type Writer interface {
	String() string
}

func Stringify(x interface{}) string {
	if writer, ok := x.(Writer); ok {
		return writer.String()
	}
	enc, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", x)
	}
	return string(enc)
}
