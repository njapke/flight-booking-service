package database

type Model interface {
	Collection() string
	Key() string
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u *User) Collection() string {
	return "users"
}

func (u *User) Key() string {
	return u.ID
}
