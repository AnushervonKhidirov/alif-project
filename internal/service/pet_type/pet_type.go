package pettype

import (
	"context"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/database"
)

type PetTypeService struct {
	repo database.PetTypeRepo
}

func NewService(r database.PetTypeRepo) *PetTypeService {
	return &PetTypeService{repo: r}
}

func (s PetTypeService) GetOne(ctx context.Context, id int) (*model.PetType, error) {
	petType, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return petType, nil
}

func (s PetTypeService) GetMany(ctx context.Context) ([]model.PetType, error) {
	petTypes, err := s.repo.GetMany(ctx)
	if err != nil {
		return nil, err
	}

	return petTypes, nil
}

func (s PetTypeService) Create(ctx context.Context, data dto.PetType) (*model.PetType, error) {
	petType, err := s.repo.Create(ctx, data.Value)
	if err != nil {
		return nil, err
	}

	return petType, nil
}

func (s PetTypeService) Update(ctx context.Context, id int, data dto.PetType) (*model.PetType, error) {
	petType, err := s.repo.Update(ctx, id, data.Value)
	if err != nil {
		return nil, err
	}

	return petType, nil
}

func (s PetTypeService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
