package entity

import (
	"csv-analyzer-api/internal/value"
	"time"
)

type User struct {
	ID             value.UserID `db:"id" bson:"id"`
	Email          string       `db:"email" bson:"email"`
	Name           string       `db:"name" bson:"name"`
	HashedPassword string       `db:"hashed_password" bson:"hashed_password"`
	Role           value.Role   `db:"role" bson:"role"`
	CreatedAt      time.Time    `db:"created_at" bson:"created_at"`
	UpdatedAt      *time.Time   `db:"updated_at" bson:"updated_at"`
}
