package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/yuldashev6267/blog-grpc/blogpb"
	"github.com/yuldashev6267/blog-grpc/db"
	"github.com/yuldashev6267/blog-grpc/models"
	"google.golang.org/grpc"
)

type server struct {
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.BlogRequest) (*blogpb.BlogResponse, error) {
	blog := req.GetBlog()

	data := models.Blog{
		// Id:       blog.GetAuthorId(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Comment:  blog.GetComment(),
	}
	_ = data
	return *blogpb.BlogResponse{
		Blog: data,
	}, nil
	return nil, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	db.DatabaseConn("grpcBlog")
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
