package models

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	User        primitive.ObjectID `bson:"user"`
	SessionName string             `bson:"sessionName"`
	UUID        string             `bson:"uuid"`
	ExpireDate  time.Time          `bson:"expireDate"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`

	UserData *User `bson:"userData,omitempty" json:"userData,omitempty"`
}

// Hooks for Qmgo
func (u *UserSession) BeforeInsert(ctx context.Context) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	// Validate
	validate := validator.New()
	return validate.Struct(u)
}

func (u *UserSession) BeforeUpdate(ctx context.Context) error {
	u.UpdatedAt = time.Now()
	return nil
}
