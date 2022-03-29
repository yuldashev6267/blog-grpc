package blog_server

import (
	"context"
	"github.com/yuldashev6267/blog-grpc/internals/blogpb"
	"github.com/yuldashev6267/blog-grpc/pkg/models"
)

type server struct{
	blogpb.UnimplementedBlogServiceServer
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.BlogRequest) (*blogpb.BlogResponse, error){
	blog := req.GetBlog()

	data := models.Blog{
		AuthorId: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Comment: blog.GetComment(),
	}

	
}



