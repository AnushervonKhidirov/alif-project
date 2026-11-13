package user

import (
	"context"
	"errors"
	"pet/internal/constant"
	"pet/internal/model"
	"pet/pkg/nullable"
	"strings"

	"github.com/jackc/pgx/v5"
)

type FilterOne struct {
	Id    *int
	Email *string
}

type FilterMany struct {
	FirstName *string
	LastName  *string
	Limit     int
	Offset    int
}

type CreateUserArgs struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     nullable.String
}

type UpdateUserArgs struct {
	FirstName    nullable.String
	LastName     nullable.String
	PasswordHash nullable.String
}

type UserRepository struct {
	conn      *pgx.Conn
	tableName string
}

func NewRepository(c *pgx.Conn) *UserRepository {
	return &UserRepository{conn: c, tableName: constant.TableUser}
}

func (r UserRepository) GetOne(ctx context.Context, filters FilterOne) (*model.User, error) {
	var user model.User

	cond, namedArgs := buildOneConditions(filters)
	sqlStr := `
		SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at
		FROM ` + r.tableName + ` WHERE
	` + cond + ";"

	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) GetMany(ctx context.Context, filters FilterMany) ([]model.User, error) {
	if filters.Limit < 0 {
		return nil, errors.New("query limit must not be negative")
	}
	if filters.Offset < 0 {
		return nil, errors.New("query offset must not be negative")
	}

	users := make([]model.User, 0)

	cond, namedArgs := buildManyConditions(filters)

	sqlStr := `
		SELECT id, email, first_name, last_name, role, created_at, updated_at
		FROM ` + r.tableName + ` WHERE
	` + cond + " ORDER BY id ASC"

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
		var u model.User

		err := rows.Scan(
			&u.Id,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Role,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r UserRepository) Create(ctx context.Context, createArgs CreateUserArgs) (*model.User, error) {
	var user model.User

	namedArgs := make(pgx.NamedArgs)
	namedArgs["email"] = createArgs.Email
	namedArgs["passwordHash"] = createArgs.PasswordHash
	namedArgs["firstName"] = createArgs.FirstName
	namedArgs["lastName"] = createArgs.LastName.Value

	sqlStr := `
		INSERT INTO ` + r.tableName + ` (email, password_hash, first_name, last_name)
		VALUES (@email, @passwordHash, @firstName, @lastName)
		RETURNING id, email, first_name, last_name, role, created_at, updated_at;
	`

	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) Update(ctx context.Context, id int, updateArgs UpdateUserArgs) (*model.User, error) {
	var user model.User

	setFields := make([]string, 0)
	namedArgs := make(pgx.NamedArgs)
	namedArgs["id"] = id

	if updateArgs.FirstName.IsSet {
		if updateArgs.FirstName.Value == nil {
			return nil, errors.New("first name must not be null")
		}

		setFields = append(setFields, "first_name = @firstName")
		namedArgs["firstName"] = updateArgs.FirstName.Value
	}

	if updateArgs.LastName.IsSet {
		setFields = append(setFields, "last_name = @lastName")
		namedArgs["lastName"] = updateArgs.LastName.Value
	}

	if updateArgs.PasswordHash.IsSet {
		setFields = append(setFields, "password_hash = @passwordHash")
		namedArgs["passwordHash"] = updateArgs.PasswordHash.Value
	}

	_, isFirstOk := namedArgs["firstName"]
	_, isLastOk := namedArgs["lastName"]
	_, isPasswordOk := namedArgs["passwordHash"]

	if !isFirstOk && !isLastOk && !isPasswordOk {
		return nil, errors.New("no fields to update")
	}

	setFields = append(setFields, "updated_at = CURRENT_TIMESTAMP")

	sqlStr := `
		UPDATE ` + r.tableName + ` SET ` + strings.Join(setFields, ", ") + `
		WHERE id = @id
		RETURNING id, email, first_name, last_name, role, created_at, updated_at;
	`

	err := r.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
		&user.Id,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) Delete(ctx context.Context, id int) error {
	sqlStr := "DELETE FROM " + r.tableName + " WHERE id = $1;"
	_, err := r.conn.Exec(ctx, sqlStr, id)
	return err
}

func (r UserRepository) GetPassword(ctx context.Context, id int) (string, error) {
	var password string
	sqlStr := "SELECT password FROM " + r.tableName + "WHERE id = $1;"
	err := r.conn.QueryRow(ctx, sqlStr, id).Scan(&password)
	return password, err
}

// helpers
func buildOneConditions(filters FilterOne) (string, pgx.NamedArgs) {
	cond := "TRUE"
	namedArgs := make(pgx.NamedArgs)

	if filters.Id != nil {
		namedArgs["id"] = *filters.Id
		cond += " AND id = @id"
	}

	if filters.Email != nil {
		namedArgs["email"] = *filters.Email
		cond += " AND email = @email"
	}

	return cond, namedArgs
}

func buildManyConditions(filters FilterMany) (string, pgx.NamedArgs) {
	cond := "TRUE"
	namedArgs := make(pgx.NamedArgs)

	if filters.FirstName != nil {
		namedArgs["firstName"] = *filters.FirstName
		cond += " AND first_name = @firstName"
	}

	if filters.LastName != nil {
		namedArgs["lastName"] = *filters.LastName
		cond += " AND last_name = @lastName"
	}

	return cond, namedArgs
}
