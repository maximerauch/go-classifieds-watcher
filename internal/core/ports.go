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
	SaveID(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}
