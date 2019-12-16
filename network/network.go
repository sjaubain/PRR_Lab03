package network

import (
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

