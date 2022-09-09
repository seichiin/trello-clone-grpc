package main

import (
	"context"
	"fmt"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetTodos(context.Context, *emptypb.Empty) (*todo.GetTodosResponse, error) {
	todos := []Todo{}
	tx := s.DB.Find(&todos)

	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error!")
	}

	resTodos := []*todo.Todo{}

	for _, todo := range todos {
		resTodos = append(resTodos, todo.Proto())
	}

	return &todo.GetTodosResponse{Todos: resTodos}, nil
}

func (s *Server) CreateTodo(ctx context.Context, req *todo.Todo) (*todo.Todo, error) {
	todo := &Todo{}
	todo.FromProto(req)

	tx := s.DB.Create(todo)

	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
	}

	return todo.Proto(), nil
}
