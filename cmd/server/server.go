package main

import (
	"fmt"
	"github.com/yuldashev6267/blog-grpc/internals/blogpb"
	"github.com/yuldashev6267/blog-grpc/internals/repository"
	"github.com/yuldashev6267/blog-grpc/pkg/api"
	"github.com/yuldashev6267/blog-grpc/pkg/blog_server"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

const (
	mongoConnectionString = "DATABASE_CONNECTION"
	dataBaseName          = "DATABASE_NAME"
	serverAddress = "GIN_SERVER"
	grpcServerAddr = "SERVER_ADDR"
)

func createDb() repository.Database {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error reading env file %v", err)
	}

	return repository.DatabaseConn(os.Getenv(mongoConnectionString), os.Getenv(dataBaseName))
}

func startServer(ch chan int){
	server:= api.New(os.Getenv(serverAddress))
	err := server.RegisterGinServer()
	if err != nil {
		log.Fatalf("Listening Server Error %v", err)
	}
	fmt.Println("Listening gin server")
	ch <- 0
}

type rpcServer struct {
	blogpb.UnimplementedBlogServiceServer
	blo blog_server.BlogService
}

func startGrpcServer(ch chan int){
	l, err := net.Listen("tcp", os.Getenv(grpcServerAddr))

	if err != nil {
		log.Fatal("Error listening server")
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)

	if err := s.Serve(l); err != nil {
		log.Fatalf("Error has been occured during rpc serve %v", err)
	}
	blogpb.RegisterBlogServiceServer(s, &rpcServer{})
	fmt.Println("Grpc Started")
	ch<-0
}

func main() {

	ginCh :=make(chan int)
	rpcCh := make(chan int,1)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listeningMong := createDb()
	go startServer(ginCh)
	defer listeningMong.Disconnect()
	go startGrpcServer(rpcCh)
	<-ginCh
	<-rpcCh
}
