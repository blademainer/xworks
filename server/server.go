package main

import (
	"net"
	"os"
	"strconv"
	"fmt"
	"github.com/blademainer/xworks/network"
	"time"
	"github.com/blademainer/xworks/proto"
	pb "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/blademainer/xworks/common"
	logger "github.com/sirupsen/logrus"
	"bufio"
)

const (
	ENV_SERVER_PORT = "SERVER_PORT"
	DEFAULT_SERVER_PORT    = "1717"
	DEFAULT_SERVER_NETWORK = "tcp"
)

type (
	Server struct {
		Network string
		Port    uint32
	}
)

func (server *Server) Start(network, address string) {
	logger.Infof("Starting server... %v: %s", network, address)
	listener, e := net.Listen(network, address)
	if e != nil {
		logger.Errorf("Failed to start server: %v", e.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		if err == nil {
			go processConn(conn)
		} else {
			logger.Errorf("Failed to accept! error: %s", err.Error())
		}
	}
}

func processConn(conn net.Conn) {
	conn.SetReadDeadline(time.Time{})
	logger.Debugf("Accepted connection: %v", conn)
	reader := bufio.NewReader(conn)
	go func() {
		for _, e := network.ReadBytes(reader, conn); e == nil; _, e = network.ReadBytes(reader, conn) {
			//fmt.Println(e.Error())
			//fmt.Println("Read: ", string(bytes))
		}
	}()
	for i := 0; i < 100; i++ {
		body := &any.Any{Value: []byte(fmt.Sprintf("%s%d", "Hello!world!", i))}
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
}



func start() {
	port, b := os.LookupEnv(ENV_SERVER_PORT)
	if !b {
		port = DEFAULT_SERVER_PORT
		logger.Warnf("Not found env %port so sets to default value: %v", ENV_SERVER_PORT, DEFAULT_SERVER_PORT)
	}
	server := &Server{}
	server.Network = DEFAULT_SERVER_NETWORK
	if p, err := strconv.ParseInt(port, 10, 32); err != nil {
		logger.Errorf("Parse port error! error: %s", err.Error())
		return
	} else {
		pp := uint32(p)
		fmt.Println(pp)
		server.Port = pp
		server.Start(server.Network, fmt.Sprintf(":%d", server.Port))
	}
}

func main() {
	common.SetLogLevel()
	start()

}
