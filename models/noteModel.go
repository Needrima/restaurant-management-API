package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id"`
	Text      string             `json:"text" bson:"text,omitempty"`
	CreatedAt time.Time          `json:"created _at" bson:"created _at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at,omitempty"`
	NoteId    string             `json:"note_id" bson:"note_id,omitempty"`
}
