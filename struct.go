package main

import (
	"fmt"
	"log"
	"DownloadBot/src/rpc"
)

// Config 是读入的配置文件的struct
type Config struct {
	Aria2Server    string `json:"aria2-server"`
	Aria2Key       string `json:"aria2-key"`
	BotKey         string `json:"bot-key"`
	UserID         string `json:"user-id"`
	MaxIndex       int    `json:"max-index"`
	Sign           string `json:"sign"`
	Language       string `json:"language"`
	DownloadFolder string `json:"downloadFolder"`
	MoveFolder     string `json:"moveFolder"`
}

// Aria2DownloadMode is set aria2 download mode
type Aria2DownloadMode struct {
	TMMode string
}

// FilesInlineKeyboards set Files Inline Key boards
type FilesInlineKeyboards struct {
	GidAndName []map[string]string
	Data       string
}

// FunctionInlineKeyboards set Files Inline Key boards
type FunctionInlineKeyboards struct {
	Describe string
	Data     string
}

// Aria2Notifier is Aria2 websocket message
type Aria2Notifier struct{}

// OnDownloadStart will be sent when a download is started. The event is of type struct and it contains following keys. The value type is string.
func (Aria2Notifier) OnDownloadStart(events []rpc.Event) {
	log.Printf(locText("onDownloadStartDes"), events)
	/*name := make([]string, 0)
	for _, file := range events {
		name = append(name, tellName(aria2Rpc.TellStatus(file.Gid)))
	}*/

	SuddenMessageChan <- fmt.Sprintf(locText("onDownloadStartDes"), events)
	TMMessageChan <- events[0].Gid
}

// OnDownloadPause will be sent when a download is paused. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadPause(events []rpc.Event) {
	log.Printf(locText("onDownloadPauseDes"), events)
	SuddenMessageChan <- fmt.Sprintf(locText("onDownloadPauseDes"), events)
}

// OnDownloadStop will be sent when a download is stopped by the user. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadStop(events []rpc.Event) {
	log.Printf(locText("onDownloadStopDes"), events)
	SuddenMessageChan <- fmt.Sprintf(locText("onDownloadStopDes"), events)
}

// OnDownloadComplete will be sent when a download is complete. For BitTorrent downloads, this notification is sent when the download is complete and seeding is over. The event is the same struct of the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadComplete(events []rpc.Event) {
	log.Printf(locText("onDownloadCompleteDes"), events)
	SuddenMessageChan <- fmt.Sprintf(locText("onDownloadCompleteDes"), events)
}

// OnDownloadError will be sent when a download is stopped due to an error. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnDownloadError(events []rpc.Event) {
	log.Printf(locText("onDownloadErrorDes"), events)
	SuddenMessageChan <- fmt.Sprintf(locText("onDownloadErrorDes"), events)
}

// OnBtDownloadComplete will be sent when a torrent download is complete but seeding is still going on. The event is the same struct as the event argument of onDownloadStart() method.
func (Aria2Notifier) OnBtDownloadComplete(events []rpc.Event) {
	log.Printf(locText("onBtDownloadCompleteDes"), events)
	//SuddenMessageChan <- fmt.Sprintf("BT 任务 %s 已完成", events)
}
