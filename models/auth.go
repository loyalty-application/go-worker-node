package models

import (
	"time"
)

type AuthLoginRequest struct {
	Email    string `json:"email" example:"sudo@soonann.dev" swaggertype:"string"`
	Password string `json:"password" example:"supersecret" swaggertype:"string"`
}

type AuthLoginResponse struct {
	UserID       string    `json:"user_id" example:"6400a..."`
	FirstName    string    `json:"first_name" example:"Soon Ann"`
	LastName     string    `json:"last_name" example:"Tan"`
	Password     string    `json:"password" example:"hashedsupersecret"`
	Email        string    `json:"email" example:"sudo@soonann.dev"`
	Phone        string    `json:"phone" example:"91234567"`
	Token        string    `json:"token" example:"eyJhb..."`
	RefreshToken string    `json:"refresh_token" example:"eyJhb..."`
	CreatedAt    time.Time `json:"created_at" example:"2023-03-02T13:10:23Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-03-02T13:10:23Z"`
}

type AuthRegistrationRequest struct {
	FirstName string `json:"first_name" example:"Soon Ann"`
	LastName  string `json:"last_name" example:"Tan"`
	Password  string `json:"password" example:"supersecret"`
	Email     string `json:"email" example:"sudo@soonann.dev"`
	Phone     string `json:"phone" example:"91234567"`
}

type AuthRegistrationResponse struct {
	UserID string `json:"InsertedID" example:"6400a..."`
}

type UserList struct {
	Users []User `json:"users" bson:",inline"`
}

type User struct {
	UserID       *string   `json:"user_id" bson:"user_id"`
	FirstName    *string   `json:"first_name" bson:"first_name" validate:"required,min=2,max=100"`
	LastName     *string   `json:"last_name" bson:"last_name" validate:"required,min=2,max=100"`
	Password     *string   `json:"password" bson:"password" validate:"required,min=6"`
	Email        *string   `json:"email" bson:"email" validate:"email,required"`
	Phone        *string   `json:"phone" bson:"phone" validate:"required"`
	Token        *string   `json:"token" bson:"token"`
	RefreshToken *string   `json:"refresh_token" bson:"refresh_token"`
	UserType     *string   `json:"user_type" bson:"user_type"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
	Points       float64   `json:"points" bson:"points"`
	Miles        float64   `json:"miles" bson:"miles"`
	Cashback     float64   `json:"cashback" bson:"cashback"`
}

type UserRecordList struct {
	UserRecords []UserRecord `json:"user_records" bson:",inline"`
}

type UserRecord struct {
	Id        string `json:"id" bson:"_id" example:"1"`
	FirstName string `json:"first_name" bson:"first_name" example:"John"`
	LastName  string `json:"last_name" bson:"last_name" example:"Tan"`
	Phone     string `json:"phone" bson:"phone" example:"6591242643"`
	Email     string `json:"email" bson:"email" example:"Dorthea.Ebert@hotmail.com"`
	CreatedAt string `json:"created_at" bson:"created_at" example:"2021-08-23 06:51:25"`
	UpdatedAt string `json:"updated_at" bson:"updated_at" example:"2021-08-23 06:51:25"`
	CardId    string `json:"card_id" bson:"card_id" example:"4111222233334444"`
	CardPan   string `json:"card_pan" bson:"card_pan" example:"xyz"`
	CardType  string `json:"card_type" bson:"card_type" example:"super_miles_card"`
}
