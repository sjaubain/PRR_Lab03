package algoCR

import (
	"PRR_Lab03/network"
	"fmt"
	"strconv"
	"strings"
)

var (
	site_id int
	apt     int
	idApt   string
	etat    string
	elu     int
)

func InitAlgo(identifiant int, apti int) {
	site_id = identifiant
	apt = apti
	idApt = strconv.Itoa(site_id) + "-" + strconv.Itoa(apt)
}

func Election() {
	network.MsgTo("A" + idApt)
	etat = "A"
}

func RcptAnnonce(list string) {
	var msg string
	if strings.Contains(list, idApt) {
		tabList := strings.Split(list, ";")
		aptMax := 0
		elu = 0
		for i := range tabList {
			id, apt := getApt(tabList[i])
			if apt > aptMax {
				aptMax = apt
				elu = id
			}
		}
		msg = "R" + strconv.Itoa(elu) + "," + strconv.Itoa(site_id)
		etat = "R"
		fmt.Println("Envoie Resultat avec elu !")
	} else {
		list += ";" + idApt
		msg = "A" + list
		etat = "A"
	}
	network.MsgTo(msg)

}

func RcptResultat(i string, list string) {
	tabProc := strings.Split(list, ";")
	for j := range tabProc {
		if strconv.Itoa(site_id) == tabProc[j] {
			fmt.Println("Fin - le processus " + strconv.Itoa(elu) + " est l'elu")
			etat = "N"
		} else if etat == "R" && strconv.Itoa(elu) != i {
			fmt.Println("Lance une nouvelle election car contradiction")
			Election()
		} else if etat == "A" {
			elu, _ = strconv.Atoi(i)
			fmt.Println("Rcpt Resultat -> Elu: " + i + " et la liste: " + list)
			list += ";" + strconv.Itoa(site_id)
			etat = "R"
			network.MsgTo("R" + strconv.Itoa(elu) + "," + list)

		}
	}
}

func GetElu() int {
	if etat == "N" {
		return elu
	}
	return -1
}

func getApt(idApt string) (int, int) {

	splitIdApt := strings.Split(idApt, "-")
	id, _ := strconv.Atoi(splitIdApt[0])
	apt, _ := strconv.Atoi(splitIdApt[1])
	return id, apt
}

func MsgHandle(net string, add string) {
	for{
		msg := network.MsgFrom(net, add)
		oppCode := msg[0]
		if oppCode == 'A' {
			fmt.Println("Rcpt Annonce: " + msg[1:])
			RcptAnnonce(msg[1:])
		} else if oppCode == 'R' {
			RcptResultat(string(msg[1]), msg[3:])
		}
	}

}

