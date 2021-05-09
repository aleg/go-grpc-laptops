package stores

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/aleg/go-grpc-laptops/pb"
	"github.com/jinzhu/copier"
)

type InMemoryLaptopStore struct {
	// There will be concurrent requests to write
	// a laptop to memory, so a mutex is needed.
	m sync.RWMutex // multiple readers, one writer
	// key: laptop ID; value: laptop object.
	data map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Implements the `Save` method of the `LaptopStore` interface.
func (st *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	st.m.Lock() // locking for writing. Also reads are blocked.
	defer st.m.Unlock()

	if _, found := st.data[laptop.GetId()]; found {
		return ErrorAlreadyExists
	}

	// Deep copy of laptop before copying it into memory.
	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	// Saving in the memory st.
	st.data[other.GetId()] = other

	return nil
}

// Implements the `Find` method of the `LaptopStore` interface.
func (st *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	st.m.RLock()
	defer st.m.RUnlock()

	laptop, found := st.data[id]
	if !found {
		return nil, nil
	}

	return deepCopy(laptop)
}

func (st *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(*pb.Laptop) error) error {
	st.m.RLock()
	defer st.m.RUnlock()

	// Going through each laptop in the store
	// and searching with `filter`. When something
	// matches, the callback function `found` is
	// called.
	for _, laptop := range st.data {
		// TODO: heavy processing
		// time.Sleep(time.Second)
		log.Print("check laptop id: ", laptop.GetId())

		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("Search stopped: context is cancelled")
			return errors.New("Context is cancelled")
		}

		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}

	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}

	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}

	return true
}

func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3 // 8 = 2^3
	case pb.Memory_KILOBYTE:
		return value << 13 // 1024 * 8 = 2^10 * 2^3 = 2^13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}

	err := copier.Copy(other, laptop) // to, from
	if err != nil {
		return nil, fmt.Errorf("Cannot copy laptop data: %w", err)
	}

	return other, nil
}
