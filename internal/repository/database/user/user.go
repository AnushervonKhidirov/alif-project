package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"pet/internal/constant"
	"pet/internal/model"
	"pet/pkg/nullable"

	"github.com/jackc/pgx/v5"
)

const (
	userFields = "u.id, u.email, u.password_hash, u.first_name, u.last_name, u.role, u.created_at, u.updated_at"
	petFields  = "p.id, p.user_id, p.name, p.about, p.gender, p.birthday, p.microchip_id, p.pet_type_id, p.breed_id, p.lost, p.created_at, p.updated_at"
	userPrefix = "u"
	petPrefix  = "p"
)

type FilterOne struct {
	Id    *int
	Email *string
}

type FilterMany struct {
	FirstName *string
	LastName  *string
	Email     *string

	Limit  int
	Offset int
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

func (r UserRepository) GetOneWithPet(ctx context.Context, filters FilterOne) (*model.User, error) {
	var user model.User

	cond, namedArgs := buildOneConditions(filters)
	sqlStr := `
		SELECT ` + userFields + `, ` + petFields + `, pt.id, pt.value, b.id, b.value, b.pet_type_id
		FROM ` + r.tableName + " AS " + userPrefix + `
		LEFT JOIN ` + constant.TablePet + " AS " + petPrefix + `
			ON u.id = p.user_id
		LEFT JOIN ` + constant.TablePetType + ` AS pt
			ON p.pet_type_id = pt.id
		LEFT JOIN ` + constant.TableBreed + ` AS b
			ON p.breed_id = b.id
		WHERE
	` + cond + ";"

	rows, err := r.conn.Query(ctx, sqlStr, namedArgs)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	user.Pets = make([]model.Pet, 0)
	found := false

	for rows.Next() {
		var petID, petUserID, petTypeID, joinedPetTypeID *int
		var petBreedID, breedID, breedPetTypeID *int
		var petName *string
		var petAbout *string
		var petGender *model.Gender
		var petBirthday *time.Time
		var petMicrochipID *string
		var petLost *bool
		var petCreatedAt, petUpdatedAt *time.Time
		var petTypeValue *string
		var breedValue *string

		err := rows.Scan(
			&user.Id,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&petID,
			&petUserID,
			&petName,
			&petAbout,
			&petGender,
			&petBirthday,
			&petMicrochipID,
			&petTypeID,
			&petBreedID,
			&petLost,
			&petCreatedAt,
			&petUpdatedAt,
			&joinedPetTypeID,
			&petTypeValue,
			&breedID,
			&breedValue,
			&breedPetTypeID,
		)
		if err != nil {
			return nil, err
		}

		found = true

		if petID == nil {
			continue
		}

		if petUserID == nil || petName == nil || petTypeID == nil || petLost == nil ||
			petCreatedAt == nil || petUpdatedAt == nil || joinedPetTypeID == nil || petTypeValue == nil {
			return nil, errors.New("incomplete pet data returned by query")
		}

		pet := model.Pet{
			Id:          *petID,
			UserId:      *petUserID,
			Name:        *petName,
			About:       petAbout,
			Gender:      petGender,
			Birthday:    petBirthday,
			MicrochipId: petMicrochipID,
			Lost:        *petLost,
			PetTypeId:   *petTypeID,
			BreedId:     petBreedID,
			CreatedAt:   *petCreatedAt,
			UpdatedAt:   *petUpdatedAt,
		}
		pet.PetType = &model.PetType{Id: *joinedPetTypeID, Value: *petTypeValue}

		if breedID != nil {
			if breedPetTypeID == nil || breedValue == nil {
				return nil, errors.New("incomplete breed data returned by query")
			}

			pet.Breed = &model.Breed{
				Id:        *breedID,
				PetTypeId: *breedPetTypeID,
				Value:     *breedValue,
			}
		}
		user.Pets = append(user.Pets, pet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if !found {
		return nil, pgx.ErrNoRows
	}

	return &user, nil
}

func (r UserRepository) GetOne(ctx context.Context, filters FilterOne) (*model.User, error) {
	var user model.User

	cond, namedArgs := buildOneConditions(filters)
	sqlStr := `
		SELECT ` + userFields + `
		FROM ` + r.tableName + " AS " + userPrefix + `
		WHERE
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
	users := make([]model.User, 0)
	cond, namedArgs := buildManyConditions(filters)

	sqlStr := `
		SELECT ` + userFields + `
		FROM ` + r.tableName + " AS " + userPrefix + `
		WHERE ` + cond + `
		ORDER BY u.id ASC
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
		var u model.User

		err := rows.Scan(
			&u.Id,
			&u.Email,
			&u.PasswordHash,
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

	if len(setFields) == 0 {
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
	sqlStr := "SELECT password_hash FROM " + r.tableName + " WHERE id = $1;"
	err := r.conn.QueryRow(ctx, sqlStr, id).Scan(&password)
	return password, err
}

// helpers
func buildOneConditions(filters FilterOne) (string, pgx.NamedArgs) {
	cond := "TRUE"
	namedArgs := make(pgx.NamedArgs)

	if filters.Id != nil {
		namedArgs["id"] = *filters.Id
		cond += " AND " + userPrefix + ".id = @id"
	}

	if filters.Email != nil {
		namedArgs["email"] = *filters.Email
		cond += " AND " + userPrefix + ".email = @email"
	}

	return cond, namedArgs
}

func buildManyConditions(filters FilterMany) (string, pgx.NamedArgs) {
	cond := "TRUE"
	namedArgs := make(pgx.NamedArgs)

	if filters.FirstName != nil {
		namedArgs["firstName"] = *filters.FirstName
		cond += " AND " + userPrefix + ".first_name = @firstName"
	}

	if filters.LastName != nil {
		namedArgs["lastName"] = *filters.LastName
		cond += " AND " + userPrefix + ".last_name = @lastName"
	}

	if filters.Email != nil {
		namedArgs["email"] = *filters.Email + "%"
		cond += " AND " + userPrefix + ".email ILIKE @email"
	}

	return cond, namedArgs
}
