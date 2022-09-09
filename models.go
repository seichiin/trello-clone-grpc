package main

import (
	"time"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	Board struct {
		ID    int64  `gorm:"primaryKey;autoIncrement:true"`
		Name  string `gorm:"not null;unique"`
		Order int32  `gorm:"unique"`
	}
	Todo struct {
		ID         int64 `gorm:"primaryKey;autoIncrement:true"`
		ExpireTime time.Time
		BoardID    int64  `gorm:"not null"`
		Name       string `gorm:"not null;unique"`
		Priority   int32  `gorm:"not null"`
		Description string `gorm:"size:120"`
	}
)

func (b *Board) Proto() *todo.Board {
	return &todo.Board{
		Id:    b.ID,
		Name:  b.Name,
		Order: b.Order,
	}
}

func (t *Todo) Proto() *todo.Todo {
	return &todo.Todo{
		Id:         t.ID,
		BoardId:    t.BoardID,
		Name:       t.Name,
		Priority:   t.Priority,
		ExpireTime: timestamppb.New(t.ExpireTime),
		Description: t.Description,
	}
}

func (t *Todo) FromProto(v1 *todo.Todo) {
	t.BoardID = v1.BoardId
	t.Name = v1.Name
	t.Priority = v1.Priority
	t.ExpireTime = v1.ExpireTime.AsTime()
	t.Description = v1.Description
}
