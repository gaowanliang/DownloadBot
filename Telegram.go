package main

import (
	"fmt"

	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SuddenMessageChan receive active requests from WebSocket
var SuddenMessageChan = make(chan string, 3)

func setCommands(bot *tgbotapi.BotAPI) {
	bot.SetMyCommands([]tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: locText("tgCommandStartDes"),
		}, {
			Command:     "myid",
			Description: locText("tgCommandMyidDes"),
		},
	})
}

// SuddenMessage receive active requests from WebSocket
func SuddenMessage(bot *tgbotapi.BotAPI) {
	for {
		a := <-SuddenMessageChan
		gid := a[2:18]
		a = strings.ReplaceAll(a, gid, tellName(aria2Rpc.TellStatus(gid)))
		myID, err := strconv.ParseInt(info.UserID, 10, 64)
		dropErr(err)
		msg := tgbotapi.NewMessage(myID, a)
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func tgBot(BotKey string, wg *sync.WaitGroup) {

	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(locText("nowDownload")),
			tgbotapi.NewKeyboardButton(locText("nowWaiting")),
			tgbotapi.NewKeyboardButton(locText("nowOver")),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(locText("pauseTask")),
			tgbotapi.NewKeyboardButton(locText("resumeTask")),
			tgbotapi.NewKeyboardButton(locText("removeTask")),
		),
	)

	bot, err := tgbotapi.NewBotAPI(BotKey)
	dropErr(err)

	bot.Debug = false

	log.Printf(locText("authorizedAccount"), bot.Self.UserName)
	defer wg.Done()
	// go receiveMessage(msgChan)
	go SuddenMessage(bot)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	setCommands(bot)
	updates, err := bot.GetUpdatesChan(u)
	dropErr(err)

	for update := range updates {
		if update.CallbackQuery != nil {
			task := strings.Split(update.CallbackQuery.Data, ":")
			log.Println(task)
			switch task[1] {
			case "1":
				aria2Rpc.Pause(task[0])
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, locText("taskNowStop")))
			case "2":
				aria2Rpc.Unpause(task[0])
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, locText("taskNowResume")))
			case "3":
				aria2Rpc.ForceRemove(task[0])
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, locText("taskNowRemove")))
			case "4":
				aria2Rpc.PauseAll()
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, locText("taskNowStopAll")))
			case "5":
				aria2Rpc.UnpauseAll()
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, locText("taskNowResumeAll")))
			}
			//fmt.Print(update)

			//bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message != nil { //

			// 创建新的MessageConfig。我们还没有文本，所以将其留空。
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "Markdown"

			switch update.Message.Text {
			case locText("nowDownload"):
				res := formatTellSomething(aria2Rpc.TellActive())
				if res != "" {
					msg.Text = res
				} else {
					// log.Println(aria2Rpc.TellStatus("42fa911166acf119"))
					msg.Text = locText("noActiveTask")
				}
			case locText("nowWaiting"):
				res := formatTellSomething(aria2Rpc.TellWaiting(0, info.MaxIndex))
				if res != "" {
					msg.Text = res
				} else {
					msg.Text = locText("noWaittingTask")
				}
			case locText("nowOver"):
				res := formatTellSomething(aria2Rpc.TellStopped(0, info.MaxIndex))
				if res != "" {
					msg.Text = res
				} else {
					msg.Text = locText("noOverTask")
				}
			case locText("pauseTask"):
				InlineKeyboards := make([]tgbotapi.InlineKeyboardButton, 0)
				for _, value := range formatGidAndName(aria2Rpc.TellActive()) {
					log.Printf("%s %s", value["GID"], value["Name"])
					InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData(value["Name"], value["GID"]+":1"))
				}
				if len(InlineKeyboards) != 0 {
					msg.Text = locText("stopWhichOne")
					if len(InlineKeyboards) > 1 {
						InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData("停止全部", "ALL:4"))
					}
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineKeyboards)
				} else {
					msg.Text = locText("noWaittingTask")
				}
			case locText("resumeTask"):
				InlineKeyboards := make([]tgbotapi.InlineKeyboardButton, 0)
				for _, value := range formatGidAndName(aria2Rpc.TellWaiting(0, info.MaxIndex)) {
					log.Printf("%s %s", value["GID"], value["Name"])
					InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData(value["Name"], value["GID"]+":2"))

				}
				if len(InlineKeyboards) != 0 {
					msg.Text = locText("resumeWhichOne")
					if len(InlineKeyboards) > 1 {
						InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData("恢复全部", "ALL:5"))
					}
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineKeyboards)
				} else {
					msg.Text = locText("noActiveTask")
				}
			case locText("removeTask"):
				InlineKeyboards := make([]tgbotapi.InlineKeyboardButton, 0)
				for _, value := range formatGidAndName(aria2Rpc.TellActive()) {
					log.Printf("%s %s", value["GID"], value["Name"])
					InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData(value["Name"], value["GID"]+":3"))
				}
				for _, value := range formatGidAndName(aria2Rpc.TellWaiting(0, info.MaxIndex)) {
					log.Printf("%s %s", value["GID"], value["Name"])
					InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData(value["Name"], value["GID"]+":3"))
				}
				if len(InlineKeyboards) != 0 {
					msg.Text = locText("removeWhichOne")
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineKeyboards)
				} else {
					msg.Text = locText("noOverTask")
				}
			default:
				if !download(update.Message.Text) {
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
				msg.ReplyMarkup = numericKeyboard

			case "help":
				msg.Text = locText("commandHelpRes")
			case "myid":
				msg.Text = fmt.Sprintf(locText("commandMyidRes"), update.Message.Chat.ID)
			}

			if msg.Text != "" {
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
