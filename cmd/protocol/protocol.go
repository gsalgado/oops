package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/vds/oops"
)

const (
	lengthSize = 8
	ACK        = "ack"
	sizeOfACK  = len(ACK)
)

var (
	bACK = []byte(ACK)
)

// Send the oops through the connection.
func SendOops(conn *net.TCPConn, encodedOops []byte) (err error) {
	oopsLength := int64(len(encodedOops))
	buf := make([]byte, lengthSize)
	binary.PutVarint(buf, oopsLength)
	bytesSent, err := conn.Write(append(buf, encodedOops...))
	if err != nil {
		return err
	}
	if bytesSent != int(oopsLength) {
		fmt.Errorf("sending oops: sent %v bytes of %v", bytesSent, oopsLength)
	}
	return
}

//  Send the oops through the connection.
func ReceiveOops(conn net.Conn) (err error) {

	buf := make([]byte, lengthSize)
	bytesReceived, err := conn.Read(buf)
	if err != nil {
		return
	}
	if bytesReceived != lengthSize {
		return fmt.Errorf("oops length: read %v byte out of %v", bytesReceived, lengthSize)
	}
	r := bytes.NewReader(buf)
	l, err := binary.ReadVarint(r)

	buf = make([]byte, l)
	bytesReceived, err = conn.Read(buf)
	if err != nil {
		return fmt.Errorf("receiving oops: %v", err)
	}
	if int64(bytesReceived) != l {
		return fmt.Errorf("receiving oops: %v byte out of %v", bytesReceived, l)
	}

	// FIXME store oops somewhere
	var oops oops.Oops
	err = oops.Unmarshal(buf)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", oops.Error)
	return
}

// Send acknowlegment of received information.
func SendAck(conn net.Conn) (err error) {
	bytesSent, err := conn.Write(bACK)
	if err != nil {
		return err
	}
	if bytesSent != sizeOfACK {
		return fmt.Errorf("sending ack: sent %v bytes of %v", bytesSent, sizeOfACK)
	}
	return
}

// Receive acknowlegment of received information.
func ReceiveAck(conn *net.TCPConn) (err error) {
	buf := make([]byte, 3)
	bytesReceived, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if bytesReceived != 3 {
		return fmt.Errorf("ack: read %v byte out of 3", bytesReceived)
	}
	if string(buf) != ACK {
		return fmt.Errorf("bad ack response: %v")
	}
	return
}