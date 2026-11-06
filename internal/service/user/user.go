package user

import (
	"context"
	"pet/internal/dto"
	"pet/internal/model"
)

type UserService struct {
}

func NewService() *UserService {
	return &UserService{}
}

func (s UserService) GetById(ctx *context.Context, id int) (*model.User, error) {
	return nil, nil
}

func (s UserService) GetByEmail(ctx *context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (s UserService) GetMany(ctx *context.Context) ([]*model.User, error) {
	return nil, nil
}

func (s UserService) Create(ctx *context.Context, data *dto.CreateUser) (*model.User, error) {
	return nil, nil
}

func (s UserService) Update(ctx *context.Context, id int, data *dto.UpdateUser) error {
	return nil
}

func (s UserService) ChangePassword(ctx *context.Context, id int, data *dto.ChangePassword) error {
	return nil
}

func (s UserService) Delete(ctx *context.Context, id int) error {
	return nil
}
