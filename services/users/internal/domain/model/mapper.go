package model

import (
	userspb "wch/gen/users/v1"

	timeconv "google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(u *User) *userspb.User {
	return &userspb.User{
		Id:           u.ID.String(),
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Surname:      u.Surname,
		AvatarUrl:    u.AvatarURL,
		Status:       u.Status,
		DepartmentId: int32(u.DepartmentID),
		CreatedAt:    timeconv.New(u.CreatedAt),
		UpdatedAt:    timeconv.New(u.UpdatedAt),
	}
}

func UserFromProto(u *userspb.User) (*User, error) {
	return &User{
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Status:       u.Status,
		DepartmentID: int(u.DepartmentId),
		CreatedAt:    u.CreatedAt.AsTime(),
		UpdatedAt:    u.UpdatedAt.AsTime(),
	}, nil
}
