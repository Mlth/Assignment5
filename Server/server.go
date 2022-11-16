package main

import (
	"context"
	"log"
	"net"

	rep "github.com/Mlth/Assignment5/proto"
	"google.golang.org/grpc"
)

type repServer struct {
	rep.ReplicationServer
}

var highestBid int32 = 0
var highestBidderId int32 = 0

func main() {

	// Vi skal lave det sådan, at man ikke har hardcoded den til port 9080
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()
	rep.RegisterReplicationServer(grpcServer, &repServer{})
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}

func (s *repServer) Bid(ctx context.Context, mess *rep.BidMessage) (*rep.AckMessage, error) {
	if mess.Amount > highestBid {
		highestBid = mess.Amount
		highestBidderId = mess.ClientId
		return &rep.AckMessage{Status: "Success"}, nil

	}
	return &rep.AckMessage{Status: "Failure"}, nil

}
func (s *repServer) Result(ctx context.Context, msg *rep.ReqMessage) (*rep.OutcomeMessage, error) {

	return &rep.OutcomeMessage{ClientId: highestBidderId, HighestBid: highestBid}, nil
}
