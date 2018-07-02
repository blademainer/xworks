package network

import (
	"bufio"
	"github.com/blademainer/xworks/logger"
	"net"
	"sync"
	"time"
)

type (
	endpoint struct {
		connection net.Conn
		writerCh   chan []byte
		readCh     chan []byte
		closed     bool
		closeOnce  sync.Once
	}
)

func InitEndpoint(conn net.Conn) (e *endpoint) {
	conn.SetReadDeadline(time.Time{})
	e = &endpoint{}
	e.connection = conn
	e.writerCh = make(chan []byte)
	e.readCh = make(chan []byte)
	e.initWriterWorker()
	e.initReaderWorker()
	return e
}

func (endpoint *endpoint) ReadChannel() (readCh chan []byte) {
	return endpoint.readCh
}

func (endpoint *endpoint) Write(data []byte) {
	if endpoint.closed {
		return
	}
	endpoint.writerCh <- data
}

func (endpoint *endpoint) Close() {
	endpoint.closeOnce.Do(func() {
		logger.Log.Warnf("Closing endpoint: %v", endpoint)
		endpoint.closed = true
		close(endpoint.readCh)
		close(endpoint.writerCh)
		if e := endpoint.connection.Close(); e != nil {
			logger.Log.Errorf("Error to close! error: %s", e.Error())
		}
		logger.Log.Warnf("Closed endpoint: %v", endpoint)

	})
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
		}
	}
	logger.Log.Warnf("Closed write worker for endpoint: %v", endpoint)

}

func (endpoint *endpoint) initReaderWorker() {
	logger.Log.Infof("Starting read worker for endpoint: %v", endpoint)
	go endpoint.processReadData()
	logger.Log.Infof("End read worker for endpoint: %v", endpoint)
}

func (endpoint *endpoint) processReadData() {
	reader := bufio.NewReader(endpoint.connection)
	for bytes, e := ReadBytes(reader, endpoint.connection); e == nil; bytes, e = ReadBytes(reader, endpoint.connection) {
		//fmt.Println(e.Error())
		logger.Log.Infof("Read: %v", bytes)
		if endpoint.closed {
			logger.Log.Warnf("Closed endpoint! %v", endpoint.connection)
			return
		}
		endpoint.readCh <- bytes
	}
	logger.Log.Warnf("Closing endpoint! %v", endpoint)
	endpoint.Close()
}
