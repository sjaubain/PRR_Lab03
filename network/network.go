package network

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func ConnTo(site int, Toaddress string, Fromaddress string){
	conn, err := net.Dial("udp", Toaddress)

	/*if err == nil {
		MsgTo("Simon", conn)
		log.Println("i am connected to site " + strconv.Itoa(site))
	} else{
		log.Fatal(err)
	}*/

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		fmt.Println(("test 1"))
		mustCopy(os.Stdout, conn)
		fmt.Println(("test 2"))

	}()
	mustCopy(conn, os.Stdin)
}
func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func MsgTo(msg string, conn net.Conn){
	//p :=  make([]byte, 2048)
	//fmt.Fprintf(conn, msg)
	//_, _ = bufio.NewReader(conn).Read(p)
	_, _ = conn.Write([]byte(msg))
}


func MsgFrom(network string, address string) {
	/*buf := make([]byte, 256)

	for{
		conn,_ := net.ListenPacket(network,address)

		fmt.Println("t es la")

		n, addr, _ := conn.ReadFrom(buf)
		fmt.Println("T a quelque chose la?")
		_ = bufio.NewScanner(bytes.NewReader(buf[0:n]))
		fmt.Println("t es la2")


		fmt.Println(addr)
	}*/

	conn, err := net.ListenPacket(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, cliAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		s := bufio.NewScanner(bytes.NewReader(buf[0:n]))
		for s.Scan() {
			s := s.Text() + " from " + cliAddr.String() + "\n"
			if _, err := conn.WriteTo([]byte(s), cliAddr); err != nil {
				log.Fatal(err)
			}
		}
	}
}