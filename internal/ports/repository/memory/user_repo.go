package memory

import (
	"errors"
	"salle-songbook-api/internal/core/user"

	"github.com/google/uuid"
)

type UserRepository struct {
	users map[string]user.User
}

func NewUserRepository() *UserRepository {
	repo := &UserRepository{users: make(map[string]user.User)}
	repo.seed()
	return repo
}

func (r *UserRepository) seed() {
	r.users["admin"] = user.User{
		ID:       uuid.NewString(),
		Username: "admin",
		Password: "adm1n@LaSalle",
		Role:     user.Admin,
	}
	r.users["pastoral-sjb"] = user.User{
		ID:       uuid.NewString(),
		Username: "pastoral-sjb",
		Password: "sJb@2025P4st",
		Role:     user.Admin,
	}
	r.users["pastoral-dls"] = user.User{
		ID:       uuid.NewString(),
		Username: "pastoral-dls",
		Password: "dLs@2025P4st",
		Role:     user.Admin,
	}
	r.users["pastoral-uls"] = user.User{
		ID:       uuid.NewString(),
		Username: "pastoral-uls",
		Password: "uLs@2025P4st",
		Role:     user.Admin,
	}
	r.users["pastoral-sjls"] = user.User{
		ID:       uuid.NewString(),
		Username: "pastoral-sjls",
		Password: "sjLs@2025P4st",
		Role:     user.Admin,
	}
	r.users["composer"] = user.User{
		ID:       uuid.NewString(),
		Username: "composer",
		Password: "composerpass",
		Role:     user.Composer,
	}
}

func (r *UserRepository) GetByUsername(username string) (user.User, error) {
	u, ok := r.users[username]
	if !ok {
		return user.User{}, errors.New("user not found")
	}
	return u, nil
}
