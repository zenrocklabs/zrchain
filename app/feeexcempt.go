package app

import "context"

// FeeExcempt is an interface that defines the methods for a fee excempt.
type FeeExcempt interface {
	IsAllowed(ctx context.Context, typeURL string) (bool, error)
}
