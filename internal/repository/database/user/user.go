package user

import (
	"context"
	"database/sql"
	"fmt"
	"pet/internal/model"

	"github.com/jackc/pgx/v5"
)

type UserFilter struct {
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
	LastName     *string
}

type UpdateUserArgs struct {
	FirstName sql.NullString
	LastName  *string
}

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(c *pgx.Conn) *UserRepository {
	return &UserRepository{conn: c}
}

func (c UserRepository) GetOne(ctx context.Context, filters UserFilter) (*model.User, error) {
	var user model.User

	cond, namedArgs := buildUserConditions(filters)
	sqlStr := `
		SELECT id, email, first_name, last_name, role, created_at, updated_at FROM users
		WHERE
	` + cond + " ORDER BY id ASC LIMIT 1;"

	err := c.conn.QueryRow(ctx, sqlStr, namedArgs).Scan(
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

func (c UserRepository) GetMany(ctx context.Context, filters UserFilter) ([]model.User, error) {
	if filters.Limit != nil && *filters.Limit < 0 {
		return nil, fmt.Errorf("user query limit must not be negative: %d", *filters.Limit)
	}
	if filters.Offset != nil && *filters.Offset < 0 {
		return nil, fmt.Errorf("user query offset must not be negative: %d", *filters.Offset)
	}

	users := make([]model.User, 0)

	cond, namedArgs := buildUserConditions(filters)

	sqlStr := `
		SELECT id, email, first_name, last_name, role, created_at, updated_at FROM users
		WHERE
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

	rows, err := c.conn.Query(ctx, sqlStr, namedArgs)
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

// TODO: complete Create
func (c UserRepository) Create(ctx context.Context, userArgs CreateUserArgs) (*model.User, error) {
	var user model.User
	sqlStr := `
		INSERT INTO users (email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, first_name, last_name, role, created_at, updated_at;
	`

	err := c.conn.QueryRow(ctx, sqlStr, userArgs.Email, userArgs.PasswordHash, userArgs.FirstName, userArgs.LastName).Scan(
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

// TODO: complete Update
func (c UserRepository) Update(ctx context.Context, userId int, userArgs UpdateUserArgs) (*model.User, error) {
	var user model.User
	sqlStr := `
		UPDATE users SET 
		WHERE id = $1
		RETURNING id, email, first_name, last_name, role, created_at, updated_at;
	`

	err := c.conn.QueryRow(ctx, sqlStr, userArgs.FirstName, userArgs.LastName).Scan(
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

func (c UserRepository) Delete(ctx context.Context, userId int) error {
	sqlStr := "DELETE FROM users WHERE id = $1;"
	_, err := c.conn.Exec(ctx, sqlStr, userId)
	return err
}

// helpers
func buildUserConditions(filters UserFilter) (string, pgx.NamedArgs) {
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
