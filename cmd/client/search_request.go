package main

import (
	"github.com/aleg/go-grpc-laptops/client"
	"github.com/aleg/go-grpc-laptops/pb"
	"github.com/aleg/go-grpc-laptops/sample"
)

func testSearchLaptop(client *client.LaptopClient) {
	// Creating a bunch of laptops.
	for i := 0; i < 10; i++ {
		client.CreateLaptop(sample.NewLaptop())
	}

	// Creating a filter and searching using it.
	ram := &pb.Memory{
		Value: 8,
		Unit:  pb.Memory_GIGABYTE,
	}
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam:      ram,
	}
	client.SearchLaptop(filter)
}
