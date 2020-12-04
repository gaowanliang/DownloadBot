package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var info Config

func dropErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func init() {
	filePtr, err := os.Open("./config.json")
	dropErr(err)
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	dropErr(err)

	locLan(info.Language)
	log.Print(locText("configCompleted"))
}

func main() {
	var wg sync.WaitGroup
	aria2Load()
	wg.Add(1)
	go tgBot(info.BotKey, &wg)
	wg.Wait()
	defer aria2Rpc.Close()
}
