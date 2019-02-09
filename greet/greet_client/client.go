package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/abelgoodwin1988/grpc-go-course-work/greet/greetpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doBidirectionalStreaming(c)
}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Start bidirectional streaming RPC")
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Failed to open stream to server: %v", err)
		return
	}

	requests := _request

	waitc := make(chan struct{})
	// We send messages to the server
	go func() {
		// Function to send messages
		for _, req := range requests {
			fmt.Printf("Sending Message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// We receive messages from the server
	go func() {
		// Function to receive messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Err while receiving from server: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block unti leverything is done
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Client Streaming RPC")
	reqStream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while opening LongGreeting Stream: %v", err)
	}
	names := [10]string{"Abel", "Joe", "Dave", "Steve", "Lauritz", "Lauren", "Loren", "Ralph", "Netty", "Ted"}
	for i := 0; i < len(names); i++ {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: names[i],
				LastName:  "Goodwin",
			},
		}
		if err := reqStream.Send(req); err != nil {
			log.Fatalf("Error while sending stream: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
		log.Printf("Req sent in stream: %v", i)
	}
	res, err := reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to either close or receive stream response: %v", err)
	}
	log.Printf("Long Greet Response: %v", res.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Server Streaming RPC")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Abel",
			LastName:  "Goodwin",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// We're reached the enf othe stream
			log.Printf("End of Stream\n")
			break
		}
		if err != nil {
			log.Fatal("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Unary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Hamp",
			LastName:  "Goodwin",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatal("Error while calling Greeting RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

var _request = []*greetpb.GreetEveryoneRequest{
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Abel",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Stephen",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Kyle",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Daniel",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Brian",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Brett",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Lee",
		},
	},
	&greetpb.GreetEveryoneRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Peyton",
		},
	},
}
