package main

import (
	"context"
	"errors"
	"example/helpers"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	DB *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return UserRepository{DB: db}
}

func (repository UserRepository) CreateUser(username, password string) (User, error) {
	newUser := User{Username: username, Password: password}
	result, err := DB.Collection("users").InsertOne(context.Background(), newUser)
	if err != nil {
		return newUser, errors.New("failed to insert user")
	}

	newUser.ID = result.InsertedID.(primitive.ObjectID)
	return newUser, nil
}

func (repository UserRepository) GetAllUsers() ([]User, error) {
	result := []User{}
	cursor, err := DB.Collection("users").Find(context.Background(), bson.M{})
	if err != nil {
		return result, errors.New("failed to get users")
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		user := User{}
		if err := cursor.Decode(&user); err != nil {
			return result, errors.New("failed to get users")
		}

		result = append(result, user)
	}

	return result, nil
}

func (repository UserRepository) GetUserByUsername(username string) (User, error) {
	user := User{}
	filter := bson.M{"username": username}
	err := DB.Collection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		fmt.Println("<---- 2", err.Error())
		return user, errors.New("user not found")
	}

	return user, nil
}

func (repository UserRepository) VerifyUserByToken(token string) (User, error) {
	defaultError := errors.New("token not verified")
	user := User{}

	if token == "" {
		return user, defaultError
	}

	claims, err := helpers.VeriyToken(token)
	if err != nil {
		return user, defaultError
	}

	uid := claims["id"].(string)
	objectId, _ := primitive.ObjectIDFromHex(uid)
	err = DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objectId}).
		Decode(&user)

	if err != nil {
		return user, defaultError
	}

	return user, nil
}
