package Tool

//解决粘包问题
import (
	"encoding/binary"
	"io"
	"net"
)

// 不写成结构体,要不然 客户端和服务端不一致

// Send 先发 4 字节长度（大端），再发数据
func Send(conn net.Conn, data string) error {
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(data)))
	if _, err := conn.Write(header); err != nil {
		return err
	}
	_, err := conn.Write([]byte(data))
	return err
}

// Recv 先读 4 字节长度，再读完整数据
func Recv(conn net.Conn) (string, error) {
	hdr := make([]byte, 4)
	if _, err := io.ReadFull(conn, hdr); err != nil {
		return "", err
	}
	sz := binary.BigEndian.Uint32(hdr)
	body := make([]byte, sz)
	if _, err := io.ReadFull(conn, body); err != nil {
		return "", err
	}
	return string(body), nil
}
