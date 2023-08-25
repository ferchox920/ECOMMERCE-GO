package services

import (
	"context"
	"log"

	"github.com/ferchox920/ecommerce-go/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	DB *mongo.Client
}

func NewUserService(client *mongo.Client) *UserService {
	return &UserService{
		DB: client,
	}
}

func (userService *UserService) CreateUser(user *models.User) error {
	collection := userService.DB.Database("ecommerce-go").Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}
	log.Println("Inserted user:", user)
	return nil
}


