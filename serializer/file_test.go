package serializer_test

import (
	"testing"

	"github.com/aleg/go-grpc-laptops/pb"
	"github.com/aleg/go-grpc-laptops/sample"
	"github.com/aleg/go-grpc-laptops/serializer"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	// var laptop *pb.Laptop
	binaryFile := "../tmp/laptop.bin"

	// Writing binary to file.
	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	// Reading from binary file.
	laptop2 := &pb.Laptop{}
	err = serializer.ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	// Writing JSON to file.
	jsonFile := "../tmp/laptop.json"
	err = serializer.WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)

	// t.Run("WriteProtobufToBinaryFile", func(t *testing.T) {
	//         laptop = sample.NewLaptop()

	//         err := serializer.WriteProtobufToBinaryFile(laptop, binaryFile)
	//         require.NoError(t, err)
	// })

	// t.Run("ReadProtobufFromBinaryFile", func(t *testing.T) {
	//         laptop1 := &pb.Laptop{}

	//         err := serializer.ReadProtobufFromBinaryFile(binaryFile, laptop1)
	//         require.NoError(t, err)
	//         require.True(t, proto.Equal(laptop, laptop1))
	// })
}
