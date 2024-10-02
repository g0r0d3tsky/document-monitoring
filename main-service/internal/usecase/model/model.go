package model

import "time"

type Content struct {
	ID        string    `bson:"_id,omitempty"`
	Filename  string    `bson:"filename"`
	Text      string    `bson:"text"`
	CreatedAt time.Time `bson:"created_at"`
}
