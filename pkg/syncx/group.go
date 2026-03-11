package syncx

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, fns ...func(context.Context) error) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, fn := range fns {
		fn := fn
		g.Go(func() error {
			return fn(ctx)
		})
	}

	return g.Wait()
}
