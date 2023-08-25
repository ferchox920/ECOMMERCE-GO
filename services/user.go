package services

import (
	"context"
	"errors"
	"github.com/ferchox920/ecommerce-go/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
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

	if err := validateEmailFormat(user.Email); err != nil {
		return err
	}

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

func (userService *UserService) FindAllUsers() ([]models.User, error) {
	collection := userService.DB.Database("ecommerce-go").Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println("error finding users:", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		log.Println("error decoding users:", err)
		return nil, err
	}

	return users, nil
}

func (userService *UserService) FindUserByID(id string) (*models.User, error) {
	collection := userService.DB.Database("ecommerce-go").Collection("users")
	user := &models.User{}
	err := collection.FindOne(context.Background(), bson.M{"id": id}).Decode(user)
	if err != nil {
		log.Println("error finding user:", err)
		return nil, err
	}
	return user, nil
}

func (userService *UserService) UpdateUser(user *models.User) error {
	collection := userService.DB.Database("ecommerce-go").Collection("users")

	// Validar que el Email tenga formato de correo electrónico
	if err := validateEmailFormat(user.Email); err != nil {
		return err
	}

	// Actualizar todas las propiedades excepto email y password
	updateFields := bson.M{
		"name":     user.Name,
		"lastname": user.Lastname,
		"adress":   user.Adress,
		"adress2":  user.Adress2,
		"city":     user.City,
		"state":    user.State,
		"zip":      user.Zip,
		"phone":    user.Phone,
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"id": user.ID}, bson.M{"$set": updateFields})
	if err != nil {
		log.Println("error updating user:", err)
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
