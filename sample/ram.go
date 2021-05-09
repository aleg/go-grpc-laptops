package sample

import (
	"github.com/aleg/go-grpc-laptops/pb"
)

func NewRAM() *pb.Memory {
	memGB := randomInt(4, 64)
	ram := &pb.Memory{
		Value: uint64(memGB),
		Unit:  pb.Memory_GIGABYTE,
	}
	return ram
}
