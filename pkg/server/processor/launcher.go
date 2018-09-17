package processor

import (
	"fmt"
	"github.com/blademainer/xworks/logger"
	"github.com/blademainer/xworks/network"
	"github.com/blademainer/xworks/proto"
	pb "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"net"
	"os"
	"strconv"
)

const (
	ENV_SERVER_PORT        = "SERVER_PORT"
	DEFAULT_SERVER_PORT    = "1717"
	DEFAULT_SERVER_NETWORK = "tcp"
)

type (
	Server struct {
		Network string
		Port    uint32
		*network.EndpointManager
	}
)

var Log = logger.NewLogger()

func (server *Server) Start(listen, address string) {
	Log.Infof("Starting server... %v: %s", listen, address)
	listener, e := net.Listen(listen, address)
	if e != nil {
		Log.Errorf("Failed to start server: %v", e.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		if err == nil {
			go server.processConn(conn)
		} else {
			Log.Errorf("Failed to accept! error: %s", err.Error())
		}
	}
}

func (server *Server) processConn(conn net.Conn) {
	Log.Debugf("Accepted connection: %v", conn.RemoteAddr())
	//reader := bufio.NewReader(conn)
	//go func() {
	//	for bytes, e := network.ReadBytes(reader, conn); e == nil; bytes, e = network.ReadBytes(reader, conn) {
	//		//fmt.Println(e.Error())
	//		Log.Debugf("Read: %v", bytes)
	//	}
	//}()
	e := server.InitEndpoint(conn)
	readCh := e.ReadChannel()
	for data := range readCh {
		request := &proto.Request{}
		if err := pb.Unmarshal(data, request); err != nil {
			Log.Errorf("Error when unmarshal data: %v error: %s", data, err.Error())
		} else {
			Log.Infof("Unmarshal data: %v to entity: %v", data, request)
		}
	}

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
			Log.Debugf("Write bytes: %v length: %d", bytes, len(bytes))
		} else {
			Log.Errorf("Failed to marshal: %v error: %v", request, e)
		}
	}
}

func Start() {
	config := logger.LoggerConfig{
		Level: logger.LOG_LEVEL_DEBUG,
	}
	if e := Log.Init(config); e != nil {
		fmt.Printf("Error to init logger! error: %s", e.Error())
	}

	port, b := os.LookupEnv(ENV_SERVER_PORT)
	if !b {
		port = DEFAULT_SERVER_PORT
		Log.Warnf("Not found env %s so sets to default value: %v", ENV_SERVER_PORT, DEFAULT_SERVER_PORT)
	}
	manager := network.InitEndpointManager()
	server := &Server{EndpointManager: manager}
	server.Network = DEFAULT_SERVER_NETWORK
	if p, err := strconv.ParseInt(port, 10, 32); err != nil {
		Log.Errorf("Parse port error! error: %s", err.Error())
		return
	} else {
		pp := uint32(p)
		fmt.Println(pp)
		server.Port = pp
		server.Start(server.Network, fmt.Sprintf(":%d", server.Port))
	}
}
