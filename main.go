package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"v2/rpc"
)

var info Config

func dropErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func configLoad() Config {
	filePtr, err := os.Open("./config.json")
	dropErr(err)
	defer filePtr.Close()
	var info Config
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	dropErr(err)
	log.Print("Configuration information loading completed!")
	return info
}

func testAria2(rpcc rpc.Client) {
	/*const targetURL = "https://nodejs.org/dist/index.json"

	g, err := rpcc.AddURI([]string{targetURL})
	if err != nil {
		log.Panic(err)
	}
	println(g)
	if _, err = rpcc.TellActive(); err != nil {
		log.Panic(err)
	}
	if _, err = rpcc.PauseAll(); err != nil {
		log.Panic(err)
	}*/
	info, err := rpcc.TellActive()
	dropErr(err)
	log.Println(info)
}

func main() {
	info = configLoad()
	var wg sync.WaitGroup
	aria2Load()
	wg.Add(1)
	go tgBot(info.BotKey, &wg)
	//testAria2(aria2Rpc)
	wg.Wait()
	defer aria2Rpc.Close()
}
