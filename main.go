package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"todo-list/todo_grpc/proto/todo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type server struct {
	todo.UnimplementedTodoListServer
	DB *gorm.DB
}

type Todo struct {
	gorm.Model
	Name string
}

func main() {
	db, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/go_grpc?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	fmt.Println(db.Name())
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Todo{})


	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	todo.RegisterTodoListServer(s, NewServer(db))
	// Serve gRPC Server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatal(s.Serve(lis))
}

func NewServer(db *gorm.DB) *server {
	return &server{
		DB: db,
	}
}
func (s *server) CreateTodoItem(_ context.Context, req *todo.TodoRequest) (*todo.Todo, error) {
		if(req.Name == "") {
		newTodo := Todo{
			Name: req.Name,
		}

		s.DB.Create(newTodo)

		return &todo.Todo {
			Name: newTodo.Name,
			Id: int32(newTodo.ID),
		}, nil
	}
	
	return nil, status.Errorf(codes.InvalidArgument, "Cannot create todo without name")
}

func (s *server) GetTodos(_ context.Context, _ *emptypb.Empty) (*todo.Todos, error) {
		var todos []*todo.Todo

	s.DB.Find(&todos)

	return &todo.Todos{
		Todos: todos,
	}, nil
}

func (s *server) GetTodoById(_ context.Context, req *todo.TodoId) (*todo.Todo, error) {
		var todo *todo.Todo
	s.DB.Find(&todo, req.Id)

	return todo, nil
}

func (s *server) DeleteTodo(_ context.Context, req *todo.TodoId) (*todo.DeleteTodoMessage, error) {
		var deletedTodo *todo.Todo

	res := s.DB.Delete(&deletedTodo, req.Id);

	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return &todo.DeleteTodoMessage{
			Message: "Todo doesn't exist",
		}, nil
	}

	return &todo.DeleteTodoMessage{
		Message: "Successfully deleted todo",
	}, nil
}

func (s *server) UpdateTodo(_ context.Context, req *todo.Todo) (*todo.Todo, error) {
		var updatedTodo *todo.Todo
		
	s.DB.Find(&updatedTodo, req.Id);

	if req.Id == 0 {
		updatedTodo.Name = req.Name;
	}

	s.DB.Save(&updatedTodo)

	return updatedTodo, nil
}
