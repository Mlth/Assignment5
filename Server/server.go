package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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
var isFirstBid bool = true
var AuctionOver bool = false

func main() {

	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	//arg2, _ := strconv.ParseInt(os.Args[2], 10, 32)
	ownPort := int32(arg1) + 9080
	ownPortStr := strconv.Itoa(int(ownPort))

	list, err := net.Listen("tcp", ":"+ownPortStr)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", ownPortStr, err)
	}
	grpcServer := grpc.NewServer()
	rep.RegisterReplicationServer(grpcServer, &repServer{})
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}

func (s *repServer) ReceiveBid(ctx context.Context, mess *rep.BidMessage) (*rep.AckMessage, error) {
	if isFirstBid {
		go func() {
			// auction closing after 1 minut
			fmt.Println("before timeout")
			time.Sleep(720 * time.Second)
			fmt.Println("time is up")
			AuctionOver = true
		}()
		isFirstBid = false
	}
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
