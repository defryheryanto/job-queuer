package queue

import "context"

type Task interface {
	GetTitle() string
	Do(ctx context.Context) error
}
