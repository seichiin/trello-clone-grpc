package todo

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
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	todo.UnimplementedTodoListServer
	DB *gorm.DB
}

func main() {
	dbConnection := "root:root@tcp(127.0.0.1:3306)/todolist_grpc?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})
	fmt.Println(db.Name())
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&todo.Todo{})

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

func NewServer(db *gorm.DB) *Server {
	return &Server{
		DB: db,
	}
}
