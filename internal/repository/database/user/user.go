package user

import (
	"context"
	"database/sql"
	"pet/internal/model"
	"strconv"

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

	cond := "1=1"
	namedArgs := make(map[string]string)

	if filters.Id != nil {
		namedArgs["id"] = strconv.Itoa(*filters.Id)
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

	sqlStr := `
		SELECT id, email, first_name, last_name, role, created_at, updated_at FROM users
		WHERE
	` + cond + " LIMIT 1;"

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
	users := make([]model.User, 0)

	cond := "1=1"
	namedArgs := make(map[string]string)

	if filters.Id != nil {
		namedArgs["id"] = strconv.Itoa(*filters.Id)
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

	if filters.Offset != nil {
		namedArgs["offset"] = strconv.Itoa(*filters.Offset)
		cond += " OFFSET = @offset"
	}

	if filters.Limit != nil {
		namedArgs["limit"] = strconv.Itoa(*filters.Limit)
		cond += " LIMIT = @limit"
	}

	sqlStr := `
		SELECT id, email, first_name, last_name, role, created_at, updated_at FROM users
		WHERE
	` + cond + ";"

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
