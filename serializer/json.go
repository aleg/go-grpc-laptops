package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON converts protocol buffer message to JSON
// representation as a byte array.
func ProtobufToJSON(message proto.Message) ([]byte, error) {
	marshaler := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		UseProtoNames:   true,
		UseEnumNumbers:  true,
		EmitUnpopulated: true,
	}

	return marshaler.Marshal(message)
}
