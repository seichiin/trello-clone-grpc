package todo

import (
	"context"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetTodos(context.Context, *emptypb.Empty) (*todo.Todos, error) {
	var todolist []*todo.Todo
	s.DB.Find(&todolist)

	return &todo.Todos{Todos: todolist}, nil
}