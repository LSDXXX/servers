package grpc

import (
	"context"

	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/servercontext"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor  grpc client interceptor
//  @param ctx
//  @param method
//  @param req
//  @param reply
//  @param cc
//  @param invoker
//  @param opts
//  @return error
func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ctx = servercontext.ContextToGrpc(ctx)
	return invoker(ctx, method, req, reply, cc, opts...)
}

// UnaryServerInterceptor grpc server interceptor
//  @param ctx
//  @param req
//  @param info
//  @param handler
//  @return resp
//  @return err
func UnaryServerInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	ctx = servercontext.ExtractFromGrpc(ctx)
	resp, err = handler(ctx, req)
	log.WithContext(ctx).Infof("grpc request info: %+v, req: %+v, resp: %+v", info, req, resp)
	return
}
