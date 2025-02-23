package grpc

import (
	"context"
	"github.com/google/uuid"
	"go-boilerplate/internal/models"
	"go-boilerplate/internal/service"
	pb "go-boilerplate/proto/message/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageServer struct {
	pb.UnimplementedMessageServiceServer
	messageService *service.MessageService
}

func NewMessageServer(messageService *service.MessageService) *MessageServer {
	return &MessageServer{
		messageService: messageService,
	}
}

func (s *MessageServer) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest) (*pb.MessageResponse, error) {
	message := &models.Message{
		Content: req.Content,
	}

	if err := s.messageService.CreateMessage(ctx, message); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	}

	return &pb.MessageResponse{
		Id:        message.ID.String(),
		Content:   message.Content,
		CreatedAt: timestamppb.New(message.CreatedAt),
	}, nil
}

func (s *MessageServer) GetMessage(ctx context.Context, req *pb.GetMessageRequest) (*pb.MessageResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid message ID: %v", err)
	}

	message, err := s.messageService.GetMessage(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get message: %v", err)
	}

	if message == nil {
		return nil, status.Error(codes.NotFound, "message not found")
	}

	return &pb.MessageResponse{
		Id:        message.ID.String(),
		Content:   message.Content,
		CreatedAt: timestamppb.New(message.CreatedAt),
	}, nil
}

func (s *MessageServer) StreamMessages(empty *emptypb.Empty, stream pb.MessageService_StreamMessagesServer) error {
	messages, err := s.messageService.ListMessages(stream.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to list messages: %v", err)
	}

	for _, msg := range messages {
		if err := stream.Send(&pb.MessageResponse{
			Id:        msg.ID.String(),
			Content:   msg.Content,
			CreatedAt: timestamppb.New(msg.CreatedAt),
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to send message: %v", err)
		}
	}

	return nil
}
