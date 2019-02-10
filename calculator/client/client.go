package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/abelgoodwin1988/grpc-go-course-work/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Print("Entered Calculator Client\n")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to Dial CalculatorService Server: %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)

	req := &calculatorpb.SumRequest{
		Sum: []int32{12, 10, 20},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatal("Failed to Sum: %v", err)
	}
	log.Printf("Result from Summing [%v]: %v", req, res)

	// GetPrimeDecomposition(c)

	// GetComputeAverage(c)

	// FindMaximum(c)

	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	// Correct Call
	doErrorCall(c, 10)
	// Error call
	doErrorCall(c, -10)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			// Actual error from gRPC (user error)
			fmt.Println(resErr.Message())
			fmt.Println(resErr.Code())
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number.")
			}
			return
		} else {
			// Other error, frameowkr.
			log.Fatalf("Big Error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Result of Sqrt of number %v: %v\n", n, res.GetNummberRoot())
}

func FindMaximum(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Started FindMaximum")
	// open stream
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Failed to open BiDi Stream to Server", err)
	}
	waitc := make(chan struct{})
	// send
	maxes := []int32{1, 15, 2, 4, 17, 2, 15, 17, 20, 21, 23, 50, 10}
	go func() {
		for _, val := range maxes {
			sendErr := stream.Send(&calculatorpb.FindMaximumRequest{
				Number: val,
			})
			if sendErr != nil {
				log.Fatalf("Failed to send value to server: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("End of responses from server stream")
				break
			}
			if err != nil {
				log.Fatalf("Failed to receive message from server stream: %v", err)
			}
			fmt.Printf("Current Max: %v\n", res.GetMax())
		}
		close(waitc)
	}()

	<-waitc
	// block
}

func GetComputeAverage(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting ComputeAverage Stream")
	reqs := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			AverageSubject: 10,
		},
		&calculatorpb.ComputeAverageRequest{
			AverageSubject: 20,
		},
		&calculatorpb.ComputeAverageRequest{
			AverageSubject: 30,
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Failed to open client stream to server: %v", err)
	}
	for _, req := range reqs {
		log.Printf("ComputeAverage send stream value to server: %v", req.GetAverageSubject())
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to close and receive client stream server response: %v", err)
	}
	log.Printf("Received ComputeAverage Response from Server: %v", res.GetAverage())
}

func GetPrimeDecomposition(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting Prime Number Decomposition")
	req := &calculatorpb.PrimeNumberDecomponsitionRequest{
		Number: 120,
	}
	resStrem, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition: %v", err)
	}
	for {
		msg, err := resStrem.Recv()
		if err == io.EOF {
			log.Print("End of Stream\n")
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Decomposition: ", msg.GetDecomposition())
	}
}
