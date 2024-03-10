package utils

import "fmt"

type ControlledPanic struct {
	Msg      string
	ExitCode int
}

func Exitf(statement string, args ...any) {
	panic(ControlledPanic{
		Msg:      fmt.Sprintf(statement, args...),
		ExitCode: 0,
	})
}

func Failf(statement string, args ...any) {
	panic(ControlledPanic{
		Msg:      fmt.Sprintf(statement, args...),
		ExitCode: 1,
	})
}

func Err(err any) {
	if err != nil {
		panic(err)
	}
}
