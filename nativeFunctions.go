package main

import "time"

type ClockFunc struct{}

func (c *ClockFunc) arity() int {
	return 0
}

func (c *ClockFunc) call(interpreter Interpreter, arguments []any) any {
	return float64(time.Now().UnixMilli())
}

func (c *ClockFunc) toString() string {
	return "<native fn>"
}
