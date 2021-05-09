package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aleg/go-grpc-laptops/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	service := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{service}
}

func (client *LaptopClient) CreateLaptop(laptop *pb.Laptop) {
	// laptop.Id = ""
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	log.Printf("Going to create laptop %s", laptop.GetId())

	// Setting the timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := client.service.CreateLaptop(ctx, req)
	alreadyExists := false
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			// Not a big deal.
			log.Print("Laptop already exists")
			alreadyExists = true
		} else {
			log.Fatal("Cannot create laptop: ", err)
		}
	}

	if !alreadyExists {
		log.Printf("Created laptop with ID %s", res.GetId())
	}
}

func (client *LaptopClient) SearchLaptop(filter *pb.Filter) {
	log.Print("Going to search laptop: ", filter)

	// Setting the timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := client.service.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("Cannot search laptop: ", err)
	}

	// Keep receiving the found laptops from
	// the stream, until no results.
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return // gracefully return if no more results
		}
		if err != nil {
			log.Fatal("Cannot receive response from search: ", err)
		}

		laptop := res.GetLaptop()
		log.Print("- Found laptop with ID ", laptop.GetId())
		log.Print("\tbrand: ", laptop.GetBrand())
		log.Print("\tname: ", laptop.GetName())
		log.Print("\tCPU: cores ", laptop.GetCpu().GetNumberCores())
		log.Print("\tCPU min ghz: ", laptop.GetCpu().GetMinGhz())
		log.Print("\tRAM: ", laptop.GetRam().GetValue(), laptop.GetRam().GetUnit())
		log.Print("\tprice: ", laptop.GetPriceUsd())
	}
}

func (client *LaptopClient) UploadImage(laptopId string, imagePath string) {
	// Opening the image file.
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Cannot open file %s", imagePath)
	}
	defer file.Close()

	// Creating the context.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Creating the stream.
	stream, err := client.service.UploadImage(ctx)
	if err != nil {
		// Also printing the "real" error
		log.Fatal("Cannot upload image", err, stream.RecvMsg(nil))
	}

	// First sending the meta-data request.
	metaReq := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.UploadImageRequest_ImageInfo{
				LaptopId:  laptopId,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}
	err = stream.Send(metaReq)
	if err != nil {
		// Also printing the "real" error
		log.Fatal("Cannot send meta request to upload image", err, stream.RecvMsg(nil))
	}
	log.Print("Meta request for uploading the image has been sent to the server")

	// Sending the image in chunks.
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1<<10)
	log.Print("Going to send the image in chunks")
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			log.Print("Finished reading the image")
			break
		}
		if err != nil {
			log.Fatal("Cannot read chunk to buffer", err)
		}

		chunkReq := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(chunkReq)
		if err != nil {
			// Also printing the "real" error
			log.Fatalf("Cannot send image chunk of size %d: %v (%v)", n, err, stream.RecvMsg(nil))
		}
		log.Print("Sent image chunk of size", n)
	}

	// Waiting for the server to answer.
	log.Print("Waiting for the server to finish uploading the image...")
	res, err := stream.CloseAndRecv()
	if err != nil {
		// Also printing the "real" error
		log.Fatal("Cannot receive response from the server", err, stream.RecvMsg(nil))
	}
	log.Printf("Image uploaded with ID %s and size %d", res.GetId(), res.GetSize())
}

func (client *LaptopClient) RateLaptop(laptopIds []string, scores []float64) error {
	log.Printf("Going to rate %d laptops: ", len(laptopIds))

	// Setting the timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Starting the client stream to send the rates.
	stream, err := client.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("Cannot rate laptops: %v", err)
	}

	// Receiving the responses from the server
	// (in a go routine to receive multiple responses as the
	// client sends rating requests).
	waitResponse := make(chan error)
	go receiveResponses(stream, waitResponse)

	// Sending the requests.
	err = sendRequests(stream, laptopIds, scores)
	if err != nil {
		return err
	}

	// Closing the stream after sending the requests.
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("Cannot close send (client stream): %v", err)
	}

	return <-waitResponse
}
