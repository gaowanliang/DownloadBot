package telegram

import (
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	"DownloadBot/internal/server/clientManage"
	"DownloadBot/tool/controller"
	"DownloadBot/tool/displayUtil/gotree"
	"DownloadBot/tool/input"
	"DownloadBot/tool/monitor"
	"DownloadBot/tool/typeTrans"
	logger "DownloadBot/tool/zap"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	tgBotApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wxnacy/wgo/arrays"
)

var deleteMes tgBotApi.DeleteMessageConfig = tgBotApi.DeleteMessageConfig{
	ChannelUsername: "",
	ChatID:          0,
	MessageID:       0,
}

func dropErr(err error) {
	if err != nil {
		logger.Panic("%w", err)
	}
}

func setCommands() tgBotApi.SetMyCommandsConfig {
	return tgBotApi.NewSetMyCommands(tgBotApi.BotCommand{
		Command:     "start",
		Description: i18nLoc.LocText("tgCommandStartDes"),
	}, tgBotApi.BotCommand{
		Command:     "myid",
		Description: i18nLoc.LocText("tgCommandMyIDDes"),
	}, tgBotApi.BotCommand{
		Command:     "change_client",
		Description: i18nLoc.LocText("watch now online client and change client"),
	})
	//, tgBotApi.BotCommand{
	//		Command:     "setMaxLength",
	//		Description: locText("tgCommandSetMaxLengthDes"),
	//	}
}

func SendTelegramAutoUpdateMessage() func(text string) {
	MessageID := 0
	myID := typeTrans.Str2Int64(config.GetTelegramUserID())
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

func SendTelegramSuddenMessage(text string) {
	myID := typeTrans.Str2Int64(config.GetTelegramUserID())
	msg := tgBotApi.NewMessage(myID, text)
	msg.ParseMode = "Markdown"
	tBot.Send(msg)
}

var tBot *tgBotApi.BotAPI

func createKeyBoardRow(texts ...string) [][]tgBotApi.KeyboardButton {
	Keyboards := make([][]tgBotApi.KeyboardButton, 0)
	for _, text := range texts {
		Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
			tgBotApi.NewKeyboardButton(text),
		))
	}
	return Keyboards
}
func createFilesInlineKeyBoardRow(filesInfos ...filesInlineKeyboards) ([][]tgBotApi.InlineKeyboardButton, string) {
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

//
// createClientSelectInlineKeyBoardRow
//  @Description: create client select inline keyboard
//  @param clientNum int "client number"
//  @return [][]tgBotApi.InlineKeyboardButton
//
func createClientSelectInlineKeyBoardRow(clientNum int, delMsgSign string) [][]tgBotApi.InlineKeyboardButton {
	Keyboards := make([][]tgBotApi.InlineKeyboardButton, 0)
	inlineKeyBoardRow := make([]tgBotApi.InlineKeyboardButton, 0)
	for i := 1; i <= clientNum; i++ {
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(fmt.Sprint(i), fmt.Sprintf("%d:10", i)))
		if i%7 == 0 {
			Keyboards = append(Keyboards, inlineKeyBoardRow)
			inlineKeyBoardRow = make([]tgBotApi.InlineKeyboardButton, 0)
		}

	}

	if len(inlineKeyBoardRow) != 0 {
		//delete this msg button
		inlineKeyBoardRow = append(inlineKeyBoardRow, tgBotApi.NewInlineKeyboardButtonData(i18nLoc.LocText("cancel"), delMsgSign+":11"))
		Keyboards = append(Keyboards, inlineKeyBoardRow)
	} else {
		Keyboards = append(Keyboards, []tgBotApi.InlineKeyboardButton{
			tgBotApi.NewInlineKeyboardButtonData(i18nLoc.LocText("cancel"), delMsgSign+":11"),
		})
	}
	return Keyboards
}

func createFunctionInlineKeyBoardRow(functionInfos ...functionInlineKeyboards) []tgBotApi.InlineKeyboardButton {
	Keyboards := make([]tgBotApi.InlineKeyboardButton, 0)
	for _, functionInfo := range functionInfos {
		Keyboards = append(Keyboards, tgBotApi.NewInlineKeyboardButtonData(functionInfo.Describe, "ALL:"+functionInfo.Describe))
	}
	return Keyboards
}

func Aria2Bot(BotKey string, wg *sync.WaitGroup) {
	Keyboards := make([][]tgBotApi.KeyboardButton, 0)
	Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("nowDownload")),
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("nowWaiting")),
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("nowOver")),
	))
	Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("pauseTask")),
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("resumeTask")),
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("removeTask")),
	))

	Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("removeDownloadFolderFiles")),
	))
	Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("uploadDownloadFolderFiles")),
	))
	Keyboards = append(Keyboards, tgBotApi.NewKeyboardButtonRow(
		tgBotApi.NewKeyboardButton(i18nLoc.LocText("moveDownloadFolderFiles")),
	))

	var numericKeyboard = tgBotApi.NewReplyKeyboard(Keyboards...)

	bot, err := tgBotApi.NewBotAPI(BotKey)
	dropErr(err)
	tBot = bot
	bot.Debug = false
	input.ToolApp.Aria2.Load(Notifier{}, func(gid string) {
		TMSelectMessageChan <- gid
	}, false)

	logger.Info(fmt.Sprintf(i18nLoc.LocText("authorizedAccount"), bot.Self.UserName))
	defer wg.Done()
	// go receiveMessage(msgChan)
	go SuddenMessage(bot)
	go Aria2TMSelectMsg(bot)
	u := tgBotApi.NewUpdate(0)
	u.Timeout = 60
	_, err = bot.Request(setCommands())
	//setCommands(bot)
	updates := bot.GetUpdatesChan(u)
	dropErr(err)
	nowCtrlClient := 0 // 0 is Master Server
	for update := range updates {
		if update.CallbackQuery != nil {
			task := strings.Split(update.CallbackQuery.Data, ":")
			//log.Println(task)
			switch task[1] {
			case "1":
				input.PauseTask(task[0])
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("taskNowStop")))
			case "2":
				input.UnpauseTask(task[0])
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("taskNowResume")))
			case "3":
				input.ForceRemoveTask(task[0])
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("taskNowRemove")))
			case "4":
				input.PauseAllTask()
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("taskNowStopAll")))
			case "5":
				input.UnpauseAllTask()
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("taskNowResumeAll")))
			case "6":
				TMSelectMessageChan <- task[0]
				b := strings.Split(task[0], "~")
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("selected")+b[1]))
			case "7":
				TMSelectMessageChan <- task[0]
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("operationSuccess")))
			case "8":
				FileControlChan <- task[0]
				b := strings.Split(task[0], "~")
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("selected")+b[1]))
			case "9":
				FileControlChan <- task[0]
				bot.Request(tgBotApi.NewCallback(update.CallbackQuery.ID, i18nLoc.LocText("operationSuccess")))
			case "10":
				// change Client
				nowCtrlClient, _ = strconv.Atoi(task[0])
				logger.Debug("change Client:" + task[0])
				logger.Debug("Client Name:" + clientManage.GetClientName(nowCtrlClient-1))

				// send change client msg
				if nowCtrlClient != 0 {
					bot.Send(tgBotApi.NewMessage(typeTrans.Str2Int64(config.GetTelegramUserID()), fmt.Sprintf(i18nLoc.LocText("Switched to client: %v"), clientManage.GetClientName(nowCtrlClient-1))))

				} else {
					bot.Send(tgBotApi.NewMessage(typeTrans.Str2Int64(config.GetTelegramUserID()), i18nLoc.LocText("Switched to master server")))
				}

				bot.Request(tgBotApi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
			case "11":
				// delete msg
				bot.Request(tgBotApi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID))
			}

			//fmt.Print(update)

			//bot.Send(tgBotApi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message != nil { //
			msg := tgBotApi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "Markdown"
			if arrays.ContainsString(strings.Split(config.GetTelegramUserID(), ","), fmt.Sprint(update.Message.Chat.ID)) != -1 {

				// 创建新的MessageConfig

				switch update.Message.Text {
				case i18nLoc.LocText("nowDownload"):
					ticker := time.NewTicker(500 * time.Millisecond)
					rand.Seed(time.Now().UnixNano())
					a := rand.Intn(100000)
					activeRefreshControl = a
					go activeRefresh(update.Message.MessageID, bot, ticker, a)
				case i18nLoc.LocText("nowWaiting"):
					res := input.ToolApp.Aria2.FormatTellWaiting()
					//res := aria2.FormatTellWaiting()
					if res != "" {
						msg.Text = res
					} else {
						msg.Text = i18nLoc.LocText("noWaitingTask")
					}
				case i18nLoc.LocText("nowOver"):
					res := input.ToolApp.Aria2.FormatTellStopped()
					if res != "" {
						msg.Text = res
					} else {
						msg.Text = i18nLoc.LocText("noOverTask")
					}
				case i18nLoc.LocText("pauseTask"):
					InlineKeyboards, text := createFilesInlineKeyBoardRow(filesInlineKeyboards{
						GidAndName: input.ToolApp.Aria2.FormatGidAndName(0),
						Data:       "1",
					})
					if len(InlineKeyboards) != 0 {
						msg.Text = i18nLoc.LocText("stopWhichOne") + "\n" + text
						if len(InlineKeyboards) > 1 {
							InlineKeyboards = append(InlineKeyboards, createFunctionInlineKeyBoardRow(functionInlineKeyboards{
								Describe: i18nLoc.LocText("StopAll"),
								Data:     "4",
							}))
						}
						msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(InlineKeyboards...)
					} else {
						msg.Text = i18nLoc.LocText("noWaitingTask")
					}
				case i18nLoc.LocText("resumeTask"):

					InlineKeyboards, text := createFilesInlineKeyBoardRow(filesInlineKeyboards{
						GidAndName: input.ToolApp.Aria2.FormatGidAndName(1),
						Data:       "2",
					})
					if len(InlineKeyboards) != 0 {
						msg.Text = i18nLoc.LocText("resumeWhichOne") + "\n" + text
						if len(InlineKeyboards) > 1 {
							InlineKeyboards = append(InlineKeyboards, createFunctionInlineKeyBoardRow(functionInlineKeyboards{
								Describe: i18nLoc.LocText("ResumeAll"),
								Data:     "5",
							}))
						}
						msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(InlineKeyboards...)
					} else {
						msg.Text = i18nLoc.LocText("noActiveTask")
					}
				case i18nLoc.LocText("removeTask"):

					InlineKeyboards, text := createFilesInlineKeyBoardRow(filesInlineKeyboards{
						GidAndName: input.ToolApp.Aria2.FormatGidAndName(0),
						Data:       "3",
					}, filesInlineKeyboards{
						GidAndName: input.ToolApp.Aria2.FormatGidAndName(1),
						Data:       "3",
					})
					if len(InlineKeyboards) != 0 {
						msg.Text = i18nLoc.LocText("removeWhichOne") + "\n" + text
						msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(InlineKeyboards...)
					} else {
						msg.Text = i18nLoc.LocText("noOverTask")
					}
				case i18nLoc.LocText("removeDownloadFolderFiles"):
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
					go removeFilesPrint(update.Message.MessageID, bot)
					FileControlChan <- "file"
				case i18nLoc.LocText("moveDownloadFolderFiles"):
					isFileChanClean := false
					for !isFileChanClean {
						select {
						case _ = <-FileControlChan:
						default:
							isFileChanClean = true
						}
					}
					FileControlChan <- "close"
					go copyFilesPrint(update.Message.MessageID, bot)
					FileControlChan <- "file"
				case i18nLoc.LocText("uploadDownloadFolderFiles"):
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
						// OneDrive auth code link
						controller.CreateDriveInfoFolder("./info/onedrive")
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
							msg.Text = i18nLoc.LocText("errorOneDriveAuthURL")
						}

					} else if strings.Contains(update.Message.Text, "4/1A") && len(update.Message.Text) == 62 {
						// Google Drive auth code
						controller.CreateDriveInfoFolder("./info/googleDrive")
						FileControlChan <- "close"
						go uploadFiles(update.Message.MessageID, update.Message.Text, bot)
						FileControlChan <- "googleDrive~create"
					} else if !input.ToolApp.Aria2.Download(update.Message.Text) {
						msg.Text = i18nLoc.LocText("unknownLink")
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
						if input.ToolApp.Aria2.Download("temp.torrent") {
							msg.Text = ""
						}
					}
					/*if update.Message.Video != nil {
						videoURL, _ := bot.GetFileDirectURL(update.Message.Video.FileID)
						log.Println(videoURL)
						videoInfo, err := bot.GetFile(tgBotApi.FileConfig{
							FileID: update.Message.Video.FileID,
						})
						dropErr(err)
						log.Println(videoInfo)
					}
					*/
				}

				// 从消息中提取命令。
				switch update.Message.Command() {
				case "start":
					msg.Text = fmt.Sprintf(i18nLoc.LocText("commandStartRes"), config.GetSign(), input.ToolApp.Aria2.GetVersion())
					if monitor.IsLocal(config.GetAria2Server()) {
						msg.Text += "\n" + i18nLoc.LocText("inLocal")
					}
					//msg.Text += "\n" + locText("nowTMMode") + locText("tmMode"+aria2Set.TMMode)
					msg.ReplyMarkup = numericKeyboard
				case "help":
					msg.Text = i18nLoc.LocText("commandHelpRes")
				case "myid":
					msg.Text = fmt.Sprintf(i18nLoc.LocText("commandMyIDRes"), update.Message.Chat.ID)
				case "setMaxLength":
					i, err := strconv.Atoi(strings.ReplaceAll(update.Message.Text, "/setMaxLength ", ""))
					if err != nil {
						msg.Text = i18nLoc.LocText("commandSetMaxLengthHelpRes")
					} else {
						goTree.SetMaxLength(i)
						msg.Text = i18nLoc.LocText("operationSuccess")
					}
				case "change_client":
					clientListStr, clientQuantity := clientManage.ShowClientList()
					if clientQuantity > 0 {
						msg.Text = i18nLoc.LocText("Now online client is:\n\n") + clientListStr + i18nLoc.LocText("\nPlease click the client number to switch")
						msg.ReplyMarkup = tgBotApi.NewInlineKeyboardMarkup(createClientSelectInlineKeyBoardRow(clientQuantity, "123456")...)
					} else {
						msg.Text = i18nLoc.LocText("There is only Master Server, no need to switch")
					}

					// ! not finished
					msg.Text = "\n\n\nThe switching function is still under development, and the current switching is invalid"
				}
			} else {
				msg.Text = i18nLoc.LocText("doNotHavePermissionControl")
				if update.Message.Command() == "myid" {
					msg.Text = fmt.Sprintf(i18nLoc.LocText("commandMyIDRes"), update.Message.Chat.ID)
				}
			}

			if msg.Text != "" {
				//bot.Send(tgBotApi.NewEditMessageText(update.Message.Chat.ID, 591, "123456"))
				_, err := bot.Send(msg)
				dropErr(err)
			}
		}
	}
	input.ToolApp.Aria2.Close()
}
