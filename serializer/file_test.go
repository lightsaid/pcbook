package serializer

import (
	"fmt"
	"os"
	"testing"

	"github.com/lightsaid/pcbook/pb"
	"github.com/lightsaid/pcbook/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop := sample.NewLaptop()

	// 写入二进制文件
	err := WriteProtobufToBinaryFile(laptop, binaryFile)
	require.NoError(t, err)

	// 从二进制文件读取到 pb.Laptop,从而对比两个laptop是否相同
	laptop2 := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop, laptop2))

	// 写入 JSON 文件
	err = WriteProtobufToJSONFile(laptop, jsonFile)
	require.NoError(t, err)

	// 对比两个文件大小
	binaryInfo, err := os.Stat(binaryFile)
	require.NoError(t, err)

	jsonInfo, err := os.Stat(jsonFile)
	require.NoError(t, err)

	require.Less(t, binaryInfo.Size(), jsonInfo.Size())

	fmt.Printf("binary file size: %d, json file size: %d\n", binaryInfo.Size(), jsonInfo.Size())
}
