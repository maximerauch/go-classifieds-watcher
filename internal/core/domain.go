package core

import "time"

type Listing struct {
	ID          string
	Title       string
	Description string
	Price       float64
	Currency    string
	Url         string
	PublishedAt time.Time
	Metadata    map[string]interface{}
}

func (l Listing) IsValid() bool {
	return l.ID != "" && l.Title != "" && l.Url != ""
}
