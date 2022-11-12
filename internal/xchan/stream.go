package xchan

import (
	"context"
)

func Stream[T any](ctx context.Context, out chan<- T, factory func(context.Context) (T, error)) error {
	for {
		v, err := factory(ctx)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case out <- v:
		}
	}
}
