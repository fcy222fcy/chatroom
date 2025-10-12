package Tool

//解决粘包问题
import (
	"encoding/binary"
	"io"
	"net"
)

// Send 先发 4 字节长度（大端），再发数据
func Send(conn net.Conn, data []byte) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))
	if _, err := conn.Write(header); err != nil {
		return err
	}
	_, err := conn.Write(data)
	return err
}

// Recv 先读 4 字节长度，再读完整数据
func Recv(conn net.Conn) ([]byte, error) {
	hdr := make([]byte, 4)
	if _, err := io.ReadFull(conn, hdr); err != nil {
		return nil, err
	}
	sz := binary.BigEndian.Uint32(hdr)
	body := make([]byte, sz)
	if _, err := io.ReadFull(conn, body); err != nil {
		return nil, err
	}
	return body, nil
}
