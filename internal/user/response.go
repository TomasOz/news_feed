package user

type UserResponse struct {
	Username string `json:"username"`
}

type TokenResponse struct {
    Token string `json:"token"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func mapUserToResponse(user *User) UserResponse {
	return UserResponse{Username: user.Username}
}