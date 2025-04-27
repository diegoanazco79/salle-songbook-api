package user

type Role string

const (
	Admin    Role = "admin"
	Composer Role = "composer"
)

type User struct {
	ID       string
	Username string
	Password string // hashed (por ahora plaintext para facilidad)
	Role     Role
}
