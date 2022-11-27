package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	rep "github.com/Mlth/Assignment5/proto"
	"google.golang.org/grpc"
)

var name string
var id int64
var totalPorts int64
var reader = bufio.NewReader(os.Stdin)

var clients []rep.ReplicationClient

func main() {
	//Loading id and total amount of ports to connect to
	id, _ = strconv.ParseInt(os.Args[1], 10, 32)
	totalPorts, _ = strconv.ParseInt(os.Args[2], 10, 32)

	//Creating .log-file for logging output from program, while still printing to the command line
	stringy := fmt.Sprintf("%v_client_output.log", id)
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

	//Creating connection to all servers
	for i := 0; i < int(totalPorts); i++ {
		// Create a virtual RPC Client Connection on port 9080 + i
		var conn *grpc.ClientConn
		var port int = 9080 + i
		portStr := strconv.Itoa(port)

		conn, err := grpc.Dial(":"+portStr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		// Defer means: When this function returns, call this method (meaning, one main is done, close connection)
		defer conn.Close()

		//  Create new Client from generated gRPC code from proto
		c := rep.NewReplicationClient(conn)
		clients = append(clients, c)
	}

	log.Print("Write your name: ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	log.Println("Logging: " + name)

	//Starting method for continuously recieving input from user
	takeInput()
}

func takeInput() {
	for {
		log.Println("Place your bid or type \"result\" to see the current highest bid")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		log.Print("Logging: " + input)

		if input == "result" || input == "Result" {
			//Calling method to recieve result from all servers
			result, _ := getResultFromAll(&rep.ReqMessage{})

			//Displaying information to user based on whether the auction is over or not
			if !result.AuctionOver {
				log.Printf("The current highest bid is: %d placed by %s (%d)\n", result.HighestBid, result.ClientName, result.ClientId)
				continue
			} else {
				log.Printf("The auction is over! Sold to %s (%d) for %d\n", result.ClientName, result.ClientId, result.HighestBid)
				break
			}
		}

		//If we reach this point in the code, we know that the user did not type 'result'. Therefore, we convert the input to
		//an int. If that is not possible, we know that the input is not valid, since it would either be 'result' or a number to bid
		intInput, err := strconv.Atoi(input)
		if err != nil {
			log.Println("Faulty input, please try again")
			continue
		}

		//Calling method for sending bid to all servers
		ack, _ := sendBidToAll(&rep.BidMessage{ClientId: int32(id), Amount: int32(intInput), ClientName: name})

		//Printing result to user, based on whether the bid was succesful
		if ack.BidPlaced {
			log.Println("Your bid has been placed!")
		} else if !ack.AuctionOver {
			log.Println("Your bid was too low!")
		} else {
			log.Println("Bid not placed, the auction is over.")
			result, _ := getResultFromAll(&rep.ReqMessage{})
			log.Printf("The auction is over! Sold to %s (%d) for %d\n", result.ClientName, result.ClientId, result.HighestBid)
			break

		}
	}
}

func sendBidToAll(message *rep.BidMessage) (*rep.AckMessage, error) {
	var ack *rep.AckMessage
	for _, c := range clients {
		tempAck, err := c.ReceiveBid(context.Background(), message)
		if err == nil {
			ack = tempAck
		}
	}
	return ack, nil
}

func getResultFromAll(message *rep.ReqMessage) (*rep.OutcomeMessage, error) {
	var result *rep.OutcomeMessage
	for _, c := range clients {
		tempResult, err := c.ReturnResult(context.Background(), message)
		if err == nil {
			result = tempResult
		}
	}
	return result, nil
}
