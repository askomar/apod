package apperror

import (
	"fmt"
	"reflect"
)

// Must is a helper that wraps a call to a function returning a value and an error
// and panics if err is error or false.
// method from samber/lo: https://github.com/samber/lo/blob/407b62d3f12eece919463f556c798661f5aabbbf/errors.go#L64
func Must[T any](val T, err any, messageArgs ...any) T {
	must(err, messageArgs...)
	return val
}

// must panics if err is error or false.
func must(err any, messageArgs ...any) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case bool:
		if !e {
			message := messageFromMsgAndArgs(messageArgs...)
			if message == "" {
				message = "not ok"
			}

			panic(message)
		}

	case error:
		message := messageFromMsgAndArgs(messageArgs...)
		if message != "" {
			panic(message + ": " + e.Error())
		} else {
			panic(e.Error())
		}

	default:
		panic("must: invalid err type '" + reflect.TypeOf(err).Name() + "', should either be a bool or an error")
	}
}

func messageFromMsgAndArgs(msgAndArgs ...any) string {
	if len(msgAndArgs) == 1 {
		if msgAsStr, ok := msgAndArgs[0].(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msgAndArgs[0])
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}
