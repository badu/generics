package generics

import (
	"context"
)

type ctxKey[T any] struct{}

func WithValue[T any](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, ctxKey[T]{}, value)
}

func Value[T any](ctx context.Context) (T, bool) {
	value, ok := ctx.Value(ctxKey[T]{}).(T)
	return value, ok
}
