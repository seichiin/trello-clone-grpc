package main

import (
	"context"
	"errors"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) SignUp(ctx context.Context, req *todo.User) (*emptypb.Empty, error) {
	username := req.Username
	password, err := HashPassword(req.Password) 

	if err != nil {
		return nil, errors.New("Some thing went wrong!")
	}

	tx := s.DB.Create(&User{
		UserName: username,
		Password: password,
	})

	if tx.Error != nil {
		return nil,	tx.Error
	}

	return &emptypb.Empty{}, nil
}