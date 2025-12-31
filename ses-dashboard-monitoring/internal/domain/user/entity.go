package user

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Don't serialize password
	Email    string `json:"email"`
	Role     string `json:"role"`
	Active   bool   `json:"active"`
}