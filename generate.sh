protoc ./pkg/internals/blogpb/blog.proto --go_out=.

protoc ./pkg/internals/blogpb/blog.proto --go-grpc_out=.

run server
  go run .\cmd\server\server.go