package main

import (
	"DownloadBot/cmd/client"
	"DownloadBot/cmd/server"
	i18nLoc "DownloadBot/i18n"
	"DownloadBot/internal/config"
	"DownloadBot/tool/input"
	"DownloadBot/tool/input/aria2"
	"DownloadBot/tool/output/telegram"
	logger "DownloadBot/tool/zap"
	"flag"
	"sync"
)

var (
	configFile = flag.String("c", "./config.json", "config file path")
	clientSign = flag.String("s", "", "client sign")
)

func init() {
	flag.Parse()
	// Init config
	config.InitConfig(*configFile, *clientSign)

	// Init i18n
	i18nLoc.LocLan(config.GetLanguage())

	// Init logger
	logger.InitLog(config.GetLogPath(), config.GetErrPath(), config.GetLogLevel(), i18nLoc.LocText)
	logger.Info(i18nLoc.LocText("configCompleted"))

}

func main() {
	var wg sync.WaitGroup
	//var s *grpc.Server

	if config.IsMasterServer() && config.IsServer() {
		go server.StartServer()
		//defer s.Stop()
		wg.Add(1)
		go telegram.Aria2Bot(config.GetTelegramBotKey(), &wg)
		wg.Wait()
	} else if !config.IsServer() {
		client.Method.InitClient()
		input.ToolApp.Aria2.Notifier = aria2.Notifier{}
		go input.SuddenlyMsg()
		input.ToolApp.Aria2.Load(input.ToolApp.Aria2.Notifier, func(gid string) {
			// client.Method.TMStop(gid)
		}, true)

	}

}
