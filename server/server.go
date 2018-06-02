package main

import (
	"net"
	"os"
	logger "github.com/sirupsen/logrus"
	"strconv"
	"fmt"
	"github.com/blademainer/xworks/network"
	"time"
)

const (
	ENV_SERVER_PORT = "SERVER_PORT"
	ENV_LOG_LEVEL   = "LOG_LEVEL"

	DEFAULT_SERVER_PORT    = "1717"
	DEFAULT_SERVER_NETWORK = "tcp"
	DEFAULT_LOG_LEVEL      = "debug"

	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
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
	go func() {
		for _, e := network.ReadBytes(conn); e == nil; _, e = network.ReadBytes(conn) {
			//fmt.Println(e.Error())
			//fmt.Println("Read: ", string(bytes))
		}
	}()
	for i := 0; i < 100; i++ {
		conn.Write([]byte("Hello world!\n"))
	}
}

func setLogLevel() {
	logLevel, _ := os.LookupEnv(ENV_LOG_LEVEL)
	switch logLevel {
	case LOG_LEVEL_DEBUG:
		logger.SetLevel(logger.DebugLevel)
	case LOG_LEVEL_INFO:
		logger.SetLevel(logger.InfoLevel)
	case LOG_LEVEL_WARN:
		logger.SetLevel(logger.WarnLevel)
	case LOG_LEVEL_ERROR:
		logger.SetLevel(logger.ErrorLevel)
	default:
		logger.SetLevel(logger.DebugLevel)
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
	setLogLevel()
	start()

}
