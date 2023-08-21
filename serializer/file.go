package serializer

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

// WriteProtobufToBinaryFile 将 protobuf 写入二进制文件
func WriteProtobufToBinaryFile(message proto.Message, filename string) error {
	// 将protobuf转byte
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	// 写入文件
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}

	return nil
}

// ReadProtobufFromBinaryFile 从二进制中读取到protobuf
func ReadProtobufFromBinaryFile(filename string, message proto.Message) error {
	// 读取protobuf二进制文件
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary data: %w", err)
	}

	// 将bytes序列化转为 protobuf
	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message: %w", err)
	}

	return nil
}

// WriteProtobufToJSONFile 将protobuf写入json文件
func WriteProtobufToJSONFile(message proto.Message, filename string) error {
	//  将 protobuf 转为 JSON string
	data, err := ProtobufToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to JSON: %w", err)
	}

	// 写入文件
	err = os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("cannot write JSON data to file: %w", err)
	}

	return nil
}
