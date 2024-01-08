package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfo struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

type LogMessage struct {
	ID         primitive.ObjectID  `bson:"_id"`
	StatusCode int                 `bson:"statusCode"`
	UserInfo   UserInfo            `bson:"userInfo"`
	Timestamp  primitive.Timestamp `bson:"timestamp"`
}
