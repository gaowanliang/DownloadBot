package clientManage

import "fmt"

type ClientInfo struct {
	ClientName string
	AlivePoint int // when client not send message to server, server will minus 1, if 0, delete client
}

var ClientList = make([]ClientInfo, 0)

func ClientSearch(client string) bool {
	for _, v := range ClientList {
		if v.ClientName == client {
			return true
		}
	}
	return false
}

func ShowClientList() (string, int) {
	var clientListStr string
	for i, v := range ClientList {
		clientListStr += fmt.Sprintf("%d:%v\n", i+1, v.ClientName)
	}
	return clientListStr, len(ClientList)
}

func GetClientName(index int) string {
	return ClientList[index].ClientName
}
