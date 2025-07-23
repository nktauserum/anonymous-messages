package article

import "time"

type Article struct {
	UUID          string    `json:"uuid,omitempty"`
	Title         string    `json:"title"`
	Text          string    `json:"text"`
	DatePublished time.Time `json:"date_published,omitempty"`
}
