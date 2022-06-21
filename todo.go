package main

import (
	"context"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetTodos(context.Context, *emptypb.Empty) (*todo.GetTodosResponse, error) {
	return nil, nil
}
