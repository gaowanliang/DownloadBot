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

// SuddenMessageChan is a channel for reminding when an emergency message occurs, such as download start and download end
var SuddenMessageChan = make(chan string, 3)

// TMSelectMessageChan is a channel for reminding when you upload a torrent file or send a magnet(TM is Torrent/Magnet)
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

// SuddenMessage is a function for dealing when an emergency message occurs, such as download start and download end
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

// TMSelectMessage is a function for dealing when you upload a torrent file or send a magnet,this function will interrupt the download and output the file list of the torrent file or magnet (TM is Torrent/Magnet)
func TMSelectMessage(bot *tgBotApi.BotAPI) {
	var MessageID []int
	myID := toInt64(info.UserID)
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
				setTMDownloadFilesAndStart(gid, selectFileList)
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
				aria2Rpc.ForceRemove(gid)
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
				biggestFileIndex := selectBiggestFile(gid)
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
				bigFilesIndex := selectBigFiles(gid)
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
				i := toInt(b[1])
				i--                         // sequence num displayed is more than the subscript,so reduce it
				if downloadFilesCount > 1 { // make sure that at least one file will be downloaded
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

			for i, val := range selectFileList { // check wheathe all files under all nodes select status
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
			fileList := formatTMFiles(gid)
			index := 1
			for i, file := range fileList {
				pathClass(fmt.Sprintf("%s|%s|%d", file[0], file[1], i+1), &directoryTree)
				//log.Printf("%s|%s|%d\n", file[0], file[1], i+1)
				index++
			}
		}

		text := fmt.Sprintf("%s %s\n", tellName(aria2Rpc.TellStatus(gid)), locText("fileDirectoryIsAsFollows"))
		Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
		inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)

		fileListGoTree, _, _ := generateGoTree(directoryTree, 0, &selectFileList)
		fileListTreeLine := strings.Split(fileListGoTree[0].Print(), "\n")
		fileListTreeLineCount := len(fileListTreeLine)
		characterCount := len(text)
		r, err := regexp.Compile(`[✅⬜](\d+)`)
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
				startAndEndIndex[1] = toInt(res[1])
			}
		}

		// log.Println(text)
		text += locText("pleaseSelectTheFileYouWantToDownload")
		if msgCount > 2 {
			text += "\n" + locText("tmFileTooMany")
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
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("selectAll"), gid+"~selectAll"+":7"))
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(locText("cancel"), gid+"~cancel"+":7"))
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
