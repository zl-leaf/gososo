package socket
import(
	"net"
	"bytes"
	"encoding/binary"
)

func Write(conn net.Conn, data []byte) (int, error){
	dataCount := len(data)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(dataCount))

	builder := bytes.Buffer{}
	builder.Write(buf)
	builder.Write(data)

	count, err := conn.Write(builder.Bytes())
	return count, err
}

func Read(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 8)
	_, err := conn.Read(buf)
	if err != nil {
		return nil,err
	}
	count := int64(binary.BigEndian.Uint64(buf))
	builder := bytes.Buffer{}
	for {
		if count <= 0 {
			break
		}
		data := make([]byte, 128)
		dataCount, err := conn.Read(data)
		if err != nil {
			return nil,err
		}
		builder.Write(data[:dataCount])
		count -= int64(dataCount)
	}

	return builder.Bytes(), nil
}