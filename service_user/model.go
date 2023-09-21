package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Username string             `bson:"username,omitempty" json:"full_name"`
	Password string             `bson:"password,omitempty" json:"-"`
}
