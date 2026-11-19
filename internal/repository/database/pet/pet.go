package pet

import (
	"context"
	"errors"
	"pet/internal/constant"
	"pet/internal/model"
	"pet/pkg/nullable"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	prefixPlaceholder = "[prefix]"
	selectFields      = "id, user_id, name, about, gender, birthday, microchip_id, pet_type_id, breed_id, lost, created_at, updated_at"

	// NOTE: p => pets, pt => pet_types, b => breeds
	selectFieldsWithJoin = "p.id, p.user_id, p.name, p.about, p.gender, p.birthday, p.microchip_id, p.pet_type_id, p.breed_id, p.lost, p.created_at, p.updated_at, pt.id, pt.value, b.id, b.value, b.pet_type_id"
)

type FilterMany struct {
	UserId      *int
	Gender      *model.Gender
	MicrochipId *string
	Lost        *bool
	PetTypeId   *int
	BreedId     *int

	Limit  int
	Offset int
}

type FilterOne struct {
	Id     int
	UserId *int
}

type CreateArgs struct {
	UserId      int
	Name        string
	About       *string
	Gender      *model.Gender
	Birthday    *time.Time
	MicrochipId *string
	PetTypeId   int
	BreedId     *int
}

type UpdateArgs struct {
	Name        nullable.String
	About       nullable.String
	Gender      nullable.Nullable[model.Gender]
	Birthday    nullable.Time
	MicrochipId nullable.String
	PetTypeId   nullable.Int
	BreedId     nullable.Int
	Lost        nullable.Bool
}

type PetRepository struct {
	conn      *pgx.Conn
	tableName string
}

func NewRepository(c *pgx.Conn) *PetRepository {
	return &PetRepository{conn: c, tableName: constant.TablePet}
}

func (r PetRepository) GetOne(ctx context.Context, id int) (*model.Pet, error) {
	var pet model.Pet
	sqlStr := `
		SELECT ` + selectFieldsWithJoin + `
		FROM ` + r.tableName + ` AS p
		LEFT JOIN ` + constant.TablePetType + ` AS pt
			ON p.pet_type_id = pt.id
		LEFT JOIN ` + constant.TableBreed + ` AS b
			ON p.breed_id = b.id
		WHERE p.id = $1;
	`

	row := r.conn.QueryRow(ctx, sqlStr, id)
	if err := collectPetData(row, &pet); err != nil {
		return nil, err
	}
	return &pet, nil
}

func (r PetRepository) GetMany(ctx context.Context, filters FilterMany) ([]model.Pet, error) {
	pets := make([]model.Pet, 0)
	petQueryPrefix := "p"
	cond, namedArgs := buildManyConditions(filters, &petQueryPrefix)

	sqlStr := `
		SELECT ` + selectFieldsWithJoin + `
		FROM ` + r.tableName + ` AS p
		LEFT JOIN ` + constant.TablePetType + ` AS pt
			ON p.pet_type_id = pt.id
		LEFT JOIN ` + constant.TableBreed + ` AS b
			ON p.breed_id = b.id
		WHERE` + cond + `
		ORDER BY p.name ASC
	`

	namedArgs["limit"] = filters.Limit
	sqlStr += " LIMIT @limit"

	namedArgs["offset"] = filters.Offset
	sqlStr += " OFFSET @offset"

	sqlStr += ";"

	rows, err := r.conn.Query(ctx, sqlStr, namedArgs)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var pet model.Pet
		if err := collectPetData(rows, &pet); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pets, nil
}

func (r PetRepository) Create(ctx context.Context, createArgs CreateArgs) (*model.Pet, error) {
	var pet model.Pet

	namedArgs := make(pgx.NamedArgs)
	namedArgs["user_id"] = createArgs.UserId

	sqlStr := `
		INSERT INTO ` + r.tableName + ` (user_id, name, about, gender, birthday, microchip_id, pet_type_id, breed_id)
		VALUES ()
		RETURNING ` + selectFields + `;
	`

	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
		&pet.Id,
		&pet.UserId,
		&pet.Name,
		&pet.About,
		&pet.Gender,
		&pet.Birthday,
		&pet.MicrochipId,
		&pet.PetTypeId,
		&pet.BreedId,
		&pet.Lost,
		&pet.CreatedAt,
		&pet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (r PetRepository) Update(ctx context.Context, filters FilterOne, updateArgs UpdateArgs) (*model.Pet, error) {
	var pet model.Pet

	setFields := make([]string, 0)
	namedArgs := make(pgx.NamedArgs)
	namedArgs["id"] = filters.Id

	cond := "id = @id"

	if filters.UserId != nil {
		namedArgs["userId"] = filters.UserId
		cond += " AND user_id = @userId"
	}

	if updateArgs.Name.IsSet {
		if updateArgs.Name.Value == nil {
			return nil, errors.New("name should not be null")
		}

		namedArgs["name"] = updateArgs.Name.Value
	}

	if updateArgs.About.IsSet {
		namedArgs["about"] = updateArgs.About.Value
	}

	if updateArgs.Gender.IsSet {
		namedArgs["gender"] = updateArgs.Gender.Value
	}

	if updateArgs.Birthday.IsSet {
		namedArgs["birthday"] = updateArgs.Birthday.Value
	}

	if updateArgs.MicrochipId.IsSet {
		namedArgs["microchipId"] = updateArgs.MicrochipId.Value
	}

	if updateArgs.PetTypeId.IsSet {
		if updateArgs.PetTypeId.Value == nil {
			return nil, errors.New("pet type id should not be null")
		}

		namedArgs["petTypeId"] = updateArgs.PetTypeId.Value
	}

	if updateArgs.BreedId.IsSet {
		namedArgs["breedId"] = updateArgs.BreedId.Value
	}

	setFields = append(setFields, "updated_at = CURRENT_TIMESTAMP")

	sqlStr := `
		UPDATE ` + r.tableName + ` SET ` + strings.Join(setFields, ", ") + `
		WHERE ` + cond + `
		RETURNING` + selectFields + `;
	`
	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
		&pet.Id,
		&pet.UserId,
		&pet.Name,
		&pet.About,
		&pet.Gender,
		&pet.Birthday,
		&pet.MicrochipId,
		&pet.PetTypeId,
		&pet.BreedId,
		&pet.Lost,
		&pet.CreatedAt,
		&pet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func (r PetRepository) Delete(ctx context.Context, filters FilterOne) error {
	namedArgs := make(pgx.NamedArgs)
	namedArgs["id"] = filters.Id
	cond := "id = @id"

	if filters.UserId != nil {
		namedArgs["userId"] = filters.UserId
		cond += " AND user_id = @userId"
	}

	sqlStr := "DELETE FROM " + r.tableName + " WHERE " + cond + ";"
	_, err := r.conn.Exec(ctx, sqlStr, namedArgs)
	return err
}

func buildManyConditions(filters FilterMany, prefix *string) (string, pgx.NamedArgs) {
	cond := "TRUE"
	namedArgs := make(pgx.NamedArgs)

	if filters.UserId != nil {
		namedArgs["userId"] = *filters.UserId
		cond += " AND " + prefixPlaceholder + "user_id = @userId"
	}

	if filters.Gender != nil {
		namedArgs["gender"] = *filters.Gender
		cond += " AND " + prefixPlaceholder + "gender = @gender"
	}

	if filters.MicrochipId != nil {
		namedArgs["microchipId"] = *filters.MicrochipId + "%"
		cond += " AND " + prefixPlaceholder + "microchip_id ILIKE @microchipId"
	}

	if filters.Lost != nil {
		namedArgs["lost"] = *filters.Lost
		cond += " AND " + prefixPlaceholder + "lost = @lost"
	}

	if filters.PetTypeId != nil {
		namedArgs["petTypeId"] = *filters.PetTypeId
		cond += " AND " + prefixPlaceholder + "pet_type_id = @petTypeId"
	}

	if filters.BreedId != nil {
		namedArgs["breedId"] = *filters.BreedId
		cond += " AND " + prefixPlaceholder + "breed_id = @breedId"
	}

	if prefix != nil {
		cond = strings.ReplaceAll(cond, prefixPlaceholder, *prefix+".")
	} else {
		cond = strings.ReplaceAll(cond, prefixPlaceholder, "")
	}

	return cond, namedArgs
}

func collectPetData(row pgx.Row, pet *model.Pet) error {
	var petType model.PetType

	var breedId *int
	var breedValue *string
	var breedPetTypeId *int

	err := row.Scan(
		&pet.Id,
		&pet.UserId,
		&pet.Name,
		&pet.About,
		&pet.Gender,
		&pet.Birthday,
		&pet.MicrochipId,
		&pet.PetTypeId,
		&pet.BreedId,
		&pet.Lost,
		&pet.CreatedAt,
		&pet.UpdatedAt,

		&petType.Id,
		&petType.Value,

		&breedId,
		&breedValue,
		&breedPetTypeId,
	)

	if err != nil {
		return err
	}

	pet.PetType = &petType

	if breedId != nil {
		breed := &model.Breed{
			Id:        *breedId,
			Value:     *breedValue,
			PetTypeId: *breedPetTypeId,
		}

		pet.Breed = breed
	}

	return nil
}
