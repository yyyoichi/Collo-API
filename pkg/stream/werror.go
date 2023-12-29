package stream

import "context"

type ErrorHook func(error)
type WithError interface {
	Error() error
}

// inputにエラーがあった場合にerrHookを呼び出す
func LineWithErrorHook[I WithError, O interface{}](
	cxt context.Context,
	errHook ErrorHook,
	inCh <-chan I,
	fn func(I) O,
) <-chan O {
	return Line[I, O](cxt, inCh, func(i I) O {
		if i.Error() != nil {
			errHook(i.Error())
		}
		return fn(i)
	})
}

// inputにエラーがあった場合にerrHookを呼び出す
func FunIOWithErrorHook[I WithError, O interface{}](
	cxt context.Context,
	errHook ErrorHook,
	inCh <-chan I,
	fn func(I) O,
) <-chan O {
	return FunIO[I, O](cxt, inCh, func(i I) O {
		if i.Error() != nil {
			errHook(i.Error())
		}
		return fn(i)
	})
}

// inputにエラーがあった場合にerrHookを呼び出す
func DemultiWithErrorHook[I WithError, O interface{}](
	cxt context.Context,
	errHook ErrorHook,
	inCh <-chan I,
	fn func(I) []O,
) <-chan O {
	return Demulti[I, O](cxt, inCh, func(i I) []O {
		if i.Error() != nil {
			errHook(i.Error())
		}
		return fn(i)
	})
}
