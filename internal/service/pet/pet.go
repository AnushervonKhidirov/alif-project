package pet

import (
	"context"
	"pet/internal/dto"
	"pet/internal/dto/query"
	"pet/internal/model"
	"pet/internal/repository/database"
	petDb "pet/internal/repository/database/pet"
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

func (s PetService) GetMany(ctx context.Context, filters query.PetFilterMany) ([]model.Pet, error) {
	dbFilters := petDb.FilterMany{
		UserId:      filters.UserId,
		Gender:      filters.Gender,
		MicrochipId: filters.MicrochipId,
		Lost:        filters.Lost,
		PetTypeId:   filters.PetTypeId,
		BreedId:     filters.BreedId,

		Limit:  filters.Limit,
		Offset: filters.Offset,
	}

	pets, err := s.repo.GetMany(ctx, dbFilters)
	if err != nil {
		return nil, err
	}

	return pets, err
}

func (s PetService) Create(ctx context.Context, userId int, data dto.CreatePet) (*model.Pet, error) {
	createArgs := petDb.CreateArgs{
		UserId:      userId,
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

func (s PetService) Update(ctx context.Context, id int, userId *int, data dto.UpdatePet) (*model.Pet, error) {
	filters := petDb.FilterOne{
		Id:     id,
		UserId: userId,
	}

	updateArgs := petDb.UpdateArgs{
		Name:        data.Name,
		About:       data.About,
		Gender:      data.Gender,
		Birthday:    data.Birthday,
		MicrochipId: data.MicrochipId,
		PetTypeId:   data.PetTypeId,
		BreedId:     data.BreedId,
	}

	pet, err := s.repo.Update(ctx, filters, updateArgs)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (s PetService) Delete(ctx context.Context, id int, userId *int) error {
	filters := petDb.FilterOne{
		Id:     id,
		UserId: userId,
	}

	err := s.repo.Delete(ctx, filters)
	return err
}
