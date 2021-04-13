package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SuddenMessageChan receive active requests from WebSocket
var SuddenMessageChan = make(chan string, 3)

// TMSelectMessageChan receive active requests from WebSocket
var TMSelectMessageChan = make(chan string, 3)

var FileControlChan = make(chan string, 3)

func setCommands(bot *tgBotApi.BotAPI) {
	_ = bot.SetMyCommands([]tgBotApi.BotCommand{
		{
			Command:     "start",
			Description: locText("tgCommandStartDes"),
		}, {
			Command:     "myid",
			Description: locText("tgCommandMyIDDes"),
		},
	})
}

// SuddenMessage receive active requests from WebSocket
func SuddenMessage(bot *tgBotApi.BotAPI) {
	for {
		a := <-SuddenMessageChan
		gid := a[2:18]
		if strings.Contains(a, "[{") {
			a = strings.Replace(a, gid, tellName(aria2Rpc.TellStatus(gid)), -1)
		}
		myID, err := strconv.ParseInt(info.UserID, 10, 64)
		dropErr(err)
		msg := tgBotApi.NewMessage(myID, a)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

// TMSelectMessage receive active requests from WebSocket
func TMSelectMessage(bot *tgBotApi.BotAPI) {
	var MessageID = 0
	var lastGid = ""
	var lastFilesInfo [][]string
	myID := toInt64(info.UserID)
	for {
		a := <-TMSelectMessageChan
		b := strings.Split(a, "~")
		gid := b[0]
		var fileList [][]string
		downloadFilesCount := 0
		if len(b) != 1 {
			if b[1] == "Start" {
				setTMDownloadFilesAndStart(gid, lastFilesInfo)
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
				continue
			}
			for i := 0; i < len(lastFilesInfo); i++ {
				if lastFilesInfo[i][2] == "true" {
					downloadFilesCount++
				}
			}
			switch b[1] {
			case "selectAll":
				for i := 0; i < len(lastFilesInfo); i++ {
					lastFilesInfo[i][2] = "true"
				}
			case "invert":
				for i := 0; i < len(lastFilesInfo); i++ {
					if lastFilesInfo[i][2] == "true" {
						lastFilesInfo[i][2] = "false"
					} else {
						lastFilesInfo[i][2] = "true"
					}
				}
				if downloadFilesCount == len(lastFilesInfo) {
					lastFilesInfo[0][2] = "true"
				}
			case "tmMode1":
				biggestFileIndex := selectBigestFile(gid)
				for i := 0; i < len(lastFilesInfo); i++ {
					if i != biggestFileIndex {
						lastFilesInfo[i][2] = "false"
					} else {
						lastFilesInfo[i][2] = "true"
					}

				}
			case "tmMode2":
				for i := 0; i < len(lastFilesInfo); i++ {
					lastFilesInfo[i][2] = "false"
				}
				bigFilesIndex := selectBigFiles(gid)
				for _, i := range bigFilesIndex {
					lastFilesInfo[i][2] = "true"
				}
			default:
				i := toInt(b[1])
				i--
				if lastFilesInfo[i][2] == "true" && downloadFilesCount > 1 {
					lastFilesInfo[i][2] = "false"
				} else {
					lastFilesInfo[i][2] = "true"
				}

			}
			fileList = lastFilesInfo
		} else {
			fileList = formatTMFiles(gid)
		}

		text := fmt.Sprintf("%s %s\n", tellName(aria2Rpc.TellStatus(gid)), locText("fileDirectoryIsAsFollows"))
		Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
		index := 1

		for i, file := range fileList {
			isSelect := "⬜"
			if file[2] == "true" {
				isSelect = "✅"
			}
			text += fmt.Sprintf("%s %d: %s    %s\n", isSelect, i+1, file[0], file[1])
			inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(index), gid+"~"+fmt.Sprint(index)+":6"))
			if index%7 == 0 {
				Keyboards = append(Keyboards, inlineKeyBoardRow)
				inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
			}
			index++
		}
		text += locText("pleaseSelectTheFileYouWantToDownload")
		if len(inlineKeyBoardRow) != 0 {
			Keyboards = append(Keyboards, inlineKeyBoardRow)
		}
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("selectAll"), gid+"~selectAll"+":7"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("invert"), gid+"~invert"+":7"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("tmMode1"), gid+"~tmMode1"+":7"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("tmMode2"), gid+"~tmMode2"+":7"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("startDownload"), gid+"~Start"+":7"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
		//myID, err := strconv.ParseInt(info.UserID, 10, 64)
		//dropErr(err)
		msg := tgBotApi.NewMessage(myID, text)
		lastFilesInfo = fileList
		if lastGid != gid {
			msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(Keyboards...)
			res, err := bot.Send(msg)
			dropErr(err)
			MessageID = res.MessageID
			lastGid = gid
		} else {
			newMsg := tgBotApi.NewEditMessageTextAndMarkup(myID, MessageID, text, tgBotApi.NewInlineKeyboardMarkup(Keyboards...))
			bot.Send(newMsg)
		}

	}
}

func removeFiles(chatMsgID int, bot *tgBotApi.BotAPI) {
	s := <-FileControlChan
	if s == "file" {
		FileControlChan <- "file"
	}
	var MessageID = 0
	var filesSelect = make(map[int]bool)
	fileList, _ := GetAllFile(info.DownloadFolder)
	myID := toInt64(info.UserID)
	_, _ = bot.DeleteMessage(tgBotApi.DeleteMessageConfig{
		ChatID:    myID,
		MessageID: chatMsgID,
	})
	if len(fileList) == 1 {
		bot.Send(tgBotApi.NewMessage(myID, locText("noFilesFound")))
		return
	}
	deleteFiles := make([]string, 0)
	for {
		a := <-FileControlChan
		if a == "close" {
			//tgBotApi.NewDeleteMessage(myID, MessageID)
			bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
			return
		}
		b := strings.Split(a, "~")
		fileTree := ""
		if len(b) == 1 {
			filesSelect = make(map[int]bool)
			for i := 1; i <= len(fileList); i++ {
				filesSelect[i] = true
			}
			fileTree, filesSelect, deleteFiles = printFolderTree(info.DownloadFolder, filesSelect, "0")
		} else {
			if b[1] == "cancel" {
				tgBotApi.NewDeleteMessage(myID, MessageID)
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
				return
			} else if b[1] == "Delete" {
				RemoveFiles(deleteFiles)
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
				bot.Send(tgBotApi.NewMessage(myID, locText("filesDeletedSuccessfully")))
				return
			}
			fileTree, filesSelect, deleteFiles = printFolderTree(info.DownloadFolder, filesSelect, b[1])
		}

		text := fmt.Sprintf("%s %s\n", info.DownloadFolder, locText("fileDirectoryIsAsFollows")) + fileTree
		Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
		index := 1
		for _, _ = range fileList {
			inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(index), "file~"+fmt.Sprint(index)+":8"))
			if index%7 == 0 {
				Keyboards = append(Keyboards, inlineKeyBoardRow)
				inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
			}
			index++
		}
		text += locText("pleaseSelectTheFileYouWantToDelete")
		if len(inlineKeyBoardRow) != 0 {
			Keyboards = append(Keyboards, inlineKeyBoardRow)
		}
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("selectAll"), "file~selectAll"+":9"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("invert"), "file~invert"+":9"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("confirmDelete"), "file~Delete"+":9"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("cancel"), "file~cancel"+":9"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)

		msg := tgBotApi.NewMessage(myID, text)
		if MessageID == 0 {
			msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(Keyboards...)
			res, err := bot.Send(msg)
			dropErr(err)
			MessageID = res.MessageID

		} else {
			newMsg := tgBotApi.NewEditMessageTextAndMarkup(myID, MessageID, text, tgBotApi.NewInlineKeyboardMarkup(Keyboards...))
			bot.Send(newMsg)
		}
	}
}

func copyFiles(chatMsgID int, bot *tgBotApi.BotAPI) {
	s := <-FileControlChan
	if s == "file" {
		FileControlChan <- "file"
	}
	var MessageID = 0
	var filesSelect = make(map[int]bool)
	fileList, _ := GetAllFile(info.DownloadFolder)
	myID := toInt64(info.UserID)
	_, _ = bot.DeleteMessage(tgBotApi.DeleteMessageConfig{
		ChatID:    myID,
		MessageID: chatMsgID,
	})
	if len(fileList) == 1 {
		bot.Send(tgBotApi.NewMessage(myID, locText("noFilesFound")))
		return
	}
	copyFiles := make([]string, 0)
	for {
		a := <-FileControlChan
		if a == "close" {
			//tgBotApi.NewDeleteMessage(myID, MessageID)
			bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
			return
		}
		b := strings.Split(a, "~")
		fileTree := ""
		if len(b) == 1 {
			filesSelect = make(map[int]bool)
			for i := 1; i <= len(fileList); i++ {
				filesSelect[i] = true
			}
			fileTree, filesSelect, copyFiles = printFolderTree(info.DownloadFolder, filesSelect, "0")
		} else {
			if b[1] == "cancel" {
				//tgBotApi.NewDeleteMessage(myID, MessageID)
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
				return
			} else if b[1] == "Copy" {
				CopyFiles(copyFiles)
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
				bot.Send(tgBotApi.NewMessage(myID, locText("filesCopySuccessfully")))
				return
			}
			fileTree, filesSelect, copyFiles = printFolderTree(info.DownloadFolder, filesSelect, b[1])
		}

		text := fmt.Sprintf("%s %s\n", info.DownloadFolder, locText("fileDirectoryIsAsFollows")) + fileTree
		Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
		index := 1
		for _, _ = range fileList {
			inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(index), "file~"+fmt.Sprint(index)+":8"))
			if index%7 == 0 {
				Keyboards = append(Keyboards, inlineKeyBoardRow)
				inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
			}
			index++
		}
		text += locText("pleaseSelectTheFileYouWantToCopy")
		if len(inlineKeyBoardRow) != 0 {
			Keyboards = append(Keyboards, inlineKeyBoardRow)
		}
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("selectAll"), "file~selectAll"+":9"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("invert"), "file~invert"+":9"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
		inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("confirmCopy"), "file~Copy"+":9"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("cancel"), "file~cancel"+":9"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)

		msg := tgBotApi.NewMessage(myID, text)
		if MessageID == 0 {
			msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(Keyboards...)
			res, err := bot.Send(msg)
			dropErr(err)
			MessageID = res.MessageID

		} else {
			newMsg := tgBotApi.NewEditMessageTextAndMarkup(myID, MessageID, text, tgBotApi.NewInlineKeyboardMarkup(Keyboards...))
			bot.Send(newMsg)
		}
	}
}

func uploadFiles(chatMsgID int, chatMsg string, bot *tgBotApi.BotAPI) {
	_, err := os.Stat("info")
	if err != nil {
		err = os.MkdirAll("info", os.ModePerm)
		dropErr(err)
	}
	s := <-FileControlChan
	if s != "close" {
		FileControlChan <- s
	}
	var MessageID = 0
	myID := toInt64(info.UserID)
	_, _ = bot.DeleteMessage(tgBotApi.DeleteMessageConfig{
		ChatID:    myID,
		MessageID: chatMsgID,
	})
	for {
		a := <-FileControlChan
		if a == "close" {
			bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
			return
		}
		b := strings.Split(a, "~")
		text := ""
		//log.Println(b)
		Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
		if len(b) != 1 {
			if b[1] == "cancel" { //当用户点击取消或者流程结束后，删除消息
				tgBotApi.NewDeleteMessage(myID, MessageID)
				bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
				return
			} else {
				switch b[0] {
				case "onedrive":
					switch b[1] {
					case "new": //向用户发送授权地址
						text = fmt.Sprintf(
							`%s https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=ad5e65fd-856d-4356-aefc-537a9700c137&response_type=code&redirect_uri=http://localhost/onedrive-login&response_mode=query&scope=offline_access%%20User.Read%%20Files.ReadWrite.All`,
							locText("oneDriveGetAccess"),
						)
					case "create": //接收用户返回的授权地址，并进行处理
						mail := getNewOneDriveInfo(chatMsg)
						text = locText("oneDriveOAuthFileCreateSuccess") + mail
					}
				case "googleDrive":
					switch b[1] {
					case "new": //向用户发送授权地址
						text = fmt.Sprintf(
							`%s %s`,
							locText("googleDriveGetAccess"), getGoogleDriveAuthCodeURL(),
						)
					case "create": //接收用户返回的授权地址，并进行处理
						mail := getNewGoogleDriveInfo(chatMsg)
						text = locText("googleDriveOAuthFileCreateSuccess") + mail
					}
				case "odInfo": //获取已登录的OneDrive用户
					uploadDFToOneDrive("./info/onedrive/" + b[1])
				case "gdInfo": // 获取已登录的 Google drive 用户
					uploadDFToGoogleDrive("./info/googleDrive/" + b[1])
				default:
					switch b[1] {
					case "1": // 选择OneDrive
						createDriveInfoFolder("./info/onedrive")

						dir, _ := ioutil.ReadDir("./info/onedrive")
						if len(dir) == 0 {
							text = locText("noOneDriveInfo")
							inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("yes"), "onedrive~new"+":9"))
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("no"), "onedrive~cancel"+":9"))
							Keyboards = append(Keyboards, inlineKeyBoardRow)
						} else {
							text = locText("accountsAreCurrentlyLogin")
							files, _text, index := getAuthInfoJson("./info/onedrive/")
							text += _text

							inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
							index = 1
							for _, name := range files {
								inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(index), "odInfo~"+name+":9"))
								if index%7 == 0 {
									Keyboards = append(Keyboards, inlineKeyBoardRow)
									inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
								}
								index++
							}
							if len(inlineKeyBoardRow) != 0 {
								Keyboards = append(Keyboards, inlineKeyBoardRow)
							}
							inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("createNewAcc"), "onedrive~new"+":9"))
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("cancel"), "upload~cancel"+":9"))
							Keyboards = append(Keyboards, inlineKeyBoardRow)
							text += locText("selectAccount")
						}

					case "2": // 选择Google drive
						createDriveInfoFolder("./info/googleDrive")
						dir, _ := ioutil.ReadDir("./info/googleDrive")
						if len(dir) == 0 {
							text = locText("noGoogleDriveInfo")
							inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("yes"), "googleDrive~new"+":9"))
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("no"), "googleDrive~cancel"+":9"))
							Keyboards = append(Keyboards, inlineKeyBoardRow)
						} else {
							text = locText("accountsAreCurrentlyLogin")
							files, _text, index := getAuthInfoJson("./info/googleDrive/")
							text += _text
							inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
							index = 1
							for _, name := range files {
								inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(index), "gdInfo~"+name+":9"))
								if index%7 == 0 {
									Keyboards = append(Keyboards, inlineKeyBoardRow)
									inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
								}
								index++
							}
							if len(inlineKeyBoardRow) != 0 {
								Keyboards = append(Keyboards, inlineKeyBoardRow)
							}
							inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("createNewAcc"), "googleDrive~new"+":9"))
							inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("cancel"), "upload~cancel"+":9"))
							Keyboards = append(Keyboards, inlineKeyBoardRow)
							text += locText("selectAccount")
						}
					case "Upload": //弃用项
						//CopyFiles(copyFiles)
						bot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
						bot.Send(tgBotApi.NewMessage(myID, locText("filesUploadSuccessfully")))
						return
					}
				}
			}

		} else { // 生成 选择网盘菜单
			text = fmt.Sprintf(locText("chooseDrive"))

			inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
			index := 1
			for i := 0; i < 4; i++ {
				inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(index), "upload~"+fmt.Sprint(index)+":9"))
				if index%7 == 0 {
					Keyboards = append(Keyboards, inlineKeyBoardRow)
					inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
				}
				index++
			}
			if len(inlineKeyBoardRow) != 0 {
				Keyboards = append(Keyboards, inlineKeyBoardRow)
			}
			inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
			inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("cancel"), "upload~cancel"+":9"))
			Keyboards = append(Keyboards, inlineKeyBoardRow)
		}

		msg := tgBotApi.NewMessage(myID, text)
		msg.ParseMode = "Markdown"
		if MessageID == 0 {
			if len(Keyboards) != 0 { // 非首次生成选择消息
				msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(Keyboards...)
			}
			res, err := bot.Send(msg)
			dropErr(err)
			MessageID = res.MessageID

		} else { // 首次生成选择消息
			if len(Keyboards) != 0 {
				newMsg := tgBotApi.NewEditMessageTextAndMarkup(myID, MessageID, text, tgBotApi.NewInlineKeyboardMarkup(Keyboards...))
				bot.Send(newMsg)
			} else {
				newMsg := tgBotApi.NewEditMessageText(myID, MessageID, text)
				bot.Send(newMsg)
			}

		}
	}
}

var activeRefreshControl = 0

// activeRefresh 刷新下载进度显示
func activeRefresh(chatMsgID int, bot *tgBotApi.BotAPI, ticker *time.Ticker, flag int) {
	var MessageID = 0
	myID := toInt64(info.UserID)

	refreshPath := func(MessageID int, myID int64, bot *tgBotApi.BotAPI, ticker *time.Ticker) int {
		res := formatTellSomething(aria2Rpc.TellActive())
		//log.Println(res, len(res))
		text := ""
		if res != "" {
			text = res
		} else {
			text = locText("noActiveTask")
		}
		if MessageID == 0 {
			msg := tgBotApi.NewMessage(myID, text)
			msg.ParseMode = "Markdown"
			res, err := bot.Send(msg)
			dropErr(err)
			if text == locText("noActiveTask") {
				ticker.Stop()
				return -1
			} else {
				return res.MessageID
			}
		} else {
			if text == locText("noActiveTask") {
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
			_, _ = bot.DeleteMessage(tgBotApi.DeleteMessageConfig{
				ChatID:    myID,
				MessageID: chatMsgID,
			})
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

var tBot *tgBotApi.BotAPI

func sendAutoUpdateMessage() func(text string) {
	MessageID := 0
	myID := toInt64(info.UserID)
	return func(text string) {
		if MessageID == 0 {
			msg := tgBotApi.NewMessage(myID, text)
			msg.ParseMode = "Markdown"
			res, err := tBot.Send(msg)
			dropErr(err)
			MessageID = res.MessageID
		} else {
			if text != "close" {
				newMsg := tgBotApi.NewEditMessageText(myID, MessageID, text)
				newMsg.ParseMode = "Markdown"
				tBot.Send(newMsg)
			} else {
				tBot.Send(tgBotApi.NewDeleteMessage(myID, MessageID))
			}
		}
		return
	}
}
func createKeyBoardRow(texts ...string) [][]tgBotApi.KeyboardButton {
	Keyboards := make([][]tgBotApi.KeyboardButton, 0)
	for _, text := range texts {
		Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
			tgBotApi.NewKeyboardButton(text),
		))
	}
	return Keyboards
}
func createFilesInlineKeyBoardRow(filesInfos ...FilesInlineKeyboards) ([][]tgBotApi.InlineKeyboardButton, string) {
	Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
	text := ""
	index := 1
	inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
	for _, filesInfo := range filesInfos {
		for _, GidAndName := range filesInfo.GidAndName {

			text += fmt.Sprintf("%d: `%s`\n", index, GidAndName["Name"])
			inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(index), GidAndName["GID"]+":"+filesInfo.Data))
			if index%7 == 0 {
				Keyboards = append(Keyboards, inlineKeyBoardRow)
				inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
			}
			index++
		}
	}
	if len(inlineKeyBoardRow) != 0 {
		Keyboards = append(Keyboards, inlineKeyBoardRow)
	}
	if text == "" {
		text = " "
	}
	return Keyboards, text[:len(text)-1]
}

func createFunctionInlineKeyBoardRow(functionInfos ...FunctionInlineKeyboards) []tgBotApi.InlineKeyboardButton {
	Keyboards := make([]tgBotApi.InlineKeyboardButton, 0)
	for _, functionInfo := range functionInfos {
		Keyboards = append(Keyboards, tgBotApi.NewInlineKeyboardButtonData(functionInfo.Describe, "ALL:"+functionInfo.Describe))
	}
	return Keyboards
}

func tgBot(BotKey string, wg *sync.WaitGroup) {
	Keyboards := make([][]tgBotApi.KeyboardButton, 0)
	Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
		tgBotApi.NewKeyboardButton(locText("nowDownload")),
		tgBotApi.NewKeyboardButton(locText("nowWaiting")),
		tgBotApi.NewKeyboardButton(locText("nowOver")),
	))
	Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
		tgBotApi.NewKeyboardButton(locText("pauseTask")),
		tgBotApi.NewKeyboardButton(locText("resumeTask")),
		tgBotApi.NewKeyboardButton(locText("removeTask")),
	))
	if isLocal(info.Aria2Server) {
		Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
			tgBotApi.NewKeyboardButton(locText("removeDownloadFolderFiles")),
		))
		Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
			tgBotApi.NewKeyboardButton(locText("uploadDownloadFolderFiles")),
		))
		Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
			tgBotApi.NewKeyboardButton(locText("moveDownloadFolderFiles")),
		))
	}

	var numericKeyboard = tgBotApi.NewReplyKeyboard(Keyboards...)

	bot, err := tgBotApi.NewBotAPI(BotKey)
	dropErr(err)
	tBot = bot
	bot.Debug = false

	log.Printf(locText("authorizedAccount"), bot.Self.UserName)
	defer wg.Done()
	// go receiveMessage(msgChan)
	go SuddenMessage(bot)
	go TMSelectMessage(bot)
	u := tgBotApi.NewUpdate(0)
	u.Timeout = 60
	setCommands(bot)
	updates, err := bot.GetUpdatesChan(u)
	dropErr(err)
	for update := range updates {
		if update.CallbackQuery != nil {
			task := strings.Split(update.CallbackQuery.Data, ":")
			//log.Println(task)
			switch task[1] {
			case "1":
				aria2Rpc.Pause(task[0])
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("taskNowStop")))
			case "2":
				aria2Rpc.Unpause(task[0])
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("taskNowResume")))
			case "3":
				aria2Rpc.ForceRemove(task[0])
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("taskNowRemove")))
			case "4":
				aria2Rpc.PauseAll()
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("taskNowStopAll")))
			case "5":
				aria2Rpc.UnpauseAll()
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("taskNowResumeAll")))
			case "6":
				TMSelectMessageChan <- task[0]
				b := strings.Split(task[0], "~")
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("selected")+b[1]))
			case "7":
				TMSelectMessageChan <- task[0]
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("operationSuccess")))
			case "8":
				FileControlChan <- task[0]
				b := strings.Split(task[0], "~")
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("selected")+b[1]))
			case "9":
				FileControlChan <- task[0]
				bot.AnswerCallbackQuery(tgBotApi.NewCallback(update.CallbackQuery.ID, locText("operationSuccess")))
			}

			//fmt.Print(update)

			//bot.Send(tgBotApi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message != nil { //
			msg := tgBotApi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "Markdown"
			if fmt.Sprint(update.Message.Chat.ID) == info.UserID {

				// 创建新的MessageConfig。我们还没有文本，所以将其留空。

				switch update.Message.Text {
				case locText("nowDownload"):
					ticker := time.NewTicker(500 * time.Millisecond)
					rand.Seed(time.Now().UnixNano())
					a := rand.Intn(100000)
					activeRefreshControl = a
					go activeRefresh(update.Message.MessageID, bot, ticker, a)
				case locText("nowWaiting"):
					res := formatTellSomething(aria2Rpc.TellWaiting(0, info.MaxIndex))
					if res != "" {
						msg.Text = res
					} else {
						msg.Text = locText("noWaitingTask")
					}
				case locText("nowOver"):
					res := formatTellSomething(aria2Rpc.TellStopped(0, info.MaxIndex))
					if res != "" {
						msg.Text = res
					} else {
						msg.Text = locText("noOverTask")
					}
				case locText("pauseTask"):
					InlineKeyboards, text := createFilesInlineKeyBoardRow(FilesInlineKeyboards{
						GidAndName: formatGidAndName(aria2Rpc.TellActive()),
						Data:       "1",
					})
					if len(InlineKeyboards) != 0 {
						msg.Text = locText("stopWhichOne") + "\n" + text
						if len(InlineKeyboards) > 1 {
							InlineKeyboards = append(InlineKeyboards, createFunctionInlineKeyBoardRow(FunctionInlineKeyboards{
								Describe: locText("StopAll"),
								Data:     "4",
							}))
						}
						msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(InlineKeyboards...)
					} else {
						msg.Text = locText("noWaitingTask")
					}
				case locText("resumeTask"):

					InlineKeyboards, text := createFilesInlineKeyBoardRow(FilesInlineKeyboards{
						GidAndName: formatGidAndName(aria2Rpc.TellWaiting(0, info.MaxIndex)),
						Data:       "2",
					})
					if len(InlineKeyboards) != 0 {
						msg.Text = locText("resumeWhichOne") + "\n" + text
						if len(InlineKeyboards) > 1 {
							InlineKeyboards = append(InlineKeyboards, createFunctionInlineKeyBoardRow(FunctionInlineKeyboards{
								Describe: locText("ResumeAll"),
								Data:     "5",
							}))
						}
						msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(InlineKeyboards...)
					} else {
						msg.Text = locText("noActiveTask")
					}
				case locText("removeTask"):

					InlineKeyboards, text := createFilesInlineKeyBoardRow(FilesInlineKeyboards{
						GidAndName: formatGidAndName(aria2Rpc.TellActive()),
						Data:       "3",
					}, FilesInlineKeyboards{
						GidAndName: formatGidAndName(aria2Rpc.TellWaiting(0, info.MaxIndex)),
						Data:       "3",
					})
					if len(InlineKeyboards) != 0 {
						msg.Text = locText("removeWhichOne") + "\n" + text
						msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(InlineKeyboards...)
					} else {
						msg.Text = locText("noOverTask")
					}
				case locText("removeDownloadFolderFiles"):
					//dropErr(removeContents(info.DownloadFolder))
					isFileChanClean := false
					for !isFileChanClean {
						select {
						case _ = <-FileControlChan:
						default:
							isFileChanClean = true
						}
					}
					FileControlChan <- "close"
					go removeFiles(update.Message.MessageID, bot)
					FileControlChan <- "file"
				case locText("moveDownloadFolderFiles"):
					isFileChanClean := false
					for !isFileChanClean {
						select {
						case _ = <-FileControlChan:
						default:
							isFileChanClean = true
						}
					}
					FileControlChan <- "close"
					go copyFiles(update.Message.MessageID, bot)
					FileControlChan <- "file"
				case locText("uploadDownloadFolderFiles"):
					isFileChanClean := false
					for !isFileChanClean {
						select {
						case _ = <-FileControlChan:
						default:
							isFileChanClean = true
						}
					}
					FileControlChan <- "close"
					go uploadFiles(update.Message.MessageID, update.Message.Text, bot)
					FileControlChan <- "upload"
				default:
					if strings.Contains(update.Message.Text, "localhost/onedrive-login") {
						//如果是OneDrive auth code 链接
						createDriveInfoFolder("./info/onedrive")
						var re *regexp.Regexp
						if len(update.Message.Text) > 100 {
							re = regexp.MustCompile(`(?m)code=(.*?)&`)
						} else {
							re = regexp.MustCompile(`(?m)code=(.*?)`)
						}
						judgeLegal := re.FindStringSubmatch(update.Message.Text)
						//log.Println(judgeLegal)
						if len(judgeLegal) >= 2 {
							isFileChanClean := false
							for !isFileChanClean {
								select {
								case _ = <-FileControlChan:
								default:
									isFileChanClean = true
								}
							}
							FileControlChan <- "close"
							go uploadFiles(update.Message.MessageID, update.Message.Text, bot)
							FileControlChan <- "onedrive~create"
						} else {
							msg.Text = locText("errorOneDriveAuthURL")
						}

					} else if strings.Contains(update.Message.Text, "4/1AY") && len(update.Message.Text) == 62 {
						//如果是Google Drive auth code
						createDriveInfoFolder("./info/googleDrive")
						FileControlChan <- "close"
						go uploadFiles(update.Message.MessageID, update.Message.Text, bot)
						FileControlChan <- "googleDrive~create"
					} else if !download(update.Message.Text) {
						msg.Text = locText("unknownLink")
					}
					if update.Message.Document != nil {
						bt, _ := bot.GetFileDirectURL(update.Message.Document.FileID)
						resp, err := http.Get(bt)
						dropErr(err)
						defer resp.Body.Close()
						out, err := os.Create("temp.torrent")
						dropErr(err)
						defer out.Close()
						_, err = io.Copy(out, resp.Body)
						dropErr(err)
						if download("temp.torrent") {
							msg.Text = ""
						}
					}

				}

				// 从消息中提取命令。
				switch update.Message.Command() {
				case "start":
					version, err := aria2Rpc.GetVersion()
					dropErr(err)
					msg.Text = fmt.Sprintf(locText("commandStartRes"), info.Sign, version.Version)
					if isLocal(info.Aria2Server) {
						msg.Text += "\n" + locText("inLocal")
					}
					//msg.Text += "\n" + locText("nowTMMode") + locText("tmMode"+aria2Set.TMMode)
					msg.ReplyMarkup = numericKeyboard
				case "help":
					msg.Text = locText("commandHelpRes")
				case "myid":
					msg.Text = fmt.Sprintf(locText("commandMyIDRes"), update.Message.Chat.ID)
				}
			} else {
				msg.Text = locText("doNotHavePermissionControl")
				if update.Message.Command() == "myid" {
					msg.Text = fmt.Sprintf(locText("commandMyIDRes"), update.Message.Chat.ID)
				}
			}

			if msg.Text != "" {
				//bot.Send(tgBotApi.NewEditMessageText(update.Message.Chat.ID, 591, "123456"))
				_, err := bot.Send(msg)
				dropErr(err)
			}
		}
	}
}
