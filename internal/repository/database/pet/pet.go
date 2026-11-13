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

const selectFields = "id, user_id, name, about, gender, birthday, microchip_id, pet_type_id, breed_id, lost, created_at, updated_at"

type Filters struct {
	UserId      *int
	Gender      *model.Gender
	MicrochipId *string
	Lost        *bool
	PetTypeId   *int
	BreedId     *int
	Limit       int
	Offset      int
}

type CreateArgs struct {
	UserId      int
	Name        string
	About       *string
	Gender      *model.Gender
	Birthday    *time.Time
	MicrochipId *string
	Lost        bool
	PetTypeId   int
	BreedId     *int
}

type UpdateArgs struct {
	Name        nullable.String
	About       nullable.String
	Gender      nullable.Nullable[model.Gender]
	Birthday    nullable.Time
	MicrochipId nullable.String
	Lost        nullable.Bool
	PetTypeId   nullable.Int
	BreedId     nullable.Int
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
	sqlStr := "SELECT " + selectFields + " FROM " + r.tableName + " WHERE id = $1;"
	err := r.conn.QueryRow(ctx, sqlStr, id).Scan(
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

func (r PetRepository) GetMany(ctx context.Context, filters *Filters) ([]model.Pet, error) {
	if filters.Limit < 0 {
		return nil, errors.New("query limit must not be negative")
	}
	if filters.Offset < 0 {
		return nil, errors.New("query offset must not be negative")
	}

	pets := make([]model.Pet, 0)

	cond, namedArgs := buildManyConditions(filters)

	sqlStr := "SELECT " + selectFields + " FROM " + r.tableName + `
		WHERE` + cond + `
		ORDER BY value ASC
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

		err := rows.Scan(
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

func (r PetRepository) Update(ctx context.Context, id int, updateArgs UpdateArgs) (*model.Pet, error) {
	var pet model.Pet

	setFields := make([]string, 0)
	namedArgs := make(pgx.NamedArgs)
	namedArgs["id"] = id

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

	if updateArgs.Lost.IsSet {
		namedArgs["lost"] = updateArgs.Lost.Value
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
		WHERE id = @id
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

func (r PetRepository) Delete(ctx context.Context, id int) error {
	sqlStr := "DELETE FROM " + r.tableName + " WHERE id = $1;"
	_, err := r.conn.Exec(ctx, sqlStr, id)
	return err
}

func buildManyConditions(filters *Filters) (string, pgx.NamedArgs) {
	cond := "TRUE"
	namedArgs := make(pgx.NamedArgs)

	if filters == nil {
		return cond, namedArgs
	}

	if filters.UserId != nil {
		namedArgs["userId"] = *filters.UserId
		cond += " AND user_id = @userId"
	}

	if filters.Gender != nil {
		namedArgs["gender"] = *filters.Gender
		cond += " AND gender = @gender"
	}

	if filters.MicrochipId != nil {
		namedArgs["microchipId"] = *filters.MicrochipId
		cond += " AND microchip_id ILIKE @microchipId%"
	}

	if filters.Lost != nil {
		namedArgs["lost"] = *filters.Lost
		cond += " AND lost = @lost"
	}

	if filters.PetTypeId != nil {
		namedArgs["petTypeId"] = *filters.PetTypeId
		cond += " AND pet_type_id = @petTypeId"
	}

	if filters.BreedId != nil {
		namedArgs["breedId"] = *filters.BreedId
		cond += " AND breed_id = @breedId"
	}

	return cond, namedArgs
}
