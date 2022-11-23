package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	rep "github.com/Mlth/Assignment5/proto"
	"google.golang.org/grpc"
)

var name string
var id int32
var reader = bufio.NewReader(os.Stdin)

func main() {
	// Create a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
	defer conn.Close()

	fmt.Print("Write your name: ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	//  Create new Client from generated gRPC code from proto
	c := rep.NewReplicationClient(conn)

	takeInput(c)

}

func takeInput(c rep.ReplicationClient) {
	for {
		fmt.Println("place your bid or type \"result\" to see the current highest bid")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "result" {
			result, _ := c.ReturnResult(context.Background(), &rep.ReqMessage{})

			if !result.AuctionOver {
				fmt.Printf("The current highest bid is: %d placed by %s (%d)\n", result.HighestBid, result.ClientName, result.ClientId)
				continue
			} else {
				fmt.Printf("The auction is over! Sold to %s (%d) for %d\n", result.ClientName, result.ClientId, result.HighestBid)
				break
			}
		}

		intInput, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Faulty input, please try again")
			continue
		}

		ack, _ := c.ReceiveBid(context.Background(), &rep.BidMessage{ClientId: id, Amount: int32(intInput), ClientName: name})

		if ack.BidPlaced {
			fmt.Println("Your bid has been placed!")
		} else if(!ack.AuctionOver) {
			fmt.Println("Your bid was to low!")
		}else{
			fmt.Println("Bid not placed, the auction is over.")
		}
	
	}

}
