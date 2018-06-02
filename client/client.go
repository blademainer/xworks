package main

import (
	"net"
	logger "github.com/sirupsen/logrus"
	"github.com/blademainer/xworks/network"
	"fmt"
)

func main() {
	if conn, e := net.Dial("tcp", "127.0.0.1:1717"); e == nil {
		done := make(chan bool)

		go func(exit chan bool) {
			for i := 0; i < 5; i++ {
				if bytes, e := network.ReadBytes(conn); e == nil {
					fmt.Println("Receive: ", string(bytes))
				} else {
					exit <- true
					break
				}
			}
			exit <- true
		}(done)

		<-done
	} else {
		logger.Errorf("Connection server error! msg: %s", e.Error())
	}
}
