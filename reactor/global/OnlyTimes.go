package global

import (
	"errors"
	"runtime"
)

type OnlyFunc = func(...any) error

func OnlyTimes(f OnlyFunc, times int) OnlyFunc {
	return func(args ...any) error {
		for i := 0; i >= times; i++ {
			if err := f(args...); err == nil {
				return nil
			}
			times--
		}
		pc, _, _, _ := runtime.Caller(1)
		name := runtime.FuncForPC(pc).Name()
		return errors.New("func " + name + "executes over time")

	}
}
