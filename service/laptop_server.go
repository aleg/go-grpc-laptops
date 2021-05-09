package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/aleg/go-grpc-laptops/pb"
	"github.com/aleg/go-grpc-laptops/stores"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

// 1MB
const maxImageSize = 1 << 20

type ServerStore struct {
	laptop stores.LaptopStore
	image  stores.ImageStore
	rating stores.RatingStore
}
type LaptopServer struct {
	store ServerStore
}

func NewLaptopServer(laptopStore stores.LaptopStore, imageStore stores.ImageStore, ratingStore stores.RatingStore) *LaptopServer {
	st := ServerStore{laptop: laptopStore, image: imageStore, rating: ratingStore}
	return &LaptopServer{store: st}
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	laptopId := laptop.GetId()

	log.Printf("Received a create-laptop request with id %s", laptopId)

	if len(laptopId) > 0 {
		// If the client generated the ID, make sure
		// it's a valid UUID.
		_, err := uuid.Parse(laptopId)
		if err != nil {
			msg := "The laptop ID provided by the client is not a valid UUID: %v"
			return nil, logError(err, codes.InvalidArgument, msg)
		}
	} else {
		// UUID is created by the server!
		id, err := uuid.NewRandom()
		if err != nil {
			msg := "Cannot generate a new UUID for the laptop: %v"
			return nil, logError(err, codes.Internal, msg)
		}
		laptopId = id.String()
		laptop.Id = laptopId
	}

	// TODO: heavy processing
	// time.Sleep(6 * time.Second)

	// Making sure there was no errors before saving to storage.
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	// Save laptop to an in-memory storage.
	err := server.store.laptop.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, stores.ErrorAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, logError(err, code, "Cannot save laptop to the store")
	}
	log.Printf("Saved laptop with id: %s", laptopId)

	response := &pb.CreateLaptopResponse{Id: laptopId}
	return response, nil
}

// SearchLaptop is a server streaming RPC to search a laptop
func (server *LaptopServer) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("Received search-laptop request with filter %v", filter)

	// Send the found laptop to the stream.
	found := func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{Laptop: laptop}
		err := stream.Send(res)
		if err != nil {
			return err
		}

		log.Printf("Sent found laptop with id: %s", laptop.GetId())
		return nil
	}

	err := server.store.laptop.Search(stream.Context(), filter, found)
	if err != nil {
		return logError(err, codes.Internal, "Unexpected error")
	}

	return nil
}

// UploadImage is a client streaming RPC to upload an image in chunks.
// `stream` is the stream of messages (meta data or chunks of the file)
// sent from the client.
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	// First, receive the image matadata.
	req, err := stream.Recv()
	if err != nil {
		return logError(err, codes.Unknown, "Cannot receive image info request")
	}

	laptopId := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("Received upload-image request for laptop %s with image type %s", laptopId, imageType)

	laptop, err := server.store.laptop.Find(laptopId)
	if err != nil {
		return logError(err, codes.Internal, "Cannot find laptop")
	}
	if laptop == nil {
		return logError(nil, codes.InvalidArgument, fmt.Sprintf("Laptop %s doesn't exist", laptopId))
	}

	// Then, start uploading in chunks.
	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		// Checking for errors before receiving more data.
		if err := contextError(stream.Context()); err != nil {
			return err
		}

		log.Print("Waiting to receive more data...")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("No more data")
			break
		}
		if err != nil {
			return logError(err, codes.Unknown, "Cannot receive chunk data")
		}

		chunk := req.GetChunkData() // getting image data from request
		size := len(chunk)          // size of the chunk
		imageSize += size
		log.Printf("Chunk received (%d bytes; current total size: %d)", size, imageSize)

		if imageSize > maxImageSize {
			msg := fmt.Sprintf("Image is too large: %d > %d", imageSize, maxImageSize)
			return logError(nil, codes.InvalidArgument, msg)
		}

		// TODO: writing slowly.
		// time.Sleep(time.Second)

		_, err = imageData.Write(chunk) // writing chunk to the buffer.
		if err != nil {
			return logError(err, codes.Internal, "Cannot write chunk data")
		}
	}

	// Saving image to file.
	imageId, err := server.store.image.Save(laptopId, imageType, imageData)
	if err != nil {
		return logError(err, codes.Internal, "Cannot save image to the store (file)")
	}
	log.Printf("Image saved with id %s, size %d", imageId, imageSize)

	// Generating the response using the image ID just generated.
	res := &pb.UploadImageResponse{Id: imageId, Size: uint32(imageSize)}
	err = stream.SendAndClose(res)
	if err != nil {
		return logError(err, codes.Unknown, "Cannot send response")
	}

	log.Print("Image response sent to the caller")
	return nil
}

// RateLaptop is a bidirectional-streaming RPC that allows clients to rate
// a stream of laptops with a score, and returns a stream of avg scores
// for each of them.
func (server *LaptopServer) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	// The client will send a stream of ratings, hence a
	// while loop is required.
	for {
		// First check that there are no errors.
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		// Receiving the rate-laptop request.
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("No more rates to receive")
			break
		}
		if err != nil {
			return logError(err, codes.Unknown, "Cannot receive request from stream")
		}

		// Extracting the data from the received request.
		laptopId := req.GetLaptopId()
		score := req.GetScore()
		log.Printf("Received a rate-laptop request: id = %s; score = %.2f", laptopId, score)

		// Searhing the laptop by its ID sent with the request.
		found, err := server.store.laptop.Find(laptopId)
		if err != nil {
			msg := fmt.Sprintf("Cannot find laptop with ID %s", laptopId)
			return logError(err, codes.Internal, msg)
		}
		if found == nil {
			msg := fmt.Sprintf("Laptop with ID %s doesn't exist", laptopId)
			return logError(nil, codes.NotFound, msg)
		}

		// Saving the rating to the store.
		rating, err := server.store.rating.Add(laptopId, score)
		if err != nil {
			return logError(err, codes.Internal, "Cannot add rating to the store")
		}

		// Building the response and sending it to the client stream.
		res := &pb.RateLaptopResponse{
			LaptopId:     laptopId,
			RatedCount:   rating.Count,
			AverageScore: rating.Sum / float64(rating.Count),
		}
		err = stream.Send(res)
		if err != nil {
			return logError(err, codes.Unknown, "Cannot send the stream response to the clinet")
		}
		log.Printf("Sent rate-laptop response: id = %s; score = %.2f", laptopId, score)
	}

	return nil
}
