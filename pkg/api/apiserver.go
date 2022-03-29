package api

import "github.com/gin-gonic/gin"

type apiServer struct {
	addr string
	server *gin.Engine
}

func New(address string) *apiServer {
	as := &apiServer{
		addr: address,
		server: gin.Default(),
	}
	return as
}


func (a *apiServer) RegisterGinServer() error{
	return a.server.Run(a.addr)
}
