package network

import (
	"net"
	//	"reflect"
	"testing"
)

const (
	address = "127.0.0.1:3006"
)


func TestListenRun(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.ParseIP(address),
	})
	if err != nil {
		t.Error("Error ", err)
	}
	defer conn.Close()
}
