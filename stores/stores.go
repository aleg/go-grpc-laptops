package stores

import (
	"bytes"
	"context"
	"errors"

	"github.com/aleg/go-grpc-laptops/pb"
	"github.com/aleg/go-grpc-laptops/users"
)

var ErrorAlreadyExists = errors.New("Record already exists")

// LaptopStore is an interface to store laptop
type LaptopStore interface {
	// Save saves the laptop to the store
	Save(laptop *pb.Laptop) error
	// Find finds a laptop by ID
	Find(id string) (*pb.Laptop, error)
	// Search searches a laptop using the provided filter,
	// and return the results one by one (through a stream)
	// using the callback function `found`.
	Search(ctx context.Context, filter *pb.Filter, found func(*pb.Laptop) error) error
}

type ImageStore interface {
	// Save saves the image to the store (and returns the ID of the saved image).
	Save(laptopId string, imageType string, imageData bytes.Buffer) (string, error)
}

type RatingStore interface {
	// Add adds a new laptop score to the store and returns its rating.
	Add(laptopId string, score float64) (*Rating, error)
}

type Rating struct {
	Count uint32  // number of times the laptop is reated
	Sum   float64 // sum of all rated scores
}

type UserStore interface {
	Save(user *users.User) error
	Find(username string) (*users.User, error)
}
