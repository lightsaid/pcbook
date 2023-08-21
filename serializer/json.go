package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON 将protobuf转string
func ProtobufToJSON(message proto.Message) (string, error) {
	// 配置 protobuf 转 json string 格式
	marshaler := protojson.MarshalOptions{
		Indent:          " ",  // 缩进，格式化输出
		UseProtoNames:   true, // 使用protobuf字段名？
		EmitUnpopulated: true, // 零值？
	}

	byte, err := marshaler.Marshal(message)

	return string(byte), err
}
