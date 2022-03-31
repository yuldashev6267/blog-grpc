package blog_server

import (
	"context"
	"github.com/yuldashev6267/blog-grpc/internals/blogpb"
	"github.com/yuldashev6267/blog-grpc/pkg/api/controllers"
	"github.com/yuldashev6267/blog-grpc/pkg/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlogController struct {
	blog controllers.BlogService
}

type BlogService interface {
	CreateBlog(ctx context.Context,req *blogpb.BlogResponse) (*blogpb.BlogResponse, error)
}

func (s *BlogController) CreateBlog(ctx context.Context, req *blogpb.BlogRequest) (*blogpb.BlogResponse, error){
	blog := req.GetBlog()

	data := models.Blog{
		AuthorId: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Comment: blog.GetComment(),
	}

	res, err := s.blog.InsertBlog(&data)

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
//
//func New(blog controllers.BlogService)BlogService {
//	return &server{blog: blog}
//}
