package breed

import (
	"context"
	"pet/internal/model"

	"github.com/jackc/pgx/v5"
)

type CreateBreedArgs struct {
	Value     string
	PetTypeId int
}

type UpdateBreedArgs struct {
	Value     *string
	PetTypeId *int
}

type BreedRepository struct {
	conn      *pgx.Conn
	tableName string
}

func NewRepository(c *pgx.Conn) *BreedRepository {
	return &BreedRepository{conn: c, tableName: "breeds"}
}

func (r BreedRepository) GetOne(ctx context.Context, id int) (*model.Breed, error) {
	var breed model.Breed
	sqlStr := "SELECT id, value, pet_type_id FROM " + r.tableName + " WHERE id = $1;"
	err := r.conn.QueryRow(ctx, sqlStr, id).Scan(&breed.Id, &breed.Value, &breed.PetTypeId)
	if err != nil {
		return nil, err
	}

	return &breed, nil
}

func (r BreedRepository) GetMany(ctx context.Context) ([]model.Breed, error) {
	breeds := make([]model.Breed, 0)

	sqlStr := "SELECT id, value, pet_type_id FROM " + r.tableName + "  ORDER BY value ASC;"
	rows, err := r.conn.Query(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var breed model.Breed

		err := rows.Scan(&breed.Id, &breed.Value, &breed.PetTypeId)

		if err != nil {
			return nil, err
		}

		breeds = append(breeds, breed)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return breeds, nil
}

func (r BreedRepository) Create(ctx context.Context, breedArgs CreateBreedArgs) (*model.Breed, error) {
	var breed model.Breed
	sqlStr := "INSERT INTO " + r.tableName + " (value, pet_type_id) VALUES ($1, $2) RETURNING id, value, pet_type_id;"
	err := r.conn.QueryRow(ctx, sqlStr, breedArgs.Value, breedArgs.PetTypeId).Scan(&breed.Id, &breed.Value, &breed.PetTypeId)
	if err != nil {
		return nil, err
	}

	return &breed, nil
}

func (r BreedRepository) Update(ctx context.Context, breedId int, updateBreedArgs UpdateBreedArgs) (*model.Breed, error) {
	var breed model.Breed

	namedArgs := make(pgx.NamedArgs)
	namedArgs["id"] = breedId
	namedArgs["value"] = updateBreedArgs.Value
	namedArgs["petTypeId"] = updateBreedArgs.PetTypeId

	sqlStr := `
		UPDATE ` + r.tableName + ` SET
			value = COALESCE(@value, value),
			pet_type_id = COALESCE(@petTypeId, pet_type_id)
		WHERE id = @id
		RETURNING id, value, pet_type_id;
	`
	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(&breed.Id, &breed.Value, &breed.PetTypeId)
	if err != nil {
		return nil, err
	}

	return &breed, nil
}

func (r BreedRepository) Delete(ctx context.Context, petTypeId int) error {
	sqlStr := "DELETE FROM " + r.tableName + " WHERE id = $1;"
	_, err := r.conn.Exec(ctx, sqlStr, petTypeId)
	return err
}
