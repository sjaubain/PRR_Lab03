package network

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	nbre_site     int
	all_add       []string
	all_apt       []int
	myId          int
	next          int
	ack           chan bool
	notConnection chan bool
)

func InitNetwork(nb_site int, site_add []string, apt_site []int, id int) {
	nbre_site = nb_site
	all_add = site_add
	all_apt = apt_site
	myId = id
	next = (myId + 1) % nbre_site
	ack = make(chan bool)
	notConnection = make(chan bool)

}

func MsgTo(msg string) {

	for i := true; i; i = <-notConnection {
		conn, err := net.Dial("udp", all_add[next])
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		//fmt.Println("\n me : " + strconv.Itoa(myId) + " next : " + strconv.Itoa(next))

		time.Sleep(time.Second) // msg transmis chaque seconde
		_, _ = fmt.Fprintf(conn, msg)

		// se mets en attente du ack après l'envoi
		// de chaque message
		buf := make([]byte, 1)
		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))

		conn.Read(buf)
		if buf[0] == 'O' {
			//fmt.Println("Recu un ACK pour le message [" + msg + "]")
			ack <- true
		}

		go ConnectionHandle(conn, buf, msg)
	}

}

func MsgFrom(network string, address string) string {
	conn, err := net.ListenPacket(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, previousSiteAddr, err := conn.ReadFrom(buf)

		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {

			msg := s.Text()

			fmt.Println(s.Text())
			// répond par un ack à la réception d'un message
			// (ack exclu car sinon boucle infinie de ack)
			if s.Text() != "O" {
				_, _ = conn.WriteTo([]byte("O"), previousSiteAddr)
			}

			return msg
		}
	}
}

func ConnectionHandle(conn net.Conn, buf []byte, msg string) {

	select {
	case <-ack:
		next = (myId + 1) % nbre_site
		notConnection <- false

	case <-time.After(2 * time.Second):
		next = (next + 1) % nbre_site

		// renvoie le message au suivant
		fmt.Println("Pas recu de ACK après 2 sec, passe au suivant")
		MsgTo(msg)
		notConnection <- true
	}
}
