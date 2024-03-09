package utils

import "fmt"

type ControlledPanic struct {
	Msg string
}

func Exitf(statement string, args ...any) {
	panic(ControlledPanic{
		Msg: fmt.Sprintf(statement, args...),
	})
}

func Err(err any) {
	if err != nil {
		panic(err)
	}
}
