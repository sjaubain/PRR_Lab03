package network

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
)

var(
	nbre_site int
	all_add []string
	all_apt []int
	myId int
	next int
)

func InitNetwork(nb_site int, site_add []string, apt_site []int, id int){
	nbre_site = nb_site
	all_add = site_add
	all_apt = apt_site
	myId = id
	next = myId + 1

}

func MsgTo(msg string){
	conn, err := net.Dial("udp", all_add[next%nbre_site])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn,msg)
}

func MsgFrom(network string, address string) string{
	conn, err := net.ListenPacket(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n,_, err := conn.ReadFrom(buf)

		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			msg := s.Text()
			return msg
		}
	}
}