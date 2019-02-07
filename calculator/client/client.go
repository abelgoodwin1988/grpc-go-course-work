package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/abelgoodwin1988/grpc-go-course-work/calculator/calculatorpb"

	"google.golang.org/grpc"
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

	GetPrimeDecomposition(c)

	GetComputeAverage(c)
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