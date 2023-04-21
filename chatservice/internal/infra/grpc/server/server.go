package server

import (
	"github.com/ahugofreire/chatservice/internal/infra/grpc/pb"
	"github.com/ahugofreire/chatservice/internal/infra/grpc/service"
	"github.com/ahugofreire/chatservice/internal/usecase/chatcompletionstream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

type GRPCServer struct {
	ChatCompletionStreamUseCase chatcompletionstream.ChatCompletionUserCase
	ChatConfigStream            chatcompletionstream.ChatCompletionConfigDTO
	ChatService                 service.ChatService
	Port                        string
	AuthToken                   string
	StreamChannel               chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewGRPCServer(
	chatCompletionStreamUseCase chatcompletionstream.ChatCompletionUserCase,
	chatConfigStream chatcompletionstream.ChatCompletionConfigDTO,
	port string,
	authToken string,
	StreamChannel chan chatcompletionstream.ChatCompletionOutputDTO,
) *GRPCServer {
	chatService := service.NewChatService(chatCompletionStreamUseCase, chatConfigStream, StreamChannel)
	return &GRPCServer{
		ChatCompletionStreamUseCase: chatCompletionStreamUseCase,
		ChatConfigStream:            chatConfigStream,
		ChatService:                 *chatService,
		Port:                        port,
		AuthToken:                   authToken,
		StreamChannel:               StreamChannel,
	}
}

// Interface Configurada para GRPC com stream de dados.
func (g *GRPCServer) AuthInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	if token[0] != g.AuthToken {
		return status.Error(codes.Unauthenticated, "authorization token is invalid")
	}

	return handler(srv, ss)
}

func (g *GRPCServer) Start() {
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(g.AuthInterceptor),
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterChatServiceServer(grpcServer, &g.ChatService)

	lis, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		panic(err.Error())
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err.Error())
	}
}
