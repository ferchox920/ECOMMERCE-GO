// auth.go

package services

import (
    "context"
    "log"

    "github.com/ferchox920/ecommerce-go/models"
    "github.com/ferchox920/ecommerce-go/utils"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    DB *mongo.Client
}

func NewAuthService(client *mongo.Client) *AuthService {
    return &AuthService{
        DB: client,
    }
}

func (authService *AuthService) Login(email, password string) (string, error) {
    collection := authService.DB.Database("ecommerce-go").Collection("users")

    user := &models.User{}
    err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(user)
    if err != nil {
        log.Println("error finding user:", err)
        return "", err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", err
    }

    token, err := utils.GenerateAuthToken(user.ID)
    if err != nil {
        log.Println("error generating auth token:", err)
        return "", err
    }

    return token, nil
}
