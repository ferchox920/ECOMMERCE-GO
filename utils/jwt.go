// utils/auth.go

package utils

import (
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/joho/godotenv"
)

func init() {
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateAuthToken(userID string) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(time.Hour * 24).Unix(), // Token expira en 1 día
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }

    return signedToken, nil
}
