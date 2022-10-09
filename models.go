package main

import (
	"time"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	User struct {
		ID       int64  `gorm:"primaryKey;autoIncrement:true"`
		UserName string `gorm:"size:40;unique"`
		Password string `gorm:"not null"`
		Email    string `gorm:"unique;not null"`
	}
	Board struct {
		ID int64 `gorm:"primaryKey;autoIncrement:true"`

		User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
		UserID int64

		Name  string `gorm:"not null;unique"`
		Order int32  `gorm:"unique"`
	}
	Todo struct {
		ID int64 `gorm:"primaryKey;autoIncrement:true"`

		Board   Board `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
		BoardID int64

		Name     string `gorm:"not null;unique"`
		Priority int32  `gorm:"not null"`

		StartTime   time.Time
		ExpireTime  time.Time
		Description string `gorm:"size:120"`
		Color       string
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
		Id:          t.ID,
		BoardId:     t.BoardID,
		Name:        t.Name,
		Priority:    t.Priority,
		StartTime:   timestamppb.New(t.ExpireTime),
		ExpireTime:  timestamppb.New(t.ExpireTime),
		Description: t.Description,
		Color:       t.Color,
	}
}

func (t *Todo) FromProto(v1 *todo.Todo) {
	t.BoardID = v1.BoardId
	t.Name = v1.Name
	t.Priority = v1.Priority
	t.StartTime = v1.ExpireTime.AsTime()
	t.ExpireTime = v1.ExpireTime.AsTime()
	t.Description = v1.Description
	t.Color = v1.Color
}
