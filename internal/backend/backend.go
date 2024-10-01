package backend

import "context"

type Expressions interface {
	Process(ctx context.Context, id, exp string) (*int, error)
}
