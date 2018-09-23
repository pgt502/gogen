package auth

import (
	"context"
	"net/http"

	"google.golang.org/grpc"
)

type AuthMiddleware interface {
	Log(h http.Handler) http.Handler
	GetIdentity(h http.Handler) http.Handler
	GrpcGetIdentity(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
	Copy(string, ...interface{}) (int, []interface{})
}
