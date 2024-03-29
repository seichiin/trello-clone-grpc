package main

import (
	"context"
	"fmt"
	"todo-list/todo_grpc/proto/todo"

	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetBoards(ctx context.Context, req *todo.GetBoardsRequest) (*todo.GetBoardsResponse, error) {
	boards := []Board{}
	tx := s.DB.Where(Board{UserID: req.UserId}).Order("`order` asc").Find(&boards)
	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error!")
	}

	resBoards := lo.Map(boards, func(board Board, _ int) *todo.Board {
		return board.Proto()
	})

	return &todo.GetBoardsResponse{
		Boards: resBoards,
	}, nil
}

func (s *Server) DeleteBoard(ctx context.Context, req *todo.BoardDetailRequest) (*emptypb.Empty, error) {
	tx := s.DB.Where("user_id = ? AND id = ?", req.UserId, req.Id).Delete(&Board{})
	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, fmt.Errorf("Board doesn't exist!")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateBoard(ctx context.Context, req *todo.UpdateBoardRequest) (*todo.Board, error) {
	maskes := []string{}
	if req.UpdateMask != nil {
		maskes = req.UpdateMask.Paths
	}

	board := Board{}
	board.FromProto(req.Board)

	tx := s.DB.Select(maskes).Updates(&board)
	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
	}

	return board.Proto(), nil
}

func (s *Server) CreateBoard(ctx context.Context, req *todo.Board) (*todo.Board, error) {
	var maxOrder int32
	row := s.DB.Table("boards").Where("user_id = ? ", req.UserId).Select("MAX(`order`)").Row()
	err := row.Scan(&maxOrder)
	if err != nil {
		maxOrder = 0
	}
	board := &Board{}
	board.FromProto(req)
	board.Order = maxOrder + 1

	tx := s.DB.Create(board)

	if tx.Error != nil {
		return nil, fmt.Errorf("Internal Server Error: %v", tx.Error)
	}

	return board.Proto(), nil
}
