package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/yuldashev6267/blog-grpc/blogpb"
	"github.com/yuldashev6267/blog-grpc/db"
	"github.com/yuldashev6267/blog-grpc/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	mongoConnectionString = "DATABASE_CONNECTION"
	dataBaseName          = "DATABASE_NAME"
)

type server struct {}

type blogService struct {
	db db.Database
}

type blogInter interface {
	InsertBlog(model *models.Blog) (*primitive.ObjectID, error)
}

func (b *blogService) InsertBlog(blog *models.Blog) (*primitive.ObjectID, error) {

	blogCollection := b.db.BlogCollection()

	res, err := blogCollection.InsertOne(context.Background(), &blog)

	if err != nil {
		return &primitive.NilObjectID, nil
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return &primitive.NilObjectID, nil
	}
	return &oid, nil

}

func (*server) CreateBlog(ctx context.Context, req *blogpb.BlogRequest) (*blogpb.BlogResponse, error) {
	blog := req.GetBlog()

	data := models.Blog{
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Comment:  blog.GetComment(),
	}
	var s blogInter
	s = &blogService{}

	id, err := s.InsertBlog(&data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error has been occured %v", err),
		)
	}

	return &blogpb.BlogResponse{
		Blog: &blogpb.Blog{
			Id:       id.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Comment:  blog.GetComment(),
		},
	}, nil
}

func createDb() db.Database {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error reading env file %v", err)
	}

	return db.DatabaseConn(os.Getenv(mongoConnectionString), os.Getenv(dataBaseName))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	listeningMong := createDb()
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

	blogServer.Stop()
	fmt.Println("Stoping grpc server")
	lis.Close()
	fmt.Println("Stoping listener")
}
