package processor

import (
	"bufio"
	"fmt"
	"github.com/blademainer/xworks/logger"
	"github.com/blademainer/xworks/network"
	"github.com/blademainer/xworks/proto"
	pb "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"net"
	"os"
	"strconv"
	"time"
)

var Logger = logger.Log

const (
	ENV_SERVER_PORT        = "SERVER_PORT"
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
	Logger.Infof("Starting server... %v: %s", network, address)
	listener, e := net.Listen(network, address)
	if e != nil {
		Logger.Errorf("Failed to start server: %v", e.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		if err == nil {
			go processConn(conn)
		} else {
			Logger.Errorf("Failed to accept! error: %s", err.Error())
		}
	}
}

func processConn(conn net.Conn) {
	conn.SetReadDeadline(time.Time{})
	Logger.Debugf("Accepted connection: %v", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	go func() {
		for bytes, e := network.ReadBytes(reader, conn); e == nil; bytes, e = network.ReadBytes(reader, conn) {
			//fmt.Println(e.Error())
			Logger.Debugf("Read: %v", bytes)
		}
	}()
	i := 0
	for ; ; i++ {
		body := &any.Any{Value: []byte(fmt.Sprintf("%s%d", "Hello!client!", i))}
		request := &proto.Request{
			Name: fmt.Sprintf("%s%d", "Hello!world!", i),
			Body: body,
		}
		if bytes, e := pb.Marshal(request); e == nil {
			//bytes = append(bytes, '\n')
			conn.Write(bytes)
			Logger.Debugf("Write bytes: %v length: %d", bytes, len(bytes))
		} else {
			Logger.Errorf("Failed to marshal: %v error: %v", request, e)
		}
	}
}

func Start() {
	port, b := os.LookupEnv(ENV_SERVER_PORT)
	if !b {
		port = DEFAULT_SERVER_PORT
		Logger.Warnf("Not found env %s so sets to default value: %v", ENV_SERVER_PORT, DEFAULT_SERVER_PORT)
	}
	server := &Server{}
	server.Network = DEFAULT_SERVER_NETWORK
	if p, err := strconv.ParseInt(port, 10, 32); err != nil {
		Logger.Errorf("Parse port error! error: %s", err.Error())
		return
	} else {
		pp := uint32(p)
		fmt.Println(pp)
		server.Port = pp
		server.Start(server.Network, fmt.Sprintf(":%d", server.Port))
	}
}
