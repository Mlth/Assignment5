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

var clients []rep.ReplicationClient

func main() {
	for i:=0 ; i<3 ; i++{
		// Create a virtual RPC Client Connection on port  9080 WithInsecure (because  of http)
		var conn *grpc.ClientConn
		var port int = 9080 + i
		portStr:= strconv.Itoa(port)
		
		conn, err := grpc.Dial(":" + portStr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		// Defer means: When this function returns, call this method (meaing, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := rep.NewReplicationClient(conn)
		clients = append(clients, c)
	}
	
	fmt.Print("Write your name: ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	takeInput()
}

func takeInput() {
	for {
		fmt.Println("place your bid or type \"result\" to see the current highest bid")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "result" {
			result, _ := getResultFromAll(&rep.ReqMessage{})

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

		ack, _ := sendBidToAll(&rep.BidMessage{ClientId: id, Amount: int32(intInput), ClientName: name})

		if ack.BidPlaced {
			fmt.Println("Your bid has been placed!")
		} else if(!ack.AuctionOver) {
			fmt.Println("Your bid was too low!")
		}else{
			fmt.Println("Bid not placed, the auction is over.")
		}
	}
}

func sendBidToAll(message *rep.BidMessage) (*rep.AckMessage, error){
	var ack *rep.AckMessage
	for _,c:= range clients{
		tempAck,err := c.ReceiveBid(context.Background(), message)
		if(err == nil){
			ack = tempAck
		}
	}
	return ack, nil
}

func getResultFromAll(message *rep.ReqMessage) (*rep.OutcomeMessage, error){
	var result *rep.OutcomeMessage
	for _,c:= range clients{
		tempResult,err := c.ReturnResult(context.Background(), message)
		if(err == nil){
			result = tempResult
		}
	}
	return result, nil
}

