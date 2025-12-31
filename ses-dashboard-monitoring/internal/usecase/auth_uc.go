package usecase

import (
	"context"
	"errors"
	"time"

	"ses-monitoring/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo user.Repository
	jwtSecret []byte
}

func NewAuthUsecase(userRepo user.Repository, jwtSecret string) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *user.User `json:"user"`
}

func (uc *AuthUsecase) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	u, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  u.ID,
		"username": u.Username,
		"role":     u.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(uc.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: tokenString,
		User:  u,
	}, nil
}

func (uc *AuthUsecase) CreateUser(ctx context.Context, user *user.User) error {
	return uc.userRepo.Create(ctx, user)
}

func (uc *AuthUsecase) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	return uc.userRepo.GetAll(ctx)
}

func (uc *AuthUsecase) ResetUserPassword(ctx context.Context, userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return uc.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}

func (uc *AuthUsecase) DisableUser(ctx context.Context, userID int) error {
	return uc.userRepo.UpdateStatus(ctx, userID, false)
}

func (uc *AuthUsecase) EnableUser(ctx context.Context, userID int) error {
	return uc.userRepo.UpdateStatus(ctx, userID, true)
}

func (uc *AuthUsecase) DeleteUser(ctx context.Context, userID int) error {
	return uc.userRepo.Delete(ctx, userID)
}

func (uc *AuthUsecase) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.userRepo.UpdatePassword(ctx, userID, string(hashedPassword))
}