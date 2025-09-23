package initial

import "github.com/baimhons/stadiumhub/internal/user"

type Validate struct {
	UserValidate user.UserValidate
}

func NewValidate() *Validate {
	return &Validate{
		UserValidate: user.NewUserValidate(),
	}
}
