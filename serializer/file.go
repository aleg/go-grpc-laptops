package serializer

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
)

func writeToFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

// WriteProtobufToJSONFile writes protocol buffer message to JSON file
func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("Cannot marshal proto message to JSON: %w", err)
	}

	// err = ioutil.WriteFile(filename, data, 0644)
	err = writeToFile(filename, data)
	if err != nil {
		return fmt.Errorf("Cannot write JSON data to file: %w", err)
	}

	return nil
}

// WriteProtobufToBinaryFile writes protocol buffer message to binary file
func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	// err = ioutil.WriteFile(filename, data, 0644)
	err = writeToFile(filename, data)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}

	return nil
}

// ReadProtobufFromBinaryFile reads protocol buffer message from binary file
func ReadProtobufFromBinaryFile(filename string, message proto.Message) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message: %w", err)
	}

	return nil
}
