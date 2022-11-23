package main

import (
	"context"
	"fmt"
	"io"
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
	log.Println("Starting server on port " + ownPortStr)
	//Creating .log-file for logging output from program, while still printing to the command line
	stringy := fmt.Sprintf("%v_server_output.log", ownPort)
	err := os.Remove(stringy)
	if err != nil {
		log.Println("No previous log file found")
	}
	f, err := os.OpenFile(stringy, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	mw := io.MultiWriter(os.Stdout, f)
	if err != nil {
		fmt.Println("Log does not work")
	}
	defer f.Close()
	log.SetOutput(mw)

	list, err := net.Listen("tcp", ":"+ownPortStr)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", ownPortStr, err)
	}
	grpcServer := grpc.NewServer()
	rep.RegisterReplicationServer(grpcServer, &repServer{})
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}

func (s *repServer) ReceiveBid(ctx context.Context, mess *rep.BidMessage) (*rep.AckMessage, error) {
	if isFirstBid {
		log.Println("First bid received, starting auction timer for 1 minute")
		go func() {
			// Auction closes after 1 minute.
			time.Sleep(60 * time.Second)
			AuctionOver = true
		}()
		isFirstBid = false
	}
	if AuctionOver {
		log.Println("Auction is over, no more bids accepted")
		return &rep.AckMessage{BidPlaced: false, AuctionOver: true}, nil
	}
	if mess.Amount > highestBid {
		log.Println("New highest bid received, updating highest bid")
		highestBid = mess.Amount
		highestBidderId = mess.ClientId
		highestBidderName = mess.ClientName
		log.Printf("New highest bid is %v from client %s (%v)", highestBid, highestBidderName, highestBidderId)
		return &rep.AckMessage{BidPlaced: true, AuctionOver: false}, nil

	}
	log.Println("Bid was not highest, not updating highest bid")
	return &rep.AckMessage{BidPlaced: false, AuctionOver: false}, nil

}
func (s *repServer) ReturnResult(ctx context.Context, msg *rep.ReqMessage) (*rep.OutcomeMessage, error) {
	log.Println("Client requested result of auction")
	log.Printf("Returning result to client: highest bid is %v from client %s (%v)", highestBid, highestBidderName, highestBidderId)
	return &rep.OutcomeMessage{ClientId: highestBidderId, HighestBid: highestBid, ClientName: highestBidderName, AuctionOver: AuctionOver}, nil
}
