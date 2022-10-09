package main

import (
	"context"
	"errors"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) SignUp(ctx context.Context, req *todo.User) (*emptypb.Empty, error) {
	if !ValidateEmail(req.Email) {
		return nil, errors.New("Invalid email, please correct your email format!")
	}
	if len(req.Username) < 6 {
		return nil, errors.New("Username must be at least 6 characters!") 
	}
	if len(req.Password) < 6 {
		return nil, errors.New("Password must be at least 6 characters!") 
	}

	username := req.Username
	password, err := HashPassword(req.Password) 
	if err != nil {
		return nil, errors.New("Some thing went wrong!")
	}

	tx := s.DB.Create(&User{
		Email: req.Email,
		UserName: Santize(username),
		Password: Santize(password),
	})
	if tx.Error != nil {
		return nil,	tx.Error
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) SignIn(ctx context.Context, req *todo.User) (*emptypb.Empty, error) {
	user := User{}
	tx := s.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&user)
	if tx.Error != nil {
		return nil,	tx.Error
	}

	err := CheckPasswordHash(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("Password is invalid!")
	}

	return &emptypb.Empty{}, nil
}