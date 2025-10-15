package response

import "github.com/google/uuid"

type LoginUserResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Role        string    `json:"role"`
}
