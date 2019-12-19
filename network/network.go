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
	T             time.Duration
)

// initialisation de la partie reseau
func InitNetwork(nb_site int, site_add []string, apt_site []int, id int) {
	nbre_site = nb_site
	all_add = site_add
	all_apt = apt_site
	myId = id
	next = (myId + 1) % nbre_site
	ack = make(chan bool, 1)
	T = 1
}

func MsgTo(msg string) {

	conn, err := net.Dial("udp", all_add[next])
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// simulation d'un temps long de transmission
	time.Sleep(T * time.Second)
	_, _ = fmt.Fprintf(conn, msg)

	// se mets en attente du ack après l'envoi
	// de chaque message
	buf := make([]byte, 1)

	// On n'écoute que pendant la période du timeout
	_ = conn.SetDeadline(time.Now().Add(2 * T * time.Second))

	conn.Read(buf)
	if buf[0] == 'O' {
		ack <- true
	}

	ConnectionHandle(conn, buf, msg)
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

			// répond par un ack à la réception d'un message
			// (ack exclu car sinon boucle infinie de ack)
			if s.Text() != "O" {
				_, err = conn.WriteTo([]byte("O"), previousSiteAddr)
			}

			return msg
		}
	}
}

func ConnectionHandle(conn net.Conn, buf []byte, msg string) {

	select {
	case <-ack:

		// réinitialise le next
		next = (myId + 1) % nbre_site

	case <-time.After(2 * T * time.Second):

		// passe au suivant
		next = (next + 1) % nbre_site

		// renvoie le message au suivant
		fmt.Println("Pas recu de ACK après 2 sec, passe au suivant")
		MsgTo(msg)
	}
}
