package algoCR

import (
	"PRR_Lab03/network"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var(
	site_id int
	apt int
	idApt string
	etat string
	elu int
)

func InitAlgo(identifiant int, apti int){
	site_id = identifiant
	apt = apti
	idApt = strconv.Itoa(site_id) + "-" + strconv.Itoa(apt)
}


func Election(){
	network.MsgTo("A" + idApt)
	etat = "A"
}

func RcptAnnonce(list string){

	/// CONTAIN NE MARCHE PAS !!!!!!!!!!!!!
	if strings.Contains(idApt,list){
		tabList := strings.Split(list,";")
		aptMax := 0
		elu = 0
		for i := range tabList{
			id,apt := getApt(tabList[i])
			if apt > aptMax{
				aptMax = apt
				elu = id
			}
		}
		network.MsgTo("R" + strconv.Itoa(elu) + "," + strconv.Itoa(site_id))
		etat = "R"
	}else{
		list += ";" + idApt
		network.MsgTo("A"+list)
		etat = "A"
	}
}



func RcptResultat(i string, list string){
	tabProc := strings.Split(list,";")
	for j:= range tabProc{
		if strconv.Itoa(site_id) == tabProc[j]{
			etat = "N"
		}else if etat == "R" && strconv.Itoa(elu) != i{
			Election()
		}else if etat == "A"{
			elu , _ = strconv.Atoi(i)
			list += ";" + strconv.Itoa(site_id)
			network.MsgTo("R" + strconv.Itoa(elu) + "," + list)
			etat = "R"
		}
	}
}

func GetElu() int{
	if etat == "N" {
		return elu
	}
	return -1
}


func getApt(idApt string) (int,int){

	splitIdApt :=  strings.Split(idApt,"-")
	id,_ := strconv.Atoi(splitIdApt[0])
	apt,_ := strconv.Atoi(splitIdApt[1])
	return id, apt
}

func MsgFrom(network string, address string) {
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

			oppCode := msg[0]

			if oppCode == 'A'{
				RcptAnnonce(msg[1:])
				fmt.Println("Rcpt Annonce: " + msg[1:])
			}else if oppCode == 'R'{
				RcptResultat(string(msg[1]),msg[3:])
				fmt.Println("Rcpt Resultat -> Elu: "+ string(msg[1]) + " et la liste: " + msg[1:])

			}
		}
	}
}
