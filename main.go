package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"todo-list/todo_grpc/proto/todo"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	db, err := gorm.Open(mysql.Open("root:akashi@tcp(127.0.0.1:3306)/go_grpc?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
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
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = todo.RegisterTodoListHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

func NewServer(db *gorm.DB) *server {
	return &server{
		DB: db,
	}
}

func (s *server) GetTodos(context.Context, *emptypb.Empty) (*todo.Todos, error) {
	var todolist []*todo.Todo
	s.DB.Find(&todolist)
	return &todo.Todos{Todos: todolist}, nil
}
