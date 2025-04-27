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
		Password: "adminpass", // en futuro hashed
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
