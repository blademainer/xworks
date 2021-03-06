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

	IService interface {
		Action(request proto.Request) proto.Response
	}

	Client interface {
		SendMessage(request proto.Request) (proto.Response, error)
	}

	Listener interface {

	}

	Server interface {
		RegisterClient(client *Client)
		RegisterListener()
	}


)
