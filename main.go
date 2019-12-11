package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"PRR_Lab03/network"
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

func main(){
	var siteId int
	var conf Conf
	LoadConfiguration(&conf)


	next := siteId + 1
	total := conf.NB_SITES

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
	fmt.Println(siteId)
	fmt.Println(next)
	fmt.Println(total)

	network.ConnTo(next,conf.SITES_ADDR[next%total],conf.SITES_ADDR[siteId])

}
