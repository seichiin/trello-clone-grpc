package main

import (
	"context"
	"errors"
	"fmt"
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
	tx := s.DB.Where("user_name = ? OR email = ?", req.Username, req.Email).Limit(1).Find(&user)
	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("User not found!")
	}

	err := CheckPasswordHash(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("Password is invalid!")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ChangePassword(ctx context.Context, req *todo.User) (*emptypb.Empty, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("Internal Server Error: %v", err)
	}

	tx := s.DB.Model(User{}).Where("user_name = ? OR email = ?", req.Username, req.Email).Update("password", hashedPassword)
	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("User not found!")
	}

	return &emptypb.Empty{}, nil
}