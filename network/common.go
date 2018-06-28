package network

import (
	"bufio"
	"errors"
	"fmt"
	logger "github.com/sirupsen/logrus"
	"io"
	"net"
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
//	buf := make([]byte, 1024)
//	//conn.SetReadDeadline(time.Now())
//	if length, err := conn.Read(buf); err != nil {
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
//		logger.Debugf("Read bytes: %d \n", length)
//		return buf[:length], nil
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

//func ReadBytes(conn net.Conn) ([]byte, error) {
//	reader := bufio.NewReader(conn)
//	if line, _, err := reader.ReadLine(); err != nil {
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
//
//	} else {
//		//var zero time.Time
//		//conn.SetReadDeadline(zero)
//		logger.Debugf("Get byte: %v", string(line))
//		return line, nil
//	}
//}

func ReadBytes(reader *bufio.Reader, conn net.Conn) ([]byte, error) {
	if line, _, err := reader.ReadLine(); err != nil {
		if err == io.EOF {
			//logger.Debugf("%s detected closed connection. error: %s type: %T", conn, err.Error(), err)
			err = &ConnectionClosedError{err.Error(), conn, err}
			logger.Warnf("Closed conn: %v, err: %s", conn, err.Error())
			//if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			//}
			conn.Close()
			conn = nil
		}
		return nil, err

	} else if len(line) <= 0 {
		e := errors.New("No data!")
		logger.Warnf(e.Error())
		return nil, e
	} else {
		//var zero time.Time
		//conn.SetReadDeadline(zero)
		//logger.Debugf("Get byte: %v", string(line))
		//bytes := make([]byte, len(line)+1)
		//copy(bytes, []byte{10})
		//copy(bytes[1:], line)
		bytes := Insert(line, []byte{'\n'}, 0)
		//line = append(line, '\n')
		return bytes, nil
	}
}

func Insert(slice, insertion []byte, index int) []byte {
	result := make([]byte, len(slice)+len(insertion))
	at := copy(result, slice[:index])
	at += copy(result[at:], insertion)
	copy(result[at:], slice[index:])
	return result
}

//func Insert(slice, insertion []interface{}, index int) []interface{} {
//	result := make([]interface{}, len(slice)+len(insertion))
//	at := copy(result, slice[:index])
//	at += copy(result[at:], insertion)
//	copy(result[at:], slice[index:])
//	return result
//}
