package main

import (
	"github.com/aleg/go-grpc-laptops/client"
	"github.com/aleg/go-grpc-laptops/sample"
)

func testCreateLaptop(client *client.LaptopClient) {
	client.CreateLaptop(sample.NewLaptop())
}
