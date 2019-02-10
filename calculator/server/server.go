package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"

	"github.com/abelgoodwin1988/grpc-go-course-work/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum Function Invoked with: %v\n", in)
	ints := in.GetSum()
	var sum int32
	for _, val := range ints {
		sum += val
	}
	res := &calculatorpb.SumResponse{
		Sum: sum,
	}
	return res, nil
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NummberRoot: math.Sqrt(float64(number)),
	}, nil
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			stream.Send(&calculatorpb.FindMaximumResponse{
				Max: max,
			})
			return nil
		}
		if err != nil {
			log.Fatalf("Failure to receive request from client stream: %v", err)
			return err
		}
		if max < req.GetNumber() {
			max = req.GetNumber()
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Max: max,
			})
			if sendErr != nil {
				log.Fatalf("Failed to send value to client via stream: $v", err)
				return err
			}
		}
	}
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecomponsitionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function invoked with: %v\n", req)
	number := req.GetNumber()
	k := int32(2)
	for {
		if number < 2 {
			break
		}
		if number%k == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecomponsitionResponse{Decomposition: k})
			number = number / k
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage was invoked with a stream")
	i := int32(1)
	sum := int32(0)
	average := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("ComputeAverage reached end of stream\n")
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Failed to receive value from stream: %\n", err)
		}
		sum += req.GetAverageSubject()
		average = sum / i
		log.Printf("New Average: %v\n from adding %v", average, req.GetAverageSubject())
		i++
	}
}

func main() {
	fmt.Printf("CalculatorServiceServer Started\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Listener Failed: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Server CalculatorService Server: %v", err)
	}
}
