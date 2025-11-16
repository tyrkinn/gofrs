package either

import "github.com/snusEbjoer/gofrs/helpers"

type Result[T any] = Either[T, error]

func NewResult[T any](v T, err error) Result[T] {
	if err != nil {
		return Result[T]{
			right:  err,
			isLeft: false,
		}
	}
	return Result[T]{
		left:   v,
		isLeft: true,
	}
}
func (e Either[L, R]) Unwrap() (L, R) {
	return Match2(
		func(t L) (L, R) { return t, helpers.GetDefault[R]() },
		func(r R) (L, R) { return helpers.GetDefault[L](), r },
	)(e)
}

func Unwrap[T any](r Result[T]) (T, error) {
	return Match2(
		func(t T) (T, error) { return t, nil },
		func(err error) (T, error) { return helpers.GetDefault[T](), err },
	)(r)
}
