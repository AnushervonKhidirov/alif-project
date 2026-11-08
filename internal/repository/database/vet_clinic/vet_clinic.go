package vetclinic

import (
	"context"
	"errors"
	"pet/internal/model"
	"pet/pkg/optional"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Filters struct {
	Name    *string
	Address *string

	Limit  *int
	Offset *int
}

type CreateVetClinicArgs struct {
	Name    string
	Address string
	About   *string
}

type UpdateVetClinicArgs struct {
	Name    optional.String
	Address optional.String
	About   optional.String
}

type VetClinicRepository struct {
	conn      *pgx.Conn
	tableName string
}

func NewRepository(c *pgx.Conn) *VetClinicRepository {
	return &VetClinicRepository{conn: c, tableName: "vet_clinics"}
}

func (r VetClinicRepository) GetOne(ctx context.Context, id int) (*model.VetClinic, error) {
	var vetClinic model.VetClinic
	sqlStr := `
		SELECT id, name, about, address, created_at, updated_at
		FROM ` + r.tableName + ` WHERE id = $1;
	`
	err := r.conn.QueryRow(ctx, sqlStr, id).Scan(
		&vetClinic.Id,
		&vetClinic.Name,
		&vetClinic.About,
		&vetClinic.Address,
		&vetClinic.CreatedAt,
		&vetClinic.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &vetClinic, nil
}

func (r VetClinicRepository) GetMany(ctx context.Context, filters Filters) ([]model.VetClinic, error) {
	if filters.Limit != nil && *filters.Limit < 0 {
		return nil, errors.New("query limit must not be negative")
	}
	if filters.Offset != nil && *filters.Offset < 0 {
		return nil, errors.New("query offset must not be negative")
	}

	vetClinics := make([]model.VetClinic, 0)
	cond, namedArgs := buildConditions(filters)
	sqlStr := `
		SELECT id, name, about, address, created_at, updated_at
		FROM ` + r.tableName + ` WHERE
	` + cond + " ORDER BY id ASC"

	if filters.Limit != nil {
		namedArgs["limit"] = *filters.Limit
		sqlStr += " LIMIT @limit"
	}

	if filters.Offset != nil {
		namedArgs["offset"] = *filters.Offset
		sqlStr += " OFFSET @offset"
	}

	sqlStr += ";"

	rows, err := r.conn.Query(ctx, sqlStr, namedArgs)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var vetClinic model.VetClinic

		err := rows.Scan(
			&vetClinic.Id,
			&vetClinic.Name,
			&vetClinic.About,
			&vetClinic.Address,
			&vetClinic.CreatedAt,
			&vetClinic.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		vetClinics = append(vetClinics, vetClinic)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return vetClinics, nil
}

func (r VetClinicRepository) Create(ctx context.Context, createArgs CreateVetClinicArgs) (*model.VetClinic, error) {
	var vetClinic model.VetClinic

	namedArgs := make(pgx.NamedArgs)
	namedArgs["name"] = createArgs.Name
	namedArgs["about"] = createArgs.About
	namedArgs["address"] = createArgs.Address

	sqlStr := `
		INSERT INTO ` + r.tableName + ` (name, about, address)
		VALUES (@name, @about, @address)
		RETURNING id, name, about, address, created_at, updated_at;
	`

	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
		&vetClinic.Id,
		&vetClinic.Name,
		&vetClinic.About,
		&vetClinic.Address,
		&vetClinic.CreatedAt,
		&vetClinic.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &vetClinic, nil
}

func (r VetClinicRepository) Update(ctx context.Context, id int, updateArgs UpdateVetClinicArgs) (*model.VetClinic, error) {
	var vetClinic model.VetClinic

	setFields := make([]string, 0)
	namedArgs := make(pgx.NamedArgs)
	namedArgs["id"] = id

	if updateArgs.Name.IsSet {
		if updateArgs.Name.Value == nil {
			return nil, errors.New("name must not be null")
		}

		setFields = append(setFields, "name = @name")
		namedArgs["name"] = updateArgs.Name.Value
	}

	if updateArgs.About.IsSet {
		setFields = append(setFields, "about = @about")
		namedArgs["about"] = updateArgs.About.Value
	}

	if updateArgs.Address.IsSet {
		if updateArgs.Address.Value == nil {
			return nil, errors.New("address must not be null")
		}

		setFields = append(setFields, "address = @address")
		namedArgs["address"] = updateArgs.Address.Value
	}

	_, isNameOk := namedArgs["name"]
	_, isAboutOk := namedArgs["about"]
	_, isAddressOk := namedArgs["address"]

	if !isNameOk && !isAboutOk && !isAddressOk {
		return nil, errors.New("no fields to update")
	}

	setFields = append(setFields, "updated_at = CURRENT_TIMESTAMP")

	sqlStr := `
		UPDATE ` + r.tableName + ` SET ` + strings.Join(setFields, ", ") + `
		WHERE id = @id
		RETURNING id, name, about, address, created_at, updated_at
	`
	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
		&vetClinic.Id,
		&vetClinic.Name,
		&vetClinic.About,
		&vetClinic.Address,
		&vetClinic.CreatedAt,
		&vetClinic.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &vetClinic, nil
}

func (r VetClinicRepository) Delete(ctx context.Context, id int) error {
	sqlStr := "DELETE FROM " + r.tableName + " WHERE id = $1;"
	_, err := r.conn.Exec(ctx, sqlStr, id)
	return err
}

// helpers
func buildConditions(filters Filters) (string, pgx.NamedArgs) {
	cond := "TRUE"
	namedArgs := make(pgx.NamedArgs)

	if filters.Name != nil {
		namedArgs["name"] = *filters.Name
		cond += " AND name ILIKE %@name%"
	}

	if filters.Address != nil {
		namedArgs["address"] = *filters.Address
		cond += " AND address ILIKE %@address%"
	}

	return cond, namedArgs
}
