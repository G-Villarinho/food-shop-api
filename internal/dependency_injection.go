package internal

import "github.com/samber/do"

type Di struct {
	injector *do.Injector
}

func NewDi() *Di {
	return &Di{
		injector: do.New(),
	}
}

func Provide[T any](d *Di, fn func(d *Di) (T, error)) {
	do.Provide(d.injector, func(i *do.Injector) (T, error) {
		return fn(d)
	})
}

func Invoke[T any](d *Di) (T, error) {
	return do.Invoke[T](d.injector)
}
