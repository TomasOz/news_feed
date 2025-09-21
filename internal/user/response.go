package user

type UserResponse struct {
	Username string `json:"username"`
}

func mapUserToResponse(user *User) UserResponse {
	return UserResponse{Username: user.Username}
}