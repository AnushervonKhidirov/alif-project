package user

import (
	"context"
	"errors"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/database"
	userRepo "pet/internal/repository/database/user"
	"pet/pkg/optional"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo database.UserRepo
}

func NewService(r database.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s UserService) GetOne(ctx context.Context, filter userRepo.Filters) (*model.User, error) {
	user, err := s.repo.GetOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s UserService) GetMany(ctx context.Context, filter userRepo.Filters) ([]model.User, error) {
	users, err := s.repo.GetMany(ctx, filter)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s UserService) Create(ctx context.Context, data *dto.CreateUser) (*model.User, error) {
	passwordHash, err := s.hashPassword(data.Password)
	if err != nil {
		return nil, errors.New("Unable to hash password")
	}

	createArgs := userRepo.CreateUserArgs{
		Email:        data.Email,
		PasswordHash: passwordHash,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
	}

	u, err := s.repo.Create(ctx, createArgs)
	if err != nil {
		return nil, errors.New("Unable to create user")
	}

	return u, nil
}

func (s UserService) Update(ctx context.Context, id int, data *dto.UpdateUser) (*model.User, error) {
	updateArgs := userRepo.UpdateUserArgs{FirstName: data.FirstName, LastName: data.LastName}
	u, err := s.repo.Update(ctx, id, updateArgs)
	if err != nil {
		return nil, errors.New("Unable to update user")
	}

	return u, nil
}

func (s UserService) ChangePassword(ctx context.Context, id int, newPassword string) error {
	newPasswordHash, err := s.hashPassword(newPassword)
	if err != nil {
		return errors.New("Unable to hash password")
	}

	passwordHash := optional.String{Value: &newPasswordHash, IsSet: true}
	_, err = s.repo.Update(ctx, id, userRepo.UpdateUserArgs{PasswordHash: passwordHash})
	return err
}

func (s UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

// helper
func (s UserService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s UserService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
