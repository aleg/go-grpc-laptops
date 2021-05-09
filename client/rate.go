package client

import (
	"fmt"
	"io"
	"log"

	"github.com/aleg/go-grpc-laptops/pb"
)

func receiveResponses(stream pb.LaptopService_RateLaptopClient, waitResponse chan<- error) {
	// Need a while loop as multiple responses can
	// be sent back from the server.
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Print("No more responses to be received.")
			waitResponse <- nil
			return
		}
		if err != nil {
			waitResponse <- fmt.Errorf("Cannot receive stream response: %v", err)
			return
		}

		log.Print("Received response: ", res)
	}
}

func sendRequests(stream pb.LaptopService_RateLaptopClient, laptopIds []string, scores []float64) error {
	for i, laptopId := range laptopIds {
		req := &pb.RateLaptopRequest{
			LaptopId: laptopId,
			Score:    scores[i],
		}
		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("Cannot send request: %v - %v", err, stream.RecvMsg(nil))
		}

		log.Print("Sent rate-laptop request: ", req)
	}

	return nil
}
