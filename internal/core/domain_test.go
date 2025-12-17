package core

import (
	"testing"
)

func TestItem_IsValid(t *testing.T) {
	tests := []struct {
		name string
		item Item
		want bool
	}{
		{
			name: "Valid Item should return true",
			item: Item{
				ID:    "123",
				Title: "Super Chien",
				Url:   "https://example.com/dog",
			},
			want: true,
		},
		{
			name: "Missing ID should return false",
			item: Item{
				ID:    "",
				Title: "Super Chien",
				Url:   "https://example.com/dog",
			},
			want: false,
		},
		{
			name: "Missing Title should return false",
			item: Item{
				ID:    "123",
				Title: "",
				Url:   "https://example.com/dog",
			},
			want: false,
		},
		{
			name: "Missing URL should return false",
			item: Item{
				ID:    "123",
				Title: "Super Chien",
				Url:   "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.IsValid(); got != tt.want {
				t.Errorf("Item.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
