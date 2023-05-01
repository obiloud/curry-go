package util

import (
	"encoding/json"
	"fmt"
)

type Writer interface {
	String() string
}

func Stringify(x interface{}) string {
	if err, ok := x.(error); ok {
		return fmt.Sprintf("error: %s;", err.Error())
	}
	if writer, ok := x.(Writer); ok {
		return writer.String()
	}
	enc, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", x)
	}
	return string(enc)
}
