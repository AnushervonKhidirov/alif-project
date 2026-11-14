package user

import (
	"context"
	"errors"
	"pet/internal/dto"
	"pet/internal/dto/query"
	"pet/internal/model"
	"pet/internal/repository/database"
	userRepo "pet/internal/repository/database/user"
	"pet/pkg/nullable"
	"pet/pkg/password"
)

type UserService struct {
	repo database.UserRepo
}

func NewService(r database.UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s UserService) GetOne(ctx context.Context, filters userRepo.FilterOne) (*model.User, error) {
	user, err := s.repo.GetOne(ctx, filters)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s UserService) GetMany(ctx context.Context, filters query.UserFilterMany) ([]model.User, error) {
	dbFilters := userRepo.FilterMany{
		FirstName: filters.FirstName,
		LastName:  filters.LastName,

		Limit:  filters.Limit,
		Offset: filters.Offset,
	}

	users, err := s.repo.GetMany(ctx, dbFilters)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s UserService) Create(ctx context.Context, data dto.CreateUser) (*model.User, error) {
	passwordHash, err := password.Hash(data.Password)
	if err != nil {
		return nil, err
	}

	createArgs := userRepo.CreateUserArgs{
		Email:        data.Email,
		PasswordHash: passwordHash,
		FirstName:    data.FirstName,
		LastName:     data.LastName,
	}

	user, err := s.repo.Create(ctx, createArgs)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s UserService) Update(ctx context.Context, id int, data dto.UpdateUser) (*model.User, error) {
	updateArgs := userRepo.UpdateUserArgs{FirstName: data.FirstName, LastName: data.LastName}
	user, err := s.repo.Update(ctx, id, updateArgs)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s UserService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s UserService) ChangePassword(ctx context.Context, id int, data dto.ChangePassword) error {
	oldPasswordHash, err := s.repo.GetPassword(ctx, id)
	if err != nil {
		return err
	}

	if isCorrect := password.CheckHash(data.OldPassword, oldPasswordHash); !isCorrect {
		return errors.New("wrong old password")
	}

	newPasswordHash, err := password.Hash(data.NewPassword)
	if err != nil {
		return err
	}

	passwordHash := nullable.String{Value: &newPasswordHash, IsSet: true}
	_, err = s.repo.Update(ctx, id, userRepo.UpdateUserArgs{PasswordHash: passwordHash})
	return err
}
