package input

import (
	"DownloadBot/cmd/client"
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	"DownloadBot/tool/input/aria2"
	logger "DownloadBot/tool/zap"
	"strings"
)

type Tool struct {
	aria2.Aria2
}

var ToolApp = new(Tool)

func PauseTask(sign string) {
	switch config.InputToolMethod() {
	case 0:
		ToolApp.Aria2.Pause(sign)
		break
	default:
		logger.Error(i18nLoc.LocText("No input tool is selected"))
	}
}
func UnpauseTask(sign string) {
	switch config.InputToolMethod() {
	case 0:
		ToolApp.Aria2.Unpause(sign)
		break
	default:
		logger.Error(i18nLoc.LocText("No input tool is selected"))
	}
}

func ForceRemoveTask(sign string) {
	switch config.InputToolMethod() {
	case 0:
		ToolApp.Aria2.ForceRemove(sign)
		break
	default:
		logger.Error(i18nLoc.LocText("No input tool is selected"))
	}
}
func PauseAllTask() {
	switch config.InputToolMethod() {
	case 0:
		ToolApp.Aria2.PauseAll()
		break
	default:
		logger.Error(i18nLoc.LocText("No input tool is selected"))
	}
}
func UnpauseAllTask() {
	switch config.InputToolMethod() {
	case 0:
		ToolApp.Aria2.UnpauseAll()
		break
	default:
		logger.Error(i18nLoc.LocText("No input tool is selected"))
	}
}

func SuddenlyMsg() {
	for {
		a := <-aria2.SuddenlyMsgChan

		gid := a[2:18]
		if strings.Contains(a, "[{") {
			a = strings.Replace(a, gid, ToolApp.Aria2.TellName(gid), -1)
		}
		client.Method.SendSuddenMessage(config.GetSign() + ": " + a)
	}

}
