package handler

import (
	"log"
	notspb "wch/gen/nots/v1"
	"wch/services/nots/internal/streams"
)

type GRPCServer struct {
	notspb.UnimplementedNotificationServiceServer
	hub *streams.StreamHub
}

func NewGRPCServer(hub *streams.StreamHub) *GRPCServer {
	return &GRPCServer{
		hub: hub,
	}
}

func (s *GRPCServer) Subscribe(
	req *notspb.SubscribeRequest,
	stream notspb.NotificationService_SubscribeServer,
) error {
	log.Println("GRPCServer.Subscribe: starting... ", req)
	userID := req.GetUserId()
	log.Println("User subscribed: ", userID)

	s.hub.Add(userID, stream)
	defer s.hub.Remove(userID)

	// keep sonnection open
	<-stream.Context().Done()

	log.Println("User disconnected: ", userID)
	return nil
}
