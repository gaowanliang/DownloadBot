package main

import (
	"fmt"
	"log"
	"v2/rpc"
)

// Config 是读入的配置文件的struct
type Config struct {
	Aria2Server string `json:"aria2-server"`
	Aria2Key    string `json:"aria2-key"`
	BotKey      string `json:"bot-key"`
	UserID      string `json:"user-id"`
	MaxIndex    int    `json:"max-index"`
	Sign        string `json:"sign"`
}

// Aria2Notifier is Aria2 websocket message
type Aria2Notifier struct{}

// OnDownloadStart will be sent when a download is started. The event is of type struct and it contains following keys. The value type is string.
func (Aria2Notifier) OnDownloadStart(events []rpc.Event) {
	log.Printf("%s started.", events)
	/*name := make([]string, 0)
	for _, file := range events {
		name = append(name, tellName(aria2Rpc.TellStatus(file.Gid)))
	}*/

	SuddenMessageChan <- fmt.Sprintf("%s 任务开始下载", events)
}

// OnDownloadPause will be sent when a download is paused. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadPause(events []rpc.Event) {
	log.Printf("%s paused.", events)
	SuddenMessageChan <- fmt.Sprintf("%s 任务暂停", events)
}

// OnDownloadStop will be sent when a download is stopped by the user. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadStop(events []rpc.Event) {
	log.Printf("%s stopped.", events)
	SuddenMessageChan <- fmt.Sprintf("%s 任务被停止", events)
}

// OnDownloadComplete will be sent when a download is complete. For BitTorrent downloads, this notification is sent when the download is complete and seeding is over. The event is the same struct of the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadComplete(events []rpc.Event) {
	log.Printf("%s completed.", events)
	SuddenMessageChan <- fmt.Sprintf("%s 任务已完成", events)
}

// OnDownloadError will be sent when a download is stopped due to an error. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadError(events []rpc.Event) {
	log.Printf("%s error.", events)
	SuddenMessageChan <- fmt.Sprintf("%s 任务发生错误", events)
}

// OnBtDownloadComplete will be sent when a torrent download is complete but seeding is still going on. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnBtDownloadComplete(events []rpc.Event) {
	log.Printf("bt %s completed.", events)
	//SuddenMessageChan <- fmt.Sprintf("BT 任务 %s 已完成", events)
}
