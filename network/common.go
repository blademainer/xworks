package network

import (
	"net"
	logger "github.com/sirupsen/logrus"
	"fmt"
	"bufio"
	"io"
)

type ConnectionClosedError struct {
	Message    string
	Connection net.Conn
	OriginErr  error
}

func (e ConnectionClosedError) Error() string {
	return fmt.Sprintf("Detected closed connection: %v, message: %s, error type: %T", e.Connection, e.Message, e.OriginErr)
}

//func ReadBytes(conn net.Conn) ([]byte, error) {
//	one := make([]byte, 32)
//	//conn.SetReadDeadline(time.Now())
//	if _, err := conn.Read(one); err != nil {
//		if err == io.EOF {
//			//logger.Debugf("%s detected closed connection. error: %s type: %T", conn, err.Error(), err)
//			err = &ConnectionClosedError{err.Error(), conn, err}
//			logger.Errorf("Closed conn: %v, err: %s", conn, err.Error())
//			//if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
//			//}
//			conn.Close()
//			conn = nil
//		}
//		return nil, err
//	} else {
//		//var zero time.Time
//		//conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
//		return one, nil
//	}
//}

//func ReadBytes(conn net.Conn) ([]byte, error) {
//	bytes := []byte{}
//	conn.SetReadDeadline(time.Now())
//	if _, err := conn.Read(bytes); err != nil {
//		logger.Debugf("%s detected closed connection. error: %s type: %T", conn, err.Error(), err)
//		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
//			logger.Errorf("Closed conn: %v", conn)
//		}
//		conn.Close()
//		conn = nil
//		return nil, err
//	} else {
//		var zero time.Time
//		conn.SetReadDeadline(zero)
//		//logger.Debugf("Get byte: %v", string(line))
//		return bytes, nil
//	}
//}

func ReadBytes(conn net.Conn) ([]byte, error) {

	if line, _, err := bufio.NewReader(conn).ReadLine(); err != nil {
		if err == io.EOF {
			//logger.Debugf("%s detected closed connection. error: %s type: %T", conn, err.Error(), err)
			err = &ConnectionClosedError{err.Error(), conn, err}
			logger.Errorf("Closed conn: %v, err: %s", conn, err.Error())
			//if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			//}
			conn.Close()
			conn = nil
		}
		return nil, err

	} else {
		//var zero time.Time
		//conn.SetReadDeadline(zero)
		logger.Debugf("Get byte: %v", string(line))
		return line, nil
	}
}
