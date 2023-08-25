package services

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/ferchox920/ecommerce-go/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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

	// Validar que el Email tenga formato de correo electrónico
	if err := validateEmailFormat(user.Email); err != nil {
		return err
	}

	// Verificar si el email ya está registrado
	existingUser := &models.User{}
	err := collection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(existingUser)
	if err == nil {
		return errors.New("email already exists")
	} else if err != mongo.ErrNoDocuments {
		log.Println("error checking existing email:", err)
		return err
	}

	id := uuid.New().String()
	user.ID = id

	// Hashear la contraseña antes de almacenarla
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error hashing password:", err)
		return err
	}
	user.Password = string(hashedPassword)

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("error creating user:", err)
		return err
	}

	return nil
}

func validateEmailFormat(email string) error {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	match, err := regexp.MatchString(emailPattern, email)
	if err != nil {
		return err
	}
	if !match {
		return errors.New("email must be in a valid format")
	}
	return nil
}