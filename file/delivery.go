package file

import "context"

type Delivery interface {
	Start(ctx context.Context, dir string) error
}
