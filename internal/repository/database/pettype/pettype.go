package pettype

import (
	"context"
	"pet/internal/constant"
	"pet/internal/model"

	"github.com/jackc/pgx/v5"
)

type PetTypeRepository struct {
	conn      *pgx.Conn
	tableName string
}

func NewRepository(c *pgx.Conn) *PetTypeRepository {
	return &PetTypeRepository{conn: c, tableName: constant.TablePetType}
}

func (r PetTypeRepository) GetOne(ctx context.Context, id int) (*model.PetType, error) {
	var petType model.PetType

	sqlStr := "SELECT id, value FROM " + r.tableName + " WHERE id = $1;"
	err := r.conn.QueryRow(ctx, sqlStr, id).Scan(&petType.Id, &petType.Value)
	if err != nil {
		return nil, err
	}

	return &petType, nil
}

func (r PetTypeRepository) GetMany(ctx context.Context) ([]model.PetType, error) {
	petTypes := make([]model.PetType, 0)

	sqlStr := "SELECT id, value FROM " + r.tableName + "  ORDER BY value ASC;"
	rows, err := r.conn.Query(ctx, sqlStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var petType model.PetType

		err := rows.Scan(
			&petType.Id,
			&petType.Value,
		)

		if err != nil {
			return nil, err
		}

		petTypes = append(petTypes, petType)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return petTypes, nil
}

func (r PetTypeRepository) Create(ctx context.Context, value string) (*model.PetType, error) {
	var petType model.PetType

	sqlStr := "INSERT INTO " + r.tableName + " (value) VALUES ($1) RETURNING id, value;"
	err := r.conn.QueryRow(ctx, sqlStr, value).Scan(&petType.Id, &petType.Value)
	if err != nil {
		return nil, err
	}

	return &petType, nil
}

func (r PetTypeRepository) Update(ctx context.Context, petTypeId int, newValue string) (*model.PetType, error) {
	var petType model.PetType
	sqlStr := `
		UPDATE ` + r.tableName + ` SET value = $2
		WHERE id = $1
		RETURNING id, value;
	`
	err := r.conn.QueryRow(ctx, sqlStr, petTypeId, newValue).Scan(&petType.Id, &petType.Value)
	if err != nil {
		return nil, err
	}

	return &petType, nil
}

func (r PetTypeRepository) Delete(ctx context.Context, petTypeId int) error {
	sqlStr := "DELETE FROM " + r.tableName + " WHERE id = $1;"
	_, err := r.conn.Exec(ctx, sqlStr, petTypeId)
	return err
}
