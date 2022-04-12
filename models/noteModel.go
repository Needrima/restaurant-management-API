package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Note struct {
	ID        primitive.ObjectID `bson:"_id"`
	Text      string             `json:"text"`
	CreatedAt time.Time          `json:"created _at"`
	UpdatedAt time.Time          `json:"updated_at"`
	NoteId    string             `json:"note_id"`
}
