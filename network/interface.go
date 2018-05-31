package network

import (
	"net"
	"github.com/blademainer/xworks/proto"
)

type (
	Context struct {
		Connection *net.Conn

	}

	Service struct {
		ServiceName string
	}

	IServer interface {
		Start(network, address string)
	}

	IService interface {
		Action(request proto.Request) proto.Response
	}

	IClient interface {
		Start(network, address string)
	}
)
