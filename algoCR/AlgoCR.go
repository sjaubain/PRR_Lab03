package algoCR

import (
	"PRR_Lab03/network"
	"fmt"
	"strconv"
	"strings"
)

var (
	site_id    int       // identifiant du site
	apt        int       // l'aptitude du site
	idApt      string    // le coupe (id-apt)
	etat       string    // Etat (R | A | N)
	elu        int       // identifiant de l'elu courant
	hdlMsgDone chan bool // channel mutex pour protéger la variable
	// elu et ne traiter qu'une réception à la fois
)

func InitAlgo(identifiant int, apti int) {
	site_id = identifiant
	apt = apti
	idApt = strconv.Itoa(site_id) + "-" + strconv.Itoa(apt)

	hdlMsgDone = make(chan bool, 1)
	hdlMsgDone <- true // channel mutex pour protéger l'accès aux variables
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
			network.MsgTo("F" + idApt) // on refait le tour pour leur donner l elu
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

func RcptFin(list string){
	if etat != "A"{
		etat = "N"
		if strings.Contains(list, idApt) {
			fmt.Println("Tout les sites ont le même elu")
		} else {
			list += ";" + idApt
			network.MsgTo(list)
		}
	}
}

func GetElu() int {
	if etat == "N" {
		return elu
	}
	return -1
}

// nous permet de parser le msg, retourne l id du site et son aptitude
func getApt(idApt string) (int, int) {

	splitIdApt := strings.Split(idApt, "-")
	id, _ := strconv.Atoi(splitIdApt[0])
	apt, _ := strconv.Atoi(splitIdApt[1])
	return id, apt
}

func MsgHandle(net string, add string) {
	for {
		msg := network.MsgFrom(net, add)

		// exécution en parallèle pour ne pas bloquer le MsgFrom
		go func() {

			<-hdlMsgDone

			oppCode := msg[0]
			if oppCode == 'A' {
				fmt.Println("Rcpt Annonce: " + msg[1:])
				RcptAnnonce(msg[1:])
			} else if oppCode == 'R' {
				RcptResultat(string(msg[1]), msg[3:])
			} else if oppCode == 'F' {
				fmt.Println("Rcpt elu: " + msg[1:])
				RcptFin(msg)
			}

			hdlMsgDone <- true
		}()
	}
}
