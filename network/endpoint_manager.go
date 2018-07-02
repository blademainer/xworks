package network

import (
	"net"
	"sync"
)

type EndpointManager struct {
	endpoints map[net.Conn]*endpoint
	sync.Mutex
}

func InitEndpointManager() *EndpointManager {
	endpointManager := &EndpointManager{endpoints: make(map[net.Conn]*endpoint)}
	return endpointManager
}

func (e *EndpointManager) InitEndpoint(connection net.Conn) *endpoint {
	e.Lock()
	defer e.Unlock()
	result := InitEndpoint(connection)
	e.endpoints[connection] = result
	return result
}

func (e *EndpointManager) CloseEndpoint(ep endpoint) {
	e.Lock()
	defer e.Unlock()
	ep.Close()
	delete(e.endpoints, ep.connection)
}
