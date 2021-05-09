package main

import (
	"github.com/aleg/go-grpc-laptops/client"
	"github.com/aleg/go-grpc-laptops/sample"
)

func testUploadImage(client *client.LaptopClient) {
	// First, creating a new laptop via RPC.
	laptop := sample.NewLaptop()
	// imgPath := "tmp/img/400-blows.jpg"
	imgPath := "tmp/img/Eternal Sunshine of the Spotless Mind.jpg"
	client.CreateLaptop(laptop)
	client.UploadImage(laptop.GetId(), imgPath)
}
