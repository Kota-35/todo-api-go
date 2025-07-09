package service

import (
	"context"
	"time"
	"todo-api-go/pkg/database"
	"todo-api-go/prisma/db"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type UserService struct {
	jwtSecret string
}

func NewUserService(jwtSecret string) *UserService {
	return &UserService{
		jwtSecret: jwtSecret,
	}
}

func (u *UserService) generateToken(userId string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (u *UserService) Register(email string, username string, password string, firstName string, lastName string) (*AuthResponse, error) {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := database.PrismaClient.User.CreateOne(
		db.User.Email.Set(email),
		db.User.Username.Set(username),
		db.User.PasswordHash.Set(string(passwordHash)),
		db.User.FirstName.Set(firstName),
		db.User.LastName.Set(lastName),
	).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	token, expiresAt, err := u.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}

func (u *UserService) Login(email string, password string) (*AuthResponse, error) {

	user, err := database.PrismaClient.User.FindUnique(db.User.Email.Equals(email)).Exec(context.Background())

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, err
	}

	token, expiresAt, err := u.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil

}
