package main

import (
	"context"
	"fmt"
	"todo-list/todo_grpc/proto/todo"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetTodos(ctx context.Context, req *todo.GetTodosRequest) (*todo.GetTodosResponse, error) {
	tx := s.DB.Where("board_id = ?", req.BoardId)
	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error!")
	}

	if req.FilterName != "" {
		tx = tx.Where("name LIKE ?", "%"+req.FilterName+"%")
		if tx.Error != nil {
			return nil, fmt.Errorf("Internal Server Error!")
		}
	}
	if req.FilterPriority != "" {
		tx = tx.Where("priority LIKE ?", "%"+req.FilterPriority+"%")
		if tx.Error != nil {
			return nil, fmt.Errorf("Internal Server Error!")
		}
	}
	if req.FilterCompleted != "" {
		tx = tx.Where("completed LIKE ?", "%"+req.FilterCompleted+"%")
		if tx.Error != nil {
			return nil, fmt.Errorf("Internal Server Error!")
		}
	}

	todos := []Todo{}
	tx = tx.Order("`order` asc").Find(&todos)

	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v",tx.Error)
	}

	resTodos := lo.Map(todos, func(todo Todo, _ int) *todo.Todo {
		return todo.Proto()
	})

	return &todo.GetTodosResponse{Todos: resTodos}, nil
}

func (s *Server) GetTodoDetail(ctx context.Context, req *todo.TodoDetailRequest) (*todo.Todo, error) {
	todo := Todo{}
	tx := s.DB.Where("board_id = ? AND id = ?", req.BoardId, req.Id).Find(&todo)

	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v",tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("Todo doesn't exist!")
	}
	
	return todo.Proto(), nil
}

func (s *Server) DeleteTodo(ctx context.Context, req *todo.TodoDetailRequest) (*emptypb.Empty, error) {
	tx := s.DB.Where("board_id = ? AND id = ?", req.BoardId, req.Id).Delete(&Todo{})

	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v",tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("Todo doesn't exist!")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateTodo(ctx context.Context, req *todo.UpdateTodoRequest) (*todo.Todo, error){
	if req.Todo.Order != 0 {
		dupTodo := Todo{}
		tx := s.DB.Where("board_id = ? AND `order` = ?", req.BoardId, req.Todo.Order).Limit(1).Find(&dupTodo)
		if tx.Error != nil {
			return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
		}
		if tx.RowsAffected != 0 {
			return nil, fmt.Errorf("Order is duplicated")
		}
	}

	maskes := []string{}
	if req.UpdateMask != nil {
		maskes = req.UpdateMask.Paths
	}
	todo := Todo{}
	todo.FromProto(req.Todo)

	tx := s.DB.Select(maskes).Updates(&todo)
	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v",tx.Error)
	}
	
	return todo.Proto(), nil
}

func (s *Server) CreateTodo(ctx context.Context, req *todo.Todo) (*todo.Todo, error) {
	var maxOrder int32
	row := s.DB.Table("todos").Where("board_id = ? ", req.BoardId).Select("MAX(`order`)").Row()
	err := row.Scan(&maxOrder)
	if err != nil {
		maxOrder = 0
	}
	todo := &Todo{}
	todo.FromProto(req)
	todo.Order = maxOrder + 1
	todo.Completed = false
	todo.Priority = "LOW"

	tx := s.DB.Create(todo)
	if tx.Error != nil {
        return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
    }

    return todo.Proto(), nil
}
