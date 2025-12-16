package core

import "context"

type Provider interface {
	Name() string
	FetchItems(ctx context.Context) ([]Item, error)
}

type Notifier interface {
	Send(ctx context.Context, item Item) error
}

type Repository interface {
	Save(ctx context.Context, item Item) error
	Exists(ctx context.Context, id string) (bool, error)
}
