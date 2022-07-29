package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AdityaPusalkarSwiggy/gRPC_Calculator/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	cc, err := grpc.Dial("localhost:50505", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calcpb.NewCalculatorServiceClient(cc)

	// Unary function - Sum
	Sum(c)

	// Server-Side Streaming function - Prime
	Prime(c)

	// Client-Side Streaming function - Average
	Average(c)

	// Bi-Directional Streaming function - Max
	Max(c)
}

func Greet(c calcpb.CalculatorServiceClient) {

	fmt.Println("Starting to do a unary GRPC....")

	req := calcpb.GreetRequest{
		Number: &calcpb.Number{
			FirstName: "John",
			LastName:  "Doe",
		},
	}

	resp, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling greet grpc unary call: %v", err)
	}

	log.Printf("Response from Greet Unary Call : %v", resp.Result)

}

func GreetManyTimes(c calcpb.CalculatorServiceClient) {
	fmt.Println("Staring ServerSide GRPC streaming ....")

	req := calcpb.GreetManyTimesRequest{

		Number: &calcpb.Number{
			FirstName: "Michael",
			LastName:  "Jackson",
		},
	}

	respStream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes server-side streaming grpc : %v", err)
	}

	for {
		msg, err := respStream.Recv()
		if err == io.EOF {
			//we have reached to the end of the file
			break
		}

		if err != nil {
			log.Fatalf("error while receving server stream : %v", err)
		}

		fmt.Println("Response From GreetManyTimes Server : ", msg.Result)
	}
}

func Average(c calcpb.CalculatorServiceClient) {

	fmt.Println("Starting Client Side Streaming over GRPC ....")

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client-side streaming : %v", err)
	}

	requests := []*calcpb.AverageRequest{
		&calcpb.AverageRequest{
			Number: &calcpb.Number{
				FirstName: "James",
				LastName:  "Bond",
			},
		},
		&calcpb.AverageRequest{
			Number: &calcpb.Number{
				Num: 5,
			},
		},
		&calcpb.AverageRequest{
			Number: &calcpb.Number{
				FirstName: "Elon",
				LastName:  "Musk",
			},
		},
		&calcpb.AverageRequest{
			Number: &calcpb.Number{
				FirstName: "Sam",
				LastName:  "Rutherford",
			},
		},
		&calcpb.AverageRequest{
			Number: &calcpb.Number{
				FirstName: "Eoin",
				LastName:  "Morgan",
			},
		},
	}

	for _, req := range requests {
		fmt.Println("\nSending Request.... : ", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from server : %v", err)
	}
	fmt.Println("\n****Response From Server : ", resp.GetResult())
}

func GreetEveryone(c calcpb.CalculatorServiceClient) {
	fmt.Println("Starting Bi-directional stream by calling GreetEveryone over GRPC......")

	requests := []*calcpb.GreetEveryoneRequest{
		&calcpb.GreetEveryoneRequest{
			Number: &calcpb.Number{
				FirstName: "James",
				LastName:  "Bond",
			},
		},
		&calcpb.GreetEveryoneRequest{
			Number: &calcpb.Number{
				FirstName: "Harry",
				LastName:  "Mitchel",
			},
		},
		&calcpb.GreetEveryoneRequest{
			Number: &calcpb.Number{
				FirstName: "Tim",
				LastName:  "Cook",
			},
		},
		&calcpb.GreetEveryoneRequest{
			Number: &calcpb.Number{
				FirstName: "Sam",
				LastName:  "Rutherford",
			},
		},
		&calcpb.GreetEveryoneRequest{
			Number: &calcpb.Number{
				FirstName: "Eoin",
				LastName:  "Morgan",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error occured while performing client side streaming : %v", err)
	}

	//wait channel to block receiver
	waitchan := make(chan struct{})

	go func(requests []*calcpb.GreetEveryoneRequest) {
		for _, req := range requests {

			fmt.Println("\nSending Request..... : ", req.Number)
			err := stream.Send(req)
			if err != nil {
				log.Fatalf("error while sending request to GreetEveryone service : %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}(requests)

	go func() {
		for {

			resp, err := stream.Recv()
			if err == io.EOF {
				close(waitchan)
				return
			}

			if err != nil {
				log.Fatalf("error receiving response from server : %v", err)
			}

			fmt.Printf("\nResponse From Server : %v", resp.GetResult())
		}
	}()

	//block until everything is finished
	<-waitchan
}
