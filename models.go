package main

import (
	"time"
	"todo-list/pkg/idgen"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type (
	User struct {
		ID       int64  `gorm:"primaryKey;autoIncrement:true"`
		UserName string `gorm:"size:40;unique:not null"`
		Password string `gorm:"not null"`
		Email    string `gorm:"unique;not null"`
	}
	Board struct {
		ID int64 `gorm:"primaryKey;autoIncrement:true"`

		User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
		UserID int64

		Name  string `gorm:"not null;size:30"`
		Order int32  `gorm:"not null"`
	}
	Todo struct {
		ID int64 `gorm:"primaryKey;autoIncrement:true"`

		Board   Board `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
		BoardID int64

		Name     string `gorm:"not null;size:30"`
		Priority string `gorm:"not null"`

		StartTime   time.Time
		ExpireTime  time.Time
		Description string `gorm:"size:120"`
		Color       string
		Order       int32 `gorm:"not null"`
		Completed      bool 
	}
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == 0 {
		u.ID = idgen.GenID()
	}
	return nil
}

func (t *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == 0 {
		t.ID = idgen.GenID()
	}
	return nil
}

func (b *Board) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == 0 {
		b.ID = idgen.GenID()
	}
	return nil
}

func (b *Board) Proto() *todo.Board {
	return &todo.Board{
		Id:    b.ID,
		Name:  b.Name,
		Order: b.Order,
		UserId: b.UserID,
	}
}

func (t *Todo) Proto() *todo.Todo {
	return &todo.Todo{
		Id:          t.ID,
		BoardId:     t.BoardID,
		Name:        t.Name,
		Priority:    todo.Todo_Priority(todo.Todo_Priority_value[t.Priority]),
		StartTime:   timestamppb.New(t.StartTime),
		ExpireTime:  timestamppb.New(t.ExpireTime),
		Description: t.Description,
		Color:       t.Color,
		Order:       t.Order,
		Completed:   t.Completed,
	}
}

func (t *Todo) FromProto(v1 *todo.Todo) {
	t.ID = v1.Id
	t.BoardID = v1.BoardId
	t.Name = v1.Name
	t.Priority = todo.Todo_Priority_name[int32(v1.Priority)]
	t.StartTime = v1.StartTime.AsTime()
	t.ExpireTime = v1.ExpireTime.AsTime()
	t.Description = v1.Description
	t.Color = v1.Color
	t.Order = v1.Order
	t.Completed = v1.Completed
}

func (b *Board) FromProto(v1 *todo.Board) {
	b.ID = v1.Id
	b.Name = v1.Name
	b.Order = v1.Order
	b.UserID = v1.UserId
}
