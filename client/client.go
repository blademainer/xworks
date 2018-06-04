package main

import (
	"net"
	logger "github.com/sirupsen/logrus"
	"github.com/blademainer/xworks/network"
	pb "github.com/golang/protobuf/proto"
	"github.com/blademainer/xworks/proto"
	"fmt"
	"time"
	"github.com/blademainer/xworks/common"
	"bufio"
	"os"
	"github.com/golang/protobuf/ptypes/any"
)

const (
	ENV_SERVER_ADDR     = "SERVER_ADDR"
	DEFAULT_SERVER_ADDR = "127.0.0.1:1717"
)

func main() {
	common.SetLogLevel()

	addr, b := os.LookupEnv(ENV_SERVER_ADDR)
	if !b {
		addr = DEFAULT_SERVER_ADDR
	}
	if conn, e := net.Dial("tcp", addr); e == nil {
		done := make(chan bool)
		var timeOut time.Time
		conn.SetReadDeadline(timeOut)
		reader := bufio.NewReader(conn)
		go func(exit chan bool) {
			for i := 0; ; i++ {
				if bytes, e := network.ReadBytes(reader, conn); e == nil {
					fmt.Printf("Receive: %v size: %d \n", bytes, len(bytes))
					request := &proto.Request{}
					if e = pb.Unmarshal(bytes, request); e == nil {
						fmt.Printf("Unmarshal: %d bytes to Request: %v \n", len(bytes), request)
					} else {
						logger.Errorf("Failed to unmarshal: %v error: %v", request, e)
					}
				} else if closedErr, closed := e.(network.ConnectionClosedError); closed {
					logger.Errorf(closedErr.Error())
					exit <- true
					break
				}
			}
			exit <- true
		}(done)
		i := 0
		for ; ; i++ {
			body := &any.Any{Value: []byte(fmt.Sprintf("%s%d", "Hello!server!", i))}
			request := &proto.Request{
				Name: fmt.Sprintf("%s%d", "Hello!world!", i),
				Body: body,
			}
			if bytes, e := pb.Marshal(request); e == nil {
				//bytes = append(bytes, '\n')
				conn.Write(bytes)
				logger.Debugf("Write bytes: %v length: %d", bytes, len(bytes))
			} else {
				logger.Errorf("Failed to marshal: %v error: %v", request, e)
			}
		}
	} else {
		logger.Errorf("Connection server error! msg: %s", e.Error())
	}
}
