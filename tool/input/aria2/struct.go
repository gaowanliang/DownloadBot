package aria2

import (
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/tool/input/aria2/rpc"
	logger "DownloadBot/tool/zap"
	"fmt"
)

// DownloadMode is set aria2 download mode
type DownloadMode struct {
	TMMode string
}

type Notifier struct {
}

// OnDownloadStart will be sent when a download is started. The event is of type struct, and it contains following keys. The value type is string.
func (Notifier) OnDownloadStart(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadStartDes"), events)
	/*name := make([]string, 0)
	for _, file := range events {
		name = append(name, tellName(aria2Rpc.TellStatus(file.Gid)))
	}*/
	SuddenlyMsgChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadStartDes"), events)

	TMMessageChan <- events[0].Gid
}

// OnDownloadPause will be sent when a download is paused. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnDownloadPause(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadPauseDes"), events)
	SuddenlyMsgChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadPauseDes"), events)
}

// OnDownloadStop will be sent when a download is stopped by the user. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnDownloadStop(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadStopDes"), events)
	SuddenlyMsgChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadStopDes"), events)
}

// OnDownloadComplete will be sent when a download is complete. For BitTorrent downloads, this notification is sent when the download is complete and seeding is over. The event is the same struct of the event argument of onDownloadStart() method.
func (Notifier) OnDownloadComplete(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadCompleteDes"), events)
	SuddenlyMsgChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadCompleteDes"), events)
}

// OnDownloadError will be sent when a download is stopped due to an error. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnDownloadError(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadErrorDes"), events)
	SuddenlyMsgChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadErrorDes"), events)
}

// OnBtDownloadComplete will be sent when a torrent download is complete but seeding is still going on. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnBtDownloadComplete(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onBtDownloadCompleteDes"), events)
	//SuddenMessageChan <- fmt.Sprintf("BT 任务 %s 已完成", events)
}
