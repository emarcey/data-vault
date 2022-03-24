package common

import (
	"time"

	bsonPrimitive "go.mongodb.org/mongo-driver/bson/primitive"
)

type Responser interface {
	GetStatusCode() int
}

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
	Type       string `json:"type"`
	SecretHash string `json:"-"`
	StatusCode int    `json:"-"`
}

func (u *User) IsAdmin() bool {
	return u.Type == "admin"
}

func (u *User) GetStatusCode() int {
	if u.StatusCode == 0 {
		return 200
	}
	return u.StatusCode
}

type AccessToken struct {
	Id         string    `json:"id"`
	UserId     string    `json:"client_id"`
	InvalidAt  time.Time `json:"invalid_at"`
	IsLatest   bool      `json:"is_latest"`
	StatusCode int       `json:"-"`
}

func (a *AccessToken) GetStatusCode() int {
	if a.StatusCode == 0 {
		return 200
	}
	return a.StatusCode
}

type Secret struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	UpdatedBy   string `json:"updated_by"`
	StatusCode  int    `json:"-"`
}

func (s *Secret) GetStatusCode() int {
	if s.StatusCode == 0 {
		return 200
	}
	return s.StatusCode
}

type EncryptedSecret struct {
	Id  string `json:"_id" bson:"_id"`
	Key string `json:"key" bson:"key"`
	Iv  string `json:"iv" bson:"iv"`
}

type AccessLog struct {
	UserId     string                 `json:"user_id" bson"user_id"`
	ActionType string                 `json:"action_type" bson:"action_type"`
	KeyName    string                 `json:"key_name" bson:"key_name"`
	AccessAt   bsonPrimitive.DateTime `json:"access_at" bson:"access_at"`
}

func NewAccessLog(userId, actionType, keyName string) *AccessLog {
	return &AccessLog{
		UserId:     userId,
		ActionType: actionType,
		KeyName:    keyName,
		AccessAt:   bsonPrimitive.NewDateTimeFromTime(time.Now().UTC()),
	}
}
