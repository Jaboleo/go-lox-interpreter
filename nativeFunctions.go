package main

import "time"

type ClockFunc struct{}

func (c *ClockFunc) arity() int {
	return 0
}

func (c *ClockFunc) call(interpreter Interpreter, arguments []interface{}) interface{} {
	return time.Now()
}

func (c *ClockFunc) toString() string {
	return "<native fn>"
}
