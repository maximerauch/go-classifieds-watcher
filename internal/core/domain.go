package core

import "time"

type Item struct {
	ID          string
	Title       string
	Description string
	Price       float64
	Currency    string
	Url         string
	PublishedAt time.Time
	Metadata    map[string]interface{}
}

func (i Item) IsValid() bool {
	return i.ID != "" && i.Title != "" && i.Url != ""
}
