package network

import (
	"bufio"
	"github.com/blademainer/xworks/pkg/logger"
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
		closeCh    chan bool
		closeOnce  sync.Once
		wg         sync.WaitGroup
	}
)

func InitEndpoint(conn net.Conn) (e *endpoint) {
	conn.SetReadDeadline(time.Time{})
	e = &endpoint{}
	e.connection = conn
	e.writerCh = make(chan []byte, 1024)
	e.readCh = make(chan []byte, 1024)
	e.closeCh = make(chan bool, 1)
	e.initWriterWorker()
	e.initReaderWorker()
	return e
}

func (endpoint *endpoint) ReadChannel() (readCh chan []byte) {
	return endpoint.readCh
}

func (endpoint *endpoint) Write(data []byte) error {
	if endpoint.closed {
		return &ConnectionClosedError{Message: "Connection closed!"}
	}
	endpoint.writerCh <- data
	return nil
}

func (endpoint *endpoint) Close() {
	endpoint.closeOnce.Do(func() {
		logger.Log.Warnf("Closing endpoint: %v", endpoint)
		// waiting for read channel and write channel done.
		endpoint.wg.Add(2)
		endpoint.closed = true
		endpoint.closeCh <- true
		endpoint.wg.Wait()
		close(endpoint.readCh)
		close(endpoint.writerCh)
		close(endpoint.closeCh)
		if e := endpoint.connection.Close(); e != nil {
			logger.Log.Errorf("Error to close! error: %s", e.Error())
		}
		logger.Log.Warnf("Closed endpoint: %v", endpoint)

	})
}

func (endpoint *endpoint) initWriterWorker() {
	go endpoint.processWriteData()
}

func (endpoint *endpoint) processWriteData() {
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
			break
		}
	}
	endpoint.wg.Done()
	logger.Log.Warnf("Closed write worker for endpoint: %v", endpoint)

}

func (endpoint *endpoint) initReaderWorker() {
	logger.Log.Infof("Starting read worker for endpoint: %v", endpoint)
	go endpoint.processReadData()
	logger.Log.Infof("End read worker for endpoint: %v", endpoint)
}

func (endpoint *endpoint) processReadData() {
	reader := bufio.NewReader(endpoint.connection)
	for !endpoint.closed {
		select {
		case <-endpoint.closeCh:
			break
		default:
			if bytes, e := ReadBytes(reader, endpoint.connection); e != nil {
				logger.Log.Warnf("Closing endpoint! %v cause error: %v", endpoint.connection, e)
				go func() {
					endpoint.Close()
				}()
				break
			} else {
				logger.Log.Debugf("Read data: %v from endpoint! %v", bytes, endpoint.connection)
				endpoint.readCh <- bytes
			}
		}

	}
	logger.Log.Warnf("Closed reader for endpoint! %v", endpoint.connection)
	endpoint.wg.Done()
}
