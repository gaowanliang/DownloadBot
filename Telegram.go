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

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⬇️ 正在下载"),
		tgbotapi.NewKeyboardButton("⌛️ 正在等待"),
		tgbotapi.NewKeyboardButton("✅ 已完成/已停止"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⏸️ 停止任务"),
		tgbotapi.NewKeyboardButton("▶️ 继续任务"),
		tgbotapi.NewKeyboardButton("❌ 移除任务"),
	),
)

func setCommands(bot *tgbotapi.BotAPI) {
	bot.SetMyCommands([]tgbotapi.BotCommand{
		{
			Command:     "start",
			Description: "获取已上线的Aria2服务器，并打开面板",
		}, {
			Command:     "myid",
			Description: "获取user-id",
		},
	})
}

// SuddenMessage receive active requests from WebSocket
func SuddenMessage(bot *tgbotapi.BotAPI) {
	for {
		a := <-SuddenMessageChan
		//log.Println("通道进入")
		//time.Sleep(time.Second * 5)
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
	bot, err := tgbotapi.NewBotAPI(BotKey)
	dropErr(err)

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)
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
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "任务已停止"))
			case "2":
				aria2Rpc.Unpause(task[0])
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "任务已恢复"))
			case "3":
				aria2Rpc.ForceRemove(task[0])
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "任务已移除"))
			case "4":
				aria2Rpc.PauseAll()
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "任务已全部停止"))
			case "5":
				aria2Rpc.UnpauseAll()
				bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "任务已全部恢复"))
			}
			//fmt.Print(update)

			//bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data))
		}

		if update.Message != nil { //

			// 创建新的MessageConfig。我们还没有文本，所以将其留空。
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "Markdown"
			// 从消息中提取命令。
			switch update.Message.Command() {
			case "start":
				version, err := aria2Rpc.GetVersion()
				dropErr(err)
				msg.Text = fmt.Sprintf("%s 当前已连接，版本: %s ，请选择一个选项", info.Sign, version.Version)
				msg.ReplyMarkup = numericKeyboard

			case "help":
				msg.Text = "🤖 一个控制你的Aria2服务器的Telegram Bot。"
			case "myid":
				msg.Text = fmt.Sprintf("你的user-id为 `%d` ", update.Message.Chat.ID)
			case "status":
				msg.Text = "I'm ok."
				//default:
				//msg.Text = "I don't know that command"
			}

			switch update.Message.Text {
			case "⬇️ 正在下载":
				res := formatTellSomething(aria2Rpc.TellActive())
				if res != "" {
					msg.Text = res
				} else {
					// log.Println(aria2Rpc.TellStatus("42fa911166acf119"))
					msg.Text = "没有正在进行的任务！"
				}
			case "⌛️ 正在等待":
				res := formatTellSomething(aria2Rpc.TellWaiting(0, info.MaxIndex))
				if res != "" {
					msg.Text = res
				} else {
					msg.Text = "没有正在等待的任务！"
				}
			case "✅ 已完成/已停止":
				res := formatTellSomething(aria2Rpc.TellStopped(0, info.MaxIndex))
				if res != "" {
					msg.Text = res
				} else {
					msg.Text = "没有已完成/已停止的任务！"
				}
			case "⏸️ 停止任务":
				InlineKeyboards := make([]tgbotapi.InlineKeyboardButton, 0)
				for _, value := range formatGidAndName(aria2Rpc.TellActive()) {
					log.Printf("%s %s", value["GID"], value["Name"])
					InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData(value["Name"], value["GID"]+":1"))
				}
				if len(InlineKeyboards) != 0 {
					msg.Text = "停止哪一个?"
					if len(InlineKeyboards) > 1 {
						InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData("停止全部", "ALL:4"))
					}
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineKeyboards)
				} else {
					msg.Text = "没有正在等待的任务！"
				}
			case "▶️ 继续任务":
				InlineKeyboards := make([]tgbotapi.InlineKeyboardButton, 0)
				for _, value := range formatGidAndName(aria2Rpc.TellWaiting(0, info.MaxIndex)) {
					log.Printf("%s %s", value["GID"], value["Name"])
					InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData(value["Name"], value["GID"]+":2"))

				}
				if len(InlineKeyboards) != 0 {
					msg.Text = "恢复哪一个?"
					if len(InlineKeyboards) > 1 {
						InlineKeyboards = append(InlineKeyboards, tgbotapi.NewInlineKeyboardButtonData("恢复全部", "ALL:5"))
					}
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineKeyboards)
				} else {
					msg.Text = "没有正在下载的任务"
				}
			case "❌ 移除任务":
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
					msg.Text = "移除哪一个?"
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineKeyboards)
				} else {
					msg.Text = "没有已完成/已停止的任务"
				}
			default:
				if !download(update.Message.Text) {
					msg.Text = "未知的下载链接，请重新检查"
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
			if msg.Text != "" {
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
