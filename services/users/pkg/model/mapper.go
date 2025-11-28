package model

import (
	"wch/gen"

	timeconv "google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(u *User) *gen.User {
	return &gen.User{
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

func UserFromProto(u *gen.User) (*User, error) {
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
