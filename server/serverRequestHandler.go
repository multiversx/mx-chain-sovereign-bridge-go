package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type serverRequestsHandler struct {
	ginHandler  *gin.Engine
	grpcHandler *grpc.Server
}

func NewServerHandler(ginHandler *gin.Engine, grpcHandler *grpc.Server) (*serverRequestsHandler, error) {
	return &serverRequestsHandler{
		ginHandler:  ginHandler,
		grpcHandler: grpcHandler,
	}, nil
}

func (h *serverRequestsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/grpc") {
		log.Trace("server handling grpc request")
		h.grpcHandler.ServeHTTP(w, req)
		return
	}

	log.Trace("server handling http request")
	h.ginHandler.ServeHTTP(w, req)
}
