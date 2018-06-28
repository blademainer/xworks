package network

import (
	"fmt"
	"github.com/blademainer/xworks/logger"
	"net"
)

type (
	endpoint struct {
		connection net.Conn
		writerCh   chan []byte
		readCh     chan []byte
		closeCh    chan bool
		closed     bool
	}
)

func InitEndpoint(conn net.Conn) (e *endpoint) {
	e = &endpoint{}
	e.connection = conn
	e.writerCh = make(chan []byte)
	e.closeCh = make(chan bool)
	e.initWriterWorker()
	return e
}

func (endpoint *endpoint) Write(data []byte) {
	if endpoint.closed {
		return
	}
	endpoint.writerCh <- data
}

func (endpoint *endpoint) Close() {
	endpoint.closed = true
	endpoint.closeCh <- true
	if e := endpoint.connection.Close(); e != nil {
		fmt.Println("Error to close! error: ", e.Error())
	}
}

func (endpoint *endpoint) initWriterWorker() {
	go endpoint.processData()
}

func (endpoint *endpoint) processData() {
	for !endpoint.closed {
		select {
		case data := <-endpoint.writerCh:
			if n, err := endpoint.connection.Write(data); err != nil {
				logger.Log.Debugf("Failed to write data error: %v data: %v ", data, err)
			} else {
				if logger.Log.IsDebugEnabled() {
					logger.Log.Debugf("Success write data, size: %d", n)
				}
			}
		case <-endpoint.closeCh:
			return
		}
	}

}

func (endpoint *endpoint) initReaderWorker() {

}
