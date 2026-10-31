package service

import (
	"context"
	"pet/internal/dto"
	"pet/internal/model"
)

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (s *userService) GetById(ctx *context.Context, id int) (*model.User, error) {
	return nil, nil
}

func (s *userService) GetByEmail(ctx *context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (s *userService) GetMany(ctx *context.Context) ([]*model.User, error) {
	return nil, nil
}

func (s *userService) Create(ctx *context.Context, data *dto.CreateUser) (*model.User, error) {
	return nil, nil
}

func (s *userService) Update(ctx *context.Context, id int, data *dto.UpdateUser) error {
	return nil
}

func (s *userService) ChangePassword(ctx *context.Context, id int, data *dto.ChangePassword) error {
	return nil
}

func (s *userService) Delete(ctx *context.Context, id int) error {
	return nil
}
