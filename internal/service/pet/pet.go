package pet

import (
	"context"
	"pet/internal/dto"
	"pet/internal/model"
	"pet/internal/repository/database"
	petDb "pet/internal/repository/database/pet"
	"pet/pkg/nullable"
)

type PetService struct {
	repo database.PetRepo
}

func NewService(r database.PetRepo) *PetService {
	return &PetService{repo: r}
}

func (s PetService) GetOne(ctx context.Context, id int) (*model.Pet, error) {
	pet, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s PetService) GetMany(ctx context.Context, filters petDb.FilterMany) ([]model.Pet, error) {
	pets, err := s.repo.GetMany(ctx, filters)
	if err != nil {
		return nil, err
	}

	return pets, err
}

func (s PetService) Create(ctx context.Context, data dto.CreatePet) (*model.Pet, error) {
	createArgs := petDb.CreateArgs{
		UserId:      data.UserId,
		Name:        data.Name,
		About:       data.About,
		Gender:      data.Gender,
		Birthday:    data.Birthday,
		MicrochipId: data.MicrochipId,
		PetTypeId:   data.PetTypeId,
		BreedId:     data.BreedId,
	}

	pet, err := s.repo.Create(ctx, createArgs)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s PetService) Update(ctx context.Context, filter petDb.FilterOne, data dto.UpdatePet) (*model.Pet, error) {
	updateArgs := petDb.UpdateArgs{
		Name:        data.Name,
		About:       data.About,
		Gender:      data.Gender,
		Birthday:    data.Birthday,
		MicrochipId: data.MicrochipId,
		PetTypeId:   data.PetTypeId,
		BreedId:     data.BreedId,
	}

	pet, err := s.repo.Update(ctx, filter, updateArgs)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s PetService) Delete(ctx context.Context, filter petDb.FilterOne) error {
	err := s.repo.Delete(ctx, filter)
	return err
}

func (s PetService) MakeLost(ctx context.Context, filter petDb.FilterOne) error {
	lost := true
	updateArgs := petDb.UpdateArgs{
		Lost: nullable.Bool{Value: &lost, IsSet: true},
	}

	_, err := s.repo.Update(ctx, filter, updateArgs)
	return err
}

func (s PetService) UnmakeLost(ctx context.Context, filter petDb.FilterOne) error {
	lost := false
	updateArgs := petDb.UpdateArgs{
		Lost: nullable.Bool{Value: &lost, IsSet: true},
	}

	_, err := s.repo.Update(ctx, filter, updateArgs)
	return err
}
