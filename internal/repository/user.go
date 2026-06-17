package repository

import (
	"database/sql"
	"errors"

	"github.com/CodingFervor/carbon-emission-management/internal/model"
)

type UserRepo struct {
	*GenericRepo[model.User]
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{GenericRepo: NewGenericRepo[model.User](db, "users", func() *model.User { return &model.User{} })}
}

// ErrUserNotFound is returned when no user matches a lookup.
var ErrUserNotFound = errors.New("user not found")

// FindByUsername loads a user by username (used for login).
func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	cols := r.columns
	q := "SELECT " + join(cols) + " FROM users WHERE username = $1"
	u := &model.User{}
	ptrs, finalize := structPtrs(u, cols)
	if err := r.DB().QueryRow(q, username).Scan(ptrs...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	finalize()
	return u, nil
}

// Create inserts a new user with a pre-hashed password.
func (r *UserRepo) Create(u *model.User) error {
	q := `INSERT INTO users (username, password, email, role, organization_id, status)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`
	return r.DB().QueryRow(q,
		u.Username, u.Password, u.Email, u.Role, u.OrganizationID, defaultStr(u.Status, "active"),
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

// Update modifies mutable user fields.
func (r *UserRepo) Update(u *model.User) error {
	q := `UPDATE users SET email=$1, role=$2, status=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$4`
	_, err := r.DB().Exec(q, u.Email, u.Role, defaultStr(u.Status, "active"), u.ID)
	return err
}

func join(cols []string) string {
	out := ""
	for i, c := range cols {
		if i > 0 {
			out += ", "
		}
		out += c
	}
	return out
}

func defaultStr(s, d string) string {
	if s == "" {
		return d
	}
	return s
}
