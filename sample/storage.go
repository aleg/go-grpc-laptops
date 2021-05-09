package sample

import (
	"github.com/aleg/go-grpc-laptops/pb"
)

func NewSSD() *pb.Storage {
	memGB := randomInt(128, 1024)

	memory := &pb.Memory{
		Value: uint64(memGB),
		Unit:  pb.Memory_GIGABYTE,
	}
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: memory,
	}

	return ssd
}

func NewHDD() *pb.Storage {
	memTB := randomInt(1, 6)

	memory := &pb.Memory{
		Value: uint64(memTB),
		Unit:  pb.Memory_TERABYTE,
	}
	hdd := &pb.Storage{
		Driver: pb.Storage_HDD,
		Memory: memory,
	}

	return hdd
}
