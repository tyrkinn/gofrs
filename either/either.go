package either

import h "github.com/snusEbjoer/gofrs/helpers"

type Either[L, R any] struct {
	left   L
	right  R
	isLeft bool
}

func Left[L, R any](l L) Either[L, R] {
	return Either[L, R]{
		left:   l,
		isLeft: true,
	}
}

func Right[L, R any](r R) Either[L, R] {
	return Either[L, R]{
		right:  r,
		isLeft: false,
	}
}

func Lift[L, R, T any](f h.F1[T, L]) h.F1[T, Either[L, R]] {
	return func(t T) Either[L, R] {
		return Left[L, R](f(t))
	}
}

func Bind[R, T, T1 any](f h.F1[T, T1]) h.F1[Either[T, R], Either[T1, R]] {
	return Match(
		func(l T) Either[T1, R] { return Left[T1, R](f(l)) },
		Right,
	)
}

func FromPredicate[L, R any](v L, predicate h.Predicate[L], onFalse func(L) R) Either[L, R] {
	if predicate(v) {
		return Either[L, R]{
			left:   v,
			isLeft: true,
		}
	}
	return Either[L, R]{
		right:  onFalse(v),
		isLeft: false,
	}
}

func Match[L, R, V any](f1 func(l L) V, f2 func(r R) V) func(Either[L, R]) V {
	return func(e Either[L, R]) V {
		if e.isLeft {
			return f1(e.left)
		}
		return f2(e.right)
	}
}

func Match2[L, R, V, V1 any](f1 func(l L) (V, V1), f2 func(r R) (V, V1)) func(Either[L, R]) (V, V1) {
	return func(e Either[L, R]) (V, V1) {
		if e.isLeft {
			return f1(e.left)
		}
		return f2(e.right)
	}
}

func Tap[L, R any](f1 func(l L), f2 func(r R)) func(Either[L, R]) Either[L, R] {
	return func(e Either[L, R]) Either[L, R] {
		if e.isLeft {
			f1(e.left)
		} else {
			f2(e.right)
		}
		return e
	}
}

func Map[L, R, V any](mapFn func(l L) V) func(Either[L, R]) Either[V, R] {
	return Match(
		h.Comp2[L, V, Either[V, R]](mapFn, Left),
		Right,
	)
}

func MapRight[L, R, V any](mapFn func(l R) V) func(Either[L, R]) Either[L, V] {
	return Match(
		Left,
		h.Comp2[R, V, Either[L, V]](mapFn, Right),
	)
}

func GetOrElse[L, R any](onLeft func(R) L) func(Either[L, R]) L {
	return Match(h.Identity, onLeft)
}

func FilterOrElse[L, R any](predicate h.Predicate[L], onFalse func(L) R) func(Either[L, R]) Either[L, R] {
	return func(e Either[L, R]) Either[L, R] {
		if predicate(e.left) {
			return e
		}
		return Either[L, R]{
			left:   e.left,
			right:  onFalse(e.left),
			isLeft: false,
		}
	}
}

func FlatMap[L, L1, R any](f h.F1[L, Either[L1, R]]) func(Either[L, R]) Either[L1, R] {
	return Match(
		func(l L) Either[L1, R] { return f(l) },
		func(r R) Either[L1, R] { return Right[L1](r) },
	)
}

func Flatten[L, R any](e Either[Either[L, R], R]) Either[L, R] {
	return Match[Either[L, R]](h.Identity, Right)(e)
}

func Swap[L, R any](e Either[L, R]) Either[R, L] {
	return Either[R, L]{
		left:   e.right,
		right:  e.left,
		isLeft: e.isLeft,
	}
}
