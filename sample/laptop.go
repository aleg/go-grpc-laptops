package sample

import (
	"github.com/aleg/go-grpc-laptops/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func randomLaptopName(brand string) string {
	switch brand {
	case APPLE:
		return randomStringFromSet("Macbook Air", "Macbook Pro")
	case DELL:
		return randomStringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	case LENOVO:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1", "Thinkpad P53")
	default:
		return "Unknown name"
	}
}

func NewLaptop() *pb.Laptop {
	brand := randomStringFromSet(APPLE, DELL, LENOVO)
	name := randomLaptopName(brand)

	gpu := NewGPU()
	ssd := NewSSD()
	hdd := NewHDD()
	weight := &pb.Laptop_WeightKg{
		WeightKg: randomFloat64(1.0, 3.0),
	}

	// Deprecated: Call the timestamppb.Now function instead.
	laptop := &pb.Laptop{
		Id:          randomID(),
		Brand:       brand,
		Name:        name,
		Cpu:         NewCPU(),
		Ram:         NewRAM(),
		Gpus:        []*pb.GPU{gpu},
		Storages:    []*pb.Storage{ssd, hdd},
		Screen:      NewScreen(),
		Keyboard:    NewKeyboard(),
		Weight:      weight,
		PriceUsd:    randomFloat64(1500, 3500),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		UpdatedAt:   timestamppb.Now(),
	}

	return laptop
}

func RandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
