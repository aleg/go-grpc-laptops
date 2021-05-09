package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aleg/go-grpc-laptops/client"
	"github.com/aleg/go-grpc-laptops/sample"
)

func testRateLaptop(client *client.LaptopClient) {
	n := 3 // # of laptops to rate.
	laptopIds := make([]string, 0, n)

	// First creating the laptops...
	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIds = append(laptopIds, laptop.GetId())
		client.CreateLaptop(laptop)
	}

	// ...then adding the scores.
	scores := make([]float64, n, n)
	for {
		fmt.Print("rate laptop? (y/n) ")
		var answer string
		fmt.Scan(&answer)

		if strings.ToLower(answer) != "y" {
			break
		}

		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}

		err := client.RateLaptop(laptopIds, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}
