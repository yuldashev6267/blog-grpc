package main

import (
	"context"
	"fmt"
	"github.com/yuldashev6267/blog-grpc/internals/blogpb"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	severAddress = "SERVER_ADDR"
)

func main() {

	fmt.Println("Client blog is called")

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error reading environmental file %v", err)
	}

	inSecure := grpc.WithInsecure()

	cc, err := grpc.Dial(os.Getenv(severAddress), inSecure)

	defer cc.Close()

	if err != nil {
		log.Fatalf("Error client server %v", err)
	}

	rpcC := blogpb.NewBlogServiceClient(cc)
	CreatBlogClient(rpcC)
}

func CreatBlogClient(b blogpb.BlogServiceClient) (string, error) {
	req := &blogpb.BlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "1",
			Title:    "1984",
			Comment:  "It is really perfect book",
		},
	}

	ctx, can := context.WithTimeout(context.Background(), time.Second*5)
	defer can()
	fmt.Println(req)
	res, err := b.CreateBlog(ctx, req)

	if err != nil {
		return "", fmt.Errorf("Error has been occured sending request", err)
	}
	
	fmt.Println(res)
	return "Blog created", nil
}
