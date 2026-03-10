package user

import userpb "github.com/Hubcher/project-management/contracts/gen/proto/user"

func toCreateUserPB(req CreateUserRequest) *userpb.CreateUserRequest {
	return &userpb.CreateUserRequest{
		Id:       req.ID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
}

func toUpdateUserPB(id string, req UpdateUserRequest) *userpb.UpdateUserRequest {
	return &userpb.UpdateUserRequest{
		Id:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
}

func toUserResponse(user *userpb.User) *UserResponse {
	return &UserResponse{
		ID:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}

func toListUsersResponse(users []*userpb.User) ListUsersResponse {
	resp := ListUsersResponse{
		Users: make([]UserResponse, 0, len(users)),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, *toUserResponse(user))
	}
	return resp
}
