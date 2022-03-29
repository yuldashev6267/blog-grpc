package main

import (
	"context"
	"fmt"
	"github.com/yuldashev6267/blog-grpc/internals/blogpb"
	"github.com/yuldashev6267/blog-grpc/internals/repository"
	"github.com/yuldashev6267/blog-grpc/pkg/models"
	"github.com/yuldashev6267/blog-grpc/pkg/api"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	mongoConnectionString = "DATABASE_CONNECTION"
	dataBaseName          = "DATABASE_NAME"
	serverAddress = "GIN_SERVER"
)

//type server struct {
//	blogpb.UnimplementedBlogServiceServer
//}

//type blogService struct {
//	db repository.Database
//}

//type blogInter interface {
//	InsertBlog(model *models.Blog) (*primitive.ObjectID, error)
//}

//func (b *blogService) InsertBlog(blog *models.Blog) (*primitive.ObjectID, error) {
//	var (
//		blogCollection = b.db.BlogCollection()
//	)
//
//	res, err := blogCollection.InsertOne(context.Background(), &blog)
//	if err != nil {
//		return &primitive.NilObjectID, fmt.Errorf("Error %v", err)
//	}
//	oid, ok := res.InsertedID.(primitive.ObjectID)
//	if !ok {
//		return &primitive.NilObjectID, fmt.Errorf("123 %v", err)
//	}
//
//	return &oid, nil
//
//}

//func (*server) CreateBlog(ctx context.Context, req *blogpb.BlogRequest) (*blogpb.BlogResponse, error) {
//	blog := req.GetBlog()
//
//	data := &models.Blog{
//		AuthorId: blog.GetAuthorId(),
//		Title:    blog.GetTitle(),
//		Comment:  blog.GetComment(),
//	}
//
//	var s blogInter
//	s = &blogService{}
//
//	id, err := s.InsertBlog(data)
//	if err != nil {
//		return nil, status.Errorf(
//			codes.Internal,
//			fmt.Sprintf("Internal error has been occured %v", err),
//		)
//	}
//
//	return &blogpb.BlogResponse{
//		Blog: &blogpb.Blog{
//			Id:       id.Hex(),
//			AuthorId: blog.GetAuthorId(),
//			Title:    blog.GetTitle(),
//			Comment:  blog.GetComment(),
//		},
//	}, nil
//}

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

func main() {

	ginCh :=make(chan int)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	listeningMong := createDb()
	go startServer(ginCh)
	defer listeningMong.Disconnect()

	if err != nil {
		log.Fatalf("There is an error listening server %v", err)
	}
	opts := []grpc.ServerOption{}
	blogServer := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(blogServer, &server{})

	go func() {
		if err := blogServer.Serve(lis); err != nil {
			log.Fatalf("Grpc Server Error %v", err)
		}
	}()

	// blocking chanel and close when ctr f5 pressed
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	<-ginCh
	blogServer.Stop()
	fmt.Println("Stoping grpc server")
	lis.Close()
	fmt.Println("Stoping listener")
}
