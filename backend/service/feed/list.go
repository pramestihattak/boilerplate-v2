package feed

import (
	pb "boilerplate-v2/gen/feed"
	"context"
	"log"
)

func (s *FeedService) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	log.Printf("Received: %v", req.GetLimit())
	return &pb.ListResponse{Message: "List feed."}, nil
}
