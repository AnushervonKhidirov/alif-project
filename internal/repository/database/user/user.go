package user

import (
	"context"
	"errors"
	"pet/internal/model"
	"pet/pkg/optional"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Filters struct {
	Id        *int
	Email     *string
	FirstName *string
	LastName  *string

	Limit  *int
	Offset *int
}

type CreateUserArgs struct {
	Email        string
	PasswordHash string
	FirstName    string
	LastName     optional.String
}

type UpdateUserArgs struct {
	FirstName    optional.String
	LastName     optional.String
	PasswordHash optional.String
}

type UserRepository struct {
	conn      *pgx.Conn
	tableName string
}

func NewRepository(c *pgx.Conn) *UserRepository {
	return &UserRepository{conn: c, tableName: "users"}
}

func (r UserRepository) GetOne(ctx context.Context, filters Filters) (*model.User, error) {
	var user model.User

	cond, namedArgs := buildConditions(filters)
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

func (r UserRepository) GetMany(ctx context.Context, filters Filters) ([]model.User, error) {
	if filters.Limit != nil && *filters.Limit < 0 {
		return nil, errors.New("query limit must not be negative")
	}
	if filters.Offset != nil && *filters.Offset < 0 {
		return nil, errors.New("query offset must not be negative")
	}

	users := make([]model.User, 0)

	cond, namedArgs := buildConditions(filters)

	sqlStr := `
		SELECT id, email, first_name, last_name, role, created_at, updated_at
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

func (r UserRepository) Create(ctx context.Context, userArgs CreateUserArgs) (*model.User, error) {
	var user model.User

	namedArgs := make(pgx.NamedArgs)
	namedArgs["email"] = userArgs.Email
	namedArgs["passwordHash"] = userArgs.PasswordHash
	namedArgs["firstName"] = userArgs.FirstName
	namedArgs["lastName"] = userArgs.LastName.Value

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

func (r UserRepository) Update(ctx context.Context, id int, userArgs UpdateUserArgs) (*model.User, error) {
	var user model.User

	setFields := make([]string, 0)
	namedArgs := make(pgx.NamedArgs)
	namedArgs["id"] = id

	if userArgs.FirstName.IsSet {
		if userArgs.FirstName.Value == nil {
			return nil, errors.New("first name must not be null")
		}

		setFields = append(setFields, "first_name = @firstName")
		namedArgs["firstName"] = userArgs.FirstName.Value
	}

	if userArgs.LastName.IsSet {
		setFields = append(setFields, "last_name = @lastName")
		namedArgs["lastName"] = userArgs.LastName.Value
	}

	if userArgs.PasswordHash.IsSet {
		setFields = append(setFields, "password_hash = @passwordHash")
		namedArgs["passwordHash"] = userArgs.PasswordHash.Value
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

// helpers
func buildConditions(filters Filters) (string, pgx.NamedArgs) {
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

	if filters.FirstName != nil {
		namedArgs["first_name"] = *filters.FirstName
		cond += " AND first_name = @first_name"
	}

	if filters.LastName != nil {
		namedArgs["last_name"] = *filters.LastName
		cond += " AND last_name = @last_name"
	}

	return cond, namedArgs
}
