package network

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
)

func ConnTo(site int, Toaddress string, Fromaddress string){
	conn, err := net.Dial("udp", Toaddress)

	if err == nil {
		MsgTo("Simon", conn)
		MsgFrom("udp",Fromaddress)
		log.Println("i am connected to site " + strconv.Itoa(site))
	} else{
		log.Fatal(err)
	}
}

func MsgTo(msg string, conn net.Conn){
	fmt.Fprintln(conn, msg)
}


func MsgFrom(network string, address string) string {
	buf := make([]byte, 256)
	conn,_ := net.ListenPacket(network,address)

	fmt.Println("t es la")

	n, _, _ := conn.ReadFrom(buf)
	s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
	fmt.Println("t es la2")


	fmt.Println(s.Text())

	return string(buf)
}