package main

import (
	"net"
	"os"
	log "github.com/sirupsen/logrus"
	"strconv"
	"fmt"
)

const (
	ENV_SERVER_PORT = "SERVER_PORT"
	ENV_LOG_LEVEL   = "LOG_LEVEL"

	DEFAULT_SERVER_PORT    = "10060"
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
	log.Infof("Starting server... %v: %s", network, address)
	listener, e := net.Listen(network, address)
	if e != nil {
		log.Errorf("Failed to start server: %v", e.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		if err == nil {
			go processConn(conn)
		} else {
			log.Errorf("Failed to accept! error: %s", err.Error())
		}
	}
}

func processConn(conn net.Conn) {
	log.Debugf("Accepted connection: %v", conn)
	conn.Write([]byte("Hello world!"))
}

func setLogLevel() {
	logLevel, _ := os.LookupEnv(ENV_LOG_LEVEL)
	switch logLevel {
	case LOG_LEVEL_DEBUG:
		log.SetLevel(log.DebugLevel)
	case LOG_LEVEL_INFO:
		log.SetLevel(log.InfoLevel)
	case LOG_LEVEL_WARN:
		log.SetLevel(log.WarnLevel)
	case LOG_LEVEL_ERROR:
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}

func start() {
	port, b := os.LookupEnv(ENV_SERVER_PORT)
	if !b {
		port = DEFAULT_SERVER_PORT
		log.Warnf("Not found env %port so sets to default value: %v", ENV_SERVER_PORT, DEFAULT_SERVER_PORT)
	}
	server := &Server{}
	server.Network = DEFAULT_SERVER_NETWORK
	if p, err := strconv.ParseInt(port, 10, 32); err != nil {
		log.Errorf("Parse port error! error: %s", err.Error())
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
