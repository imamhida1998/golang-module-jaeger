// Package grpc — adapter TIPIS di atas tracing.Inject/Extract, khusus untuk gRPC
// unary calls. Opsional: pakai ini kalau service kamu gRPC. Tidak menambah dependency
// baru selain google.golang.org/grpc (tidak butuh otelgrpc).
package grpc

import (
	"context"

	"github.com/funxdofficial/golang-module-jaeger/tracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor men-extract trace context dari metadata gRPC masuk dan
// menaruhnya di ctx yang diteruskan ke handler.
//
//	grpc.NewServer(grpc.UnaryInterceptor(grpctransport.UnaryServerInterceptor()))
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		headers := map[string]string{}
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for k, v := range md {
				if len(v) > 0 {
					headers[k] = v[0]
				}
			}
		}
		ctx = tracing.Extract(ctx, headers)
		return handler(ctx, req)
	}
}

// UnaryClientInterceptor menyalin trace context dari ctx ke metadata gRPC keluar,
// supaya server penerima bisa extract dan lanjutkan trace yang sama.
//
//	conn, _ := grpc.Dial(addr, grpc.WithUnaryInterceptor(grpctransport.UnaryClientInterceptor()))
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		headers := tracing.Inject(ctx)
		md := metadata.New(headers)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
