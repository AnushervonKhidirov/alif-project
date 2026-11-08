package user

import (
	"context"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/database"
	breedRepo "pet/internal/repository/database/breed"
)

type BreedService struct {
	repo database.BreedRepo
}

func NewService(r database.BreedRepo) *BreedService {
	return &BreedService{repo: r}
}

func (s BreedService) GetOne(ctx context.Context, id int) (*model.Breed, error) {
	breed, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return breed, nil
}

func (s BreedService) GetMany(ctx context.Context) ([]model.Breed, error) {
	breeds, err := s.repo.GetMany(ctx)
	if err != nil {
		return nil, err
	}

	return breeds, nil
}

func (s BreedService) Create(ctx context.Context, data dto.CreateBreed) (*model.Breed, error) {
	createArgs := breedRepo.CreateBreedArgs{
		Value:     data.Value,
		PetTypeId: data.PetTypeId,
	}

	breed, err := s.repo.Create(ctx, createArgs)
	if err != nil {
		return nil, err
	}

	return breed, nil
}

func (s BreedService) Update(ctx context.Context, id int, data dto.UpdateBreed) (*model.Breed, error) {
	updateArgs := breedRepo.UpdateBreedArgs{
		Value:     data.Value,
		PetTypeId: data.PetTypeId,
	}

	user, err := s.repo.Update(ctx, id, updateArgs)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s BreedService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
