package core

import "context"

type Provider interface {
	Name() string
	FetchListings(ctx context.Context) ([]Listing, error)
}

type Notifier interface {
	Send(ctx context.Context, listing Listing) error
}

type Repository interface {
	Save(ctx context.Context, listing Listing) error
	Exists(ctx context.Context, id string) (bool, error)
}
