package request

type RegisterUser struct {
	Username        string `json:"username" validate:"required,min=5,max=20,alphanum"`
	FirstName       string `json:"first_name" validate:"required,min=3,max=30,alphaunicode"`
	LastName        string `json:"last_name" validate:"required,min=3,max=30,alphaunicode"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,eqfield=Password"`
	PhoneNumber     string `json:"phone_number" validate:"required,numeric,len=10"`
}

type LoginUser struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required,min=6"`
}

type UpdateUser struct {
	Username    string `json:"username" validate:"required,min=5,max=20,alphanum"`
	FirstName   string `json:"first_name" validate:"required,min=3,max=30,alphaunicode"`
	LastName    string `json:"last_name" validate:"required,min=3,max=30,alphaunicode"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric,len=10"`
}
