package telegram

import (
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	"DownloadBot/tool/displayUtil/gotree"
	"DownloadBot/tool/input"
	"DownloadBot/tool/input/aria2"
	"DownloadBot/tool/input/aria2/rpc"
	"DownloadBot/tool/typeTrans"
	logger "DownloadBot/tool/zap"
	"fmt"
	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"sort"
	"strings"
	"time"
)

// SuddenMessageChan is a channel for reminding when an emergency message occurs, such as download start and download end
var SuddenMessageChan = make(chan string, 3)

// TMSelectMessageChan is a channel for reminding when you upload a torrent file or send a magnet(TM is Torrent/Magnet)
var TMSelectMessageChan = make(chan string, 3)

// SuddenMessage is a function for dealing when an emergency message occurs, such as download start and download end
func SuddenMessage(bot *tgBotApi.BotAPI) {
	for {
		a := <-SuddenMessageChan
		gid := a[2:18]
		if strings.Contains(a, "[{") {
			a = strings.Replace(a, gid, input.ToolApp.Aria2.TellName(gid), -1)
		}
		myID := typeTrans.Str2Int64(config.GetTelegramUserID())
		msg := tgBotApi.NewMessage(myID, a)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

// Aria2TMSelectMsg is a function for dealing when you upload a torrent file or send a magnet,this function will interrupt the download and output the file list of the torrent file or magnet (TM is Torrent/Magnet)
func Aria2TMSelectMsg(bot *tgBotApi.BotAPI) {
	var MessageID []int
	myID := typeTrans.Str2Int64(config.GetTelegramUserID())
	selectFileList := make([][2]int, 0)
	directoryTree := make(map[string]interface{}, 0)
	for {
		a := <-TMSelectMessageChan
		b := strings.Split(a, "~")
		// 0 is gid, 1 is control sign(Start/SelectAll...) / selected num
		gid := b[0]

		downloadFilesCount := 0
		if len(b) != 1 {
			if b[1] == "Start" { //click start
				input.ToolApp.Aria2.SetTMDownloadFilesAndStart(gid, selectFileList)
				for _, val := range MessageID {
					bot.Send(tgBotApi.NewDeleteMessage(myID, val))
				}
				MessageID = make([]int, 0)
				selectFileList = make([][2]int, 0)
				directoryTree = make(map[string]interface{}, 0)
				continue
			} else if b[1] == "cancel" {
				for _, val := range MessageID {
					bot.Send(tgBotApi.NewDeleteMessage(myID, val))
				}
				MessageID = make([]int, 0)
				selectFileList = make([][2]int, 0)
				directoryTree = make(map[string]interface{}, 0)
				input.ToolApp.Aria2.ForceRemove(gid)
				continue
			}
			for i := 0; i < len(selectFileList); i++ {
				if selectFileList[i][0] == 1 && selectFileList[i][1] >= 0 {
					downloadFilesCount++
				}
			}
			switch b[1] {
			case "selectAll": //click select all
				for i := 0; i < len(selectFileList); i++ {
					selectFileList[i][0] = 1
				}
			case "tmMode1": //click selectBiggestFile
				biggestFileIndex := input.ToolApp.Aria2.SelectBiggestFile(gid)
				for i := 0; i < len(selectFileList); i++ {
					if selectFileList[i][1] != biggestFileIndex {
						selectFileList[i][0] = 0
					} else {
						selectFileList[i][0] = 1
						if i != 0 && selectFileList[i-1][1] == -1 {
							// if there is only one file in the folder to which this file belongs, select it as well
							selectFileList[i-1][0] = 1
						}

					}
				}
			case "tmMode2": // click smart ....
				bigFilesIndex := input.ToolApp.Aria2.SelectBigFiles(gid)
				index2Original := make(map[int]int, 0)
				for i := 0; i < len(selectFileList); i++ {
					selectFileList[i][0] = 0
					if selectFileList[i][1] > 0 { // node do not need records
						index2Original[selectFileList[i][1]] = i
					}
				}
				for _, i := range bigFilesIndex {
					selectFileList[index2Original[i]][0] = 1
				}
			default:
				i := typeTrans.Str2Int(b[1]) - 1 // sequence num displayed is more than the subscript,so reduce it
				if downloadFilesCount > 1 {      // make sure that at least one file will be downloaded
					if selectFileList[i][1] >= 0 {
						if selectFileList[i][0] == 1 {
							selectFileList[i][0] = 0
						} else {
							selectFileList[i][0] = 1
						}
					} else { //if select node
						if selectFileList[i][0] == 1 {
							selectFileList[i][0] = 0
							for s := 1; s < selectFileList[i][1]*-1; s++ {
								selectFileList[i+s][0] = 0
								downloadFilesCount--
							}
							if downloadFilesCount < 1 {
								selectFileList[i][0] = 1
								for s := 1; s < selectFileList[i][1]*-1; s++ {
									selectFileList[i+s][0] = 1
									downloadFilesCount++
								}
							}
						} else {
							selectFileList[i][0] = 1
							for s := 1; s < selectFileList[i][1]*-1; s++ {
								selectFileList[i+s][0] = 1
							}
						}
					}
				}
			}

			for i, val := range selectFileList { // check weather all files under all nodes select status
				if val[1] < 0 {
					allSelected := true
					allNotSelected := true
					// log.Println(i+1,i+1+val[1]*-1-1,val[1],)
					for _, val1 := range selectFileList[i+1 : i+1+val[1]*-1-1] {
						if val1[0] == 0 {
							allSelected = false
						} else if val1[0] == 1 {
							allNotSelected = false
						}
						if !allSelected && !allNotSelected {
							break
						}
					}
					if !allSelected && !allNotSelected { // partially selected, partially not be selected
						selectFileList[i][0] = -1
					} else if allSelected {
						selectFileList[i][0] = 1
					} else {
						selectFileList[i][0] = 0
					}

				}
			}
		} else {
			fileList := input.ToolApp.Aria2.FormatTMFiles(gid)
			index := 1
			for i, file := range fileList {
				pathClass(fmt.Sprintf("%s|%s|%d", file[0], file[1], i+1), &directoryTree)
				//logger.Info("%s|%s|%d\n", file[0], file[1], i+1)
				index++
			}
		}

		text := fmt.Sprintf("%s %s\n", input.ToolApp.Aria2.TellName(gid), i18nLoc.LocText("fileDirectoryIsAsFollows"))
		Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)

		fileListGoTree, _, _ := generateGoTree(directoryTree, 0, &selectFileList)
		fileListTreeLine := strings.Split(fileListGoTree[0].Print(), "\n")
		fileListTreeLineCount := len(fileListTreeLine)
		characterCount := len(text)
		r, err := regexp.Compile(`[✅⬜〰](\d+)`)
		dropErr(err)
		startAndEndIndex := []int{1, 0}
		msgCount := 0

		for i, line := range fileListTreeLine {
			if fileListTreeLineCount != i+1 && characterCount+len(fileListTreeLine[i+1]) > 4096 {

				msg := tgBotApi.NewMessage(myID, text)
				msgCount++
				//lastFilesInfo = fileList
				index := 0
				for j := startAndEndIndex[0]; j <= startAndEndIndex[1]; j++ {
					inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(j), gid+"~"+fmt.Sprint(j)+":6"))
					index++
					if index%7 == 0 {
						Keyboards = append(Keyboards, inlineKeyBoardRow)
						inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
					}

				}
				if len(inlineKeyBoardRow) != 0 {
					Keyboards = append(Keyboards, inlineKeyBoardRow)
				}
				if msgCount > len(MessageID) {
					msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(Keyboards...)
					res, err := bot.Send(msg)
					MessageID = append(MessageID, res.MessageID)
					dropErr(err)
				} else {
					newMsg := tgBotApi.NewEditMessageTextAndMarkup(myID, MessageID[msgCount-1], text, tgBotApi.NewInlineKeyboardMarkup(Keyboards...))
					bot.Send(newMsg)
				}

				Keyboards = make([][]tgBotApi.InlineKeyboardButton, 0)
				inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)

				//lastGid = gid
				characterCount = 0
				text = "·"
				startAndEndIndex[0] = startAndEndIndex[1] + 1
			}

			if text == "·" {
				text += line[1:] + "\n"
			} else {
				text += line + "\n"
			}
			characterCount += len(line + "\n")
			res := r.FindStringSubmatch(line)
			if res != nil {
				startAndEndIndex[1] = typeTrans.Str2Int(res[1])
			}
		}

		// log.Println(text)
		text += i18nLoc.LocText("pleaseSelectTheFileYouWantToDownload")
		if msgCount > 2 {
			text += "\n" + i18nLoc.LocText("tmFileTooMany")
		}
		index := 0
		for j := startAndEndIndex[0]; j <= startAndEndIndex[1]; j++ {
			inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(j), gid+"~"+fmt.Sprint(j)+":6"))
			index++
			if index%7 == 0 {
				Keyboards = append(Keyboards, inlineKeyBoardRow)
				inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
			}

		}

		if len(inlineKeyBoardRow) != 0 {
			Keyboards = append(Keyboards, inlineKeyBoardRow)
		}
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(i18nLoc.LocText("selectAll"), gid+"~selectAll"+":7"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(i18nLoc.LocText("cancel"), gid+"~cancel"+":7"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(i18nLoc.LocText("tmMode1"), gid+"~tmMode1"+":7"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(i18nLoc.LocText("tmMode2"), gid+"~tmMode2"+":7"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(i18nLoc.LocText("startDownload"), gid+"~Start"+":7"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)

		//myID, err := strconv.ParseInt(info.UserID, 10, 64)
		//dropErr(err)

		//lastFilesInfo = fileList
		if len(MessageID) < msgCount+1 {
			msg := tgBotApi.NewMessage(myID, text)
			msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(Keyboards...)
			res, err := bot.Send(msg)
			dropErr(err)
			MessageID = append(MessageID, res.MessageID)
		} else {
			newMsg := tgBotApi.NewEditMessageTextAndMarkup(myID, MessageID[msgCount], text, tgBotApi.NewInlineKeyboardMarkup(Keyboards...))
			bot.Send(newMsg)
		}

	}
}

const FileMarker = "<files>"

//pathClass is a function that can generate a directory tree structure when giving a file path list. val is path, trunk is the directory structure
func pathClass(branch string, trunk *map[string]interface{}) {
	parts := strings.SplitN(branch, "/", 2)
	if len(parts) == 1 {
		if (*trunk)[FileMarker] == nil {
			fileNameList := make([]string, 0)
			fileNameList = append(fileNameList, parts[0])
			(*trunk)[FileMarker] = fileNameList
		} else {
			fileNameList := (*trunk)[FileMarker].([]string)
			fileNameList = append(fileNameList, parts[0])
			(*trunk)[FileMarker] = fileNameList
		}
	} else {
		node, other := parts[0], parts[1]
		if _, ok := (*trunk)[node]; !ok {
			(*trunk)[node] = map[string]interface{}{}
		}
		_temp := (*trunk)[node].(map[string]interface{})
		pathClass(other, &_temp)
	}
}

var activeRefreshControl = 0

// activeRefresh refresh the download info
func activeRefresh(chatMsgID int, bot *tgBotApi.BotAPI, ticker *time.Ticker, flag int) {
	var MessageID = 0
	myID := typeTrans.Str2Int64(config.GetTelegramUserID())

	refreshPath := func(MessageID int, myID int64, bot *tgBotApi.BotAPI, ticker *time.Ticker) int {
		res := input.ToolApp.Aria2.FormatTellActive()
		//log.Println(res, len(res))
		text := ""
		if res != "" {
			text = res
		} else {
			text = i18nLoc.LocText("noActiveTask")
		}
		if MessageID == 0 {
			msg := tgBotApi.NewMessage(myID, text)
			msg.ParseMode = "Markdown"
			res, err := bot.Send(msg)
			dropErr(err)
			if text == i18nLoc.LocText("noActiveTask") {
				ticker.Stop()
				return -1
			} else {
				return res.MessageID
			}
		} else {
			if text == i18nLoc.LocText("noActiveTask") {
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
				ticker.Stop()
				return -1
			} else {
				newMsg := tgBotApi.NewEditMessageText(myID, MessageID, text)
				newMsg.ParseMode = "Markdown"
				bot.Send(newMsg)
				return newMsg.MessageID
			}

		}
	}

	for {
		if activeRefreshControl != flag {
			if MessageID != 0 {
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
			}
			ticker.Stop()
			msgToDelete := tgBotApi.DeleteMessageConfig{
				ChatID:    myID,
				MessageID: chatMsgID,
			}
			_, _ = bot.Request(msgToDelete)
			return
		} else {
			if MessageID != 0 {
				select {
				case _ = <-ticker.C:
					MessageID = refreshPath(MessageID, myID, bot, ticker)
					if MessageID == -1 {
						return
					}
				}
			} else {
				MessageID = refreshPath(MessageID, myID, bot, ticker)
				if MessageID == -1 {
					return
				}
			}
		}
	}
}

//generateGoTree is a function that receive the directory tree structure generated by pathClass(),and file list that user want to select ,to generate goTree,return both goTree and the list of selected files
func generateGoTree(m map[string]interface{}, index int, selectFileList *[][2]int) ([]goTree.Tree, int, int) {
	var artist []goTree.Tree
	allFilesOwnedByNodeCount := 0
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		var _artist goTree.Tree
		filesOwnedByNodeCount := 0
		switch vv := m[k].(type) {
		case map[string]interface{}:
			//fmt.Println(k, "is an map:")
			var isEmpty = false
			var nodeIndex = index
			index++
			if len(*selectFileList) < index {
				_selectAndOriginalNum := [2]int{1, 0} //node has no original num
				*selectFileList = append(*selectFileList, _selectAndOriginalNum)
				isEmpty = true
			}
			if (*selectFileList)[index-1][0] == 1 {
				_artist = goTree.New(fmt.Sprintf("✅%d: %s", index, k))
			} else if (*selectFileList)[index-1][0] == 0 {
				_artist = goTree.New(fmt.Sprintf("⬜%d: %s", index, k))
			} else {
				_artist = goTree.New(fmt.Sprintf("〰%d: %s", index, k))
			}

			if vs, ok := vv[FileMarker]; ok {
				ls := vs.([]string)
				//flag=true
				for _, u := range ls {
					index++
					filesOwnedByNodeCount++
					_fileInfo := strings.Split(u, "|")
					//0 is file name,1 is file size, 2 is original num

					if isEmpty {
						_selectAndOriginalNum := [2]int{1, typeTrans.Str2Int(_fileInfo[2])}
						*selectFileList = append(*selectFileList, _selectAndOriginalNum)
					}

					if (*selectFileList)[index-1][0] == 1 {
						_artist.Add(fmt.Sprintf("✅%d: %s|%s", index, _fileInfo[0], _fileInfo[1]))
					} else {
						_artist.Add(fmt.Sprintf("⬜%d: %s|%s", index, _fileInfo[0], _fileInfo[1]))
					}
				}
			}
			res, _index, _filesOwnedByNodeCount := generateGoTree(vv, index, selectFileList)
			index = _index
			filesOwnedByNodeCount += _filesOwnedByNodeCount + 1
			if res != nil {
				for _, vk := range res {
					_artist.AddTree(vk)
				}
			}
			artist = append(artist, _artist)
			(*selectFileList)[nodeIndex][1] = filesOwnedByNodeCount * -1 // node has no original num, so we can set the second to the number of file subordinate to this node,but to prevent confusion with original num, we will take a negative
		}
		allFilesOwnedByNodeCount += filesOwnedByNodeCount
	}
	return artist, index, allFilesOwnedByNodeCount
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

	SuddenMessageChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadStartDes"), events)
	aria2.TMMessageChan <- events[0].Gid
}

// OnDownloadPause will be sent when a download is paused. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnDownloadPause(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadPauseDes"), events)
	SuddenMessageChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadPauseDes"), events)
}

// OnDownloadStop will be sent when a download is stopped by the user. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnDownloadStop(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadStopDes"), events)
	SuddenMessageChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadStopDes"), events)
}

// OnDownloadComplete will be sent when a download is complete. For BitTorrent downloads, this notification is sent when the download is complete and seeding is over. The event is the same struct of the event argument of onDownloadStart() method.
func (Notifier) OnDownloadComplete(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadCompleteDes"), events)
	SuddenMessageChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadCompleteDes"), events)
}

// OnDownloadError will be sent when a download is stopped due to an error. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnDownloadError(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onDownloadErrorDes"), events)
	SuddenMessageChan <- fmt.Sprintf(i18nLoc.LocText("onDownloadErrorDes"), events)
}

// OnBtDownloadComplete will be sent when a torrent download is complete but seeding is still going on. The event is the same struct as the event argument of onDownloadStart() method.
func (Notifier) OnBtDownloadComplete(events []rpc.Event) {
	logger.Info(i18nLoc.LocText("onBtDownloadCompleteDes"), events)
	//SuddenMessageChan <- fmt.Sprintf("BT 任务 %s 已完成", events)
}
