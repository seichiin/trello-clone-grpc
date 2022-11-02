package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"todo-list/pkg/idgen"
	"todo-list/todo_grpc/proto/todo"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	todo.UnimplementedUserServiceServer
	todo.UnimplementedTodoServiceServer
	todo.UnimplementedBoardServiceServer
	DB *gorm.DB
}

func allowedOrigin(origin string) bool {
    if viper.GetString("cors") == "*" {
        return true
    }
    if matched, _ := regexp.MatchString(viper.GetString("cors"), origin); matched {
        return true
    }
    return false
}

func cors(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if allowedOrigin(r.Header.Get("Origin")) {
            w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
            w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType")
        }
        if r.Method == "OPTIONS" {
            return
        }
        h.ServeHTTP(w, r)
    })
}

func main() {
	dbName := "todolist_grpc"
	dbConnection := fmt.Sprintf("root:root@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbName)
	db, err := gorm.Open(mysql.Open(dbConnection), &gorm.Config{})
	
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Todo{}, &Board{}, &User{})

	lis, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	idgen.MustInit(1)

	s := grpc.NewServer()
	todo.RegisterTodoServiceServer(s, &Server{
		DB: db,
	})
	todo.RegisterBoardServiceServer(s, &Server{
		DB: db,
	})
	todo.RegisterUserServiceServer(s, &Server{
		DB: db,
	})
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
		log.Fatalln("Failed tooo dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = todo.RegisterTodoServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	err = todo.RegisterBoardServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	err = todo.RegisterUserServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: cors(gwmux),
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
