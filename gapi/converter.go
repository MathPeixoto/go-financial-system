package gapi

import (
	db "github.com/MathPeixoto/go-financial-system/db/sqlc"
	"github.com/MathPeixoto/go-financial-system/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func converter(dbUser db.User) *pb.User {
	return &pb.User{
		Username:          dbUser.Username,
		FullName:          dbUser.FullName,
		Email:             dbUser.Email,
		PasswordChangedAt: timestamppb.New(dbUser.PasswordChangedAt),
		CreatedAt:         timestamppb.New(dbUser.CreatedAt),
	}
}
