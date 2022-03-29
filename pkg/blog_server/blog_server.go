package blog_server

import (
	"context"
	"github.com/yuldashev6267/blog-grpc/internals/blogpb"
	"github.com/yuldashev6267/blog-grpc/pkg/models"
	"github.com/yuldashev6267/blog-grpc/pkg/api/controllers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{
	blogpb.UnimplementedBlogServiceServer
	blog controllers.BlogService
}

type

func (s *server) CreateBlog(ctx context.Context, req *blogpb.BlogRequest) (*blogpb.BlogResponse, error){
	blog := req.GetBlog()

	data := models.Blog{
		AuthorId: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Comment: blog.GetComment(),
	}
	res, err := s.blog.InsertBlog(data)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"Errror inserting document data base %v", err,
			)
	}

	return &blogpb.BlogResponse{
		Blog: &blogpb.Blog{
			Id:res.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Comment: blog.GetComment(),
		},
	}, nil
}


