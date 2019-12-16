package main

import (
	"PRR_Lab03/algoCR"
	"PRR_Lab03/network"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Conf struct {
	NB_SITES   int
	SITES_ADDR []string
	APT_SITES []int
}

func LoadConfiguration(conf *Conf) {
	file, _ := os.Open("configuration/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&conf)
}

var fin chan bool
func main(){
	var siteId int
	var conf Conf
	LoadConfiguration(&conf)


	// parse command line args
	if len(os.Args) == 1 {
		log.Println("you have to provide a site id")
		return
	} else {
		siteId, _ = strconv.Atoi(os.Args[1])
		if !(0 <= siteId && siteId <= conf.NB_SITES) {
			log.Println("invalid site id")
			return
		}
	}
	go algoCR.MsgFrom("udp", conf.SITES_ADDR[siteId])


	algoCR.InitAlgo(siteId,conf.APT_SITES[siteId])
	network.InitNetwork(conf.NB_SITES,conf.SITES_ADDR,conf.APT_SITES,siteId)




	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nEnter text: [E (Election)]\n")

		// read the input of the user
		cmd, _ := reader.ReadString('\n')

		// if W, the site do an ask
		if cmd == "E\n" {
			algoCR.Election()
		} else {
			fmt.Println("unknown command " + cmd)
		}
	}

	<- fin

}
