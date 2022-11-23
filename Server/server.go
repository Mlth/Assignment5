package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	rep "github.com/Mlth/Assignment5/proto"
	"google.golang.org/grpc"
)

type repServer struct {
	rep.ReplicationServer
}

var highestBid int32
var highestBidderId int32
var highestBidderName string

var AuctionOver bool = false

func main() {
	go func() {
		// auction closing after 1 minut
		fmt.Println("before timeout")
		time.Sleep(30 * time.Second)
		fmt.Println("time is up")
		AuctionOver = true
	}()
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

func (s *repServer) ReceiveBid(ctx context.Context, mess *rep.BidMessage) (*rep.AckMessage, error) {

	if AuctionOver {
		return &rep.AckMessage{BidPlaced: false, AuctionOver: true}, nil
	}
	if mess.Amount > highestBid {
		highestBid = mess.Amount
		highestBidderId = mess.ClientId
		highestBidderName = mess.ClientName
		return &rep.AckMessage{BidPlaced: true, AuctionOver: false}, nil

	}
	return &rep.AckMessage{BidPlaced: false, AuctionOver: false}, nil

}
func (s *repServer) ReturnResult(ctx context.Context, msg *rep.ReqMessage) (*rep.OutcomeMessage, error) {

	return &rep.OutcomeMessage{ClientId: highestBidderId, HighestBid: highestBid, ClientName: highestBidderName, AuctionOver: AuctionOver}, nil
}
