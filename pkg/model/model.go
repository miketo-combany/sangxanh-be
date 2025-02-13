package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IModel interface {
	SetID(id primitive.ObjectID)
	CollectionName() string
	GetID() primitive.ObjectID
}

type Model struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (m *Model) SetID(id primitive.ObjectID) {
	m.ID = id
}

func (m *Model) GetID() primitive.ObjectID {
	return m.ID
}
