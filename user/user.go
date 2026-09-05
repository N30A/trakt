package user

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Role         UserRole
}

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)
