package main

import (
	"context"
	"fmt"
	"todo-list/todo_grpc/proto/todo"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetTodos(context.Context, *emptypb.Empty) (*todo.GetTodosResponse, error) {
    todos := []Todo{}
    tx := s.DB.Find(&todos)

    if tx.Error != nil {
        return nil, fmt.Errorf("Internal Server Error!")
    }

    resTodos := lo.Map(todos, func(todo Todo, _ int) *todo.Todo {
		return todo.Proto()
	})

    return &todo.GetTodosResponse{Todos: resTodos}, nil
}
